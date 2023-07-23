/*
 * This file is part of the kiagnose project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2023 Red Hat, Inc.
 *
 */

package trex

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	expect "github.com/google/goexpect"

	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/kiagnose/kubevirt-dpdk-checkup/pkg/internal/checkup/executor/console"
)

type Client struct {
	vmiSerialClient      vmiSerialConsoleClient
	namespace            string
	verbosePrintsEnabled bool
}

func NewClient(vmiSerialClient vmiSerialConsoleClient, namespace string, verbosePrintsEnabled bool) Client {
	return Client{
		vmiSerialClient:      vmiSerialClient,
		namespace:            namespace,
		verbosePrintsEnabled: verbosePrintsEnabled,
	}
}

func (t Client) WaitForServerToBeReady(ctx context.Context, vmiName string) error {
	const (
		interval = 5 * time.Second
		timeout  = time.Minute
	)
	var err error
	ctxWithNewDeadline, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	conditionFn := func(ctx context.Context) (bool, error) {
		if t.isServerRunning(vmiName) {
			log.Printf("trex-server is now ready")
			return true, nil
		}
		if t.verbosePrintsEnabled {
			log.Printf("trex-server is not yet ready...")
		}
		return false, nil
	}
	if err = wait.PollImmediateUntilWithContext(ctxWithNewDeadline, interval, conditionFn); err != nil {
		if !errors.Is(err, wait.ErrWaitTimeout) {
			return err
		}
		if t.verbosePrintsEnabled {
			if logErr := t.printTrexServiceFailLogs(vmiName); logErr != nil {
				return logErr
			}
		}
		return fmt.Errorf("timeout waiting for trex-server to be ready")
	}
	return nil
}

func (t Client) isServerRunning(vmiName string) bool {
	const helpSubstring = "Console Commands"
	resp, err := t.runTrexConsoleCmd(vmiName, "help")
	if err != nil || !strings.Contains(resp, helpSubstring) {
		return false
	}
	return true
}

func (t Client) printTrexServiceFailLogs(vmiName string) error {
	var err error
	trexServiceStatus, err := t.getTrexServiceStatus(vmiName)
	if err != nil {
		return fmt.Errorf("failed gathering systemctl service status after trex-server timeout: %w", err)
	}
	trexJournalctlLogs, err := t.getTrexServiceJournalctl(vmiName)
	if err != nil {
		return fmt.Errorf("failed gathering trex.service related joutnalctl logs after trex-server timeout: %w", err)
	}
	log.Printf("timeout waiting for trex-server to be ready\n"+
		"systemd service status:\n%s\n"+
		"joutnalctl logs:\n%s", trexServiceStatus, trexJournalctlLogs)
	return nil
}

func (t Client) getTrexServiceStatus(vmiName string) (string, error) {
	command := fmt.Sprintf("systemctl status %s | cat", SystemdUnitFileName)
	resp, err := console.SafeExpectBatchWithResponse(t.vmiSerialClient, t.namespace, vmiName,
		[]expect.Batcher{
			&expect.BSnd{S: command + "\n"},
			&expect.BExp{R: shellPrompt},
		},
		batchTimeout,
	)
	return resp[0].Output, err
}

func (t Client) getTrexServiceJournalctl(vmiName string) (string, error) {
	command := fmt.Sprintf("journalctl | grep %s", SystemdUnitFileName)
	resp, err := console.SafeExpectBatchWithResponse(t.vmiSerialClient, t.namespace, vmiName,
		[]expect.Batcher{
			&expect.BSnd{S: command + "\n"},
			&expect.BExp{R: shellPrompt},
		},
		batchTimeout,
	)
	return resp[0].Output, err
}

func (t Client) runTrexConsoleCmd(vmiName, command string) (string, error) {
	shellCommand := fmt.Sprintf("cd %s && echo %q | ./trex-console -q", BinDirectory, command)
	resp, err := console.SafeExpectBatchWithResponse(t.vmiSerialClient, t.namespace, vmiName,
		[]expect.Batcher{
			&expect.BSnd{S: shellCommand + "\n"},
			&expect.BExp{R: shellPrompt},
		},
		batchTimeout,
	)

	if err != nil {
		return "", err
	}
	return cleanStdout(resp[0].Output), nil
}

func cleanStdout(rawStdout string) string {
	stdout := strings.Replace(rawStdout, "Using 'python3' as Python interpeter", "", -1)
	stdout = strings.Replace(stdout, "-=TRex Console v3.0=-", "", -1)
	stdout = strings.Replace(stdout, "Type 'help' or '?' for supported actions", "", -1)
	stdout = strings.Replace(stdout, "trex>Global Statistitcs", "", -1)
	stdout = strings.Replace(stdout, "trex>", "", -1)

	return stdout
}
