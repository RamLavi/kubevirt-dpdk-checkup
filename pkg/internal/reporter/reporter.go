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

package reporter

import (
	"fmt"

	"k8s.io/client-go/kubernetes"

	kreporter "github.com/kiagnose/kiagnose/kiagnose/reporter"

	"github.com/kiagnose/kubevirt-dpdk-checkup/pkg/internal/status"
)

const (
	TrafficGeneratorTxPacketsKey       = "trafficGeneratorTxPackets"
	TrafficGeneratorOutErrorPacketsKey = "trafficGeneratorOutputErrorPackets"
	TrafficGeneratorInErrorPacketsKey  = "trafficGeneratorInErrorPackets"
	DPDKRxTestPacketsKey               = "DPDKRxTestPackets"
	DPDKRxDropsKey                     = "DPDKRxPacketDrops"
	DPDKTxDropsKey                     = "DPDKTxPacketDrops"
	TrafficGeneratorNodeKey            = "trafficGeneratorNode"
	DPDKVMNodeKey                      = "DPDKVMNode"
)

type Reporter struct {
	kreporter.Reporter
}

func New(c kubernetes.Interface, configMapNamespace, configMapName string) *Reporter {
	r := kreporter.New(c, configMapNamespace, configMapName)
	return &Reporter{*r}
}

func (r *Reporter) Report(checkupStatus status.Status) error {
	if !r.HasData() {
		return r.Reporter.Report(checkupStatus.Status)
	}

	checkupStatus.Succeeded = len(checkupStatus.FailureReason) == 0

	checkupStatus.Status.Results = formatResults(checkupStatus)

	return r.Reporter.Report(checkupStatus.Status)
}

func formatResults(checkupStatus status.Status) map[string]string {
	var emptyResults status.Results
	if checkupStatus.Results == emptyResults {
		return map[string]string{}
	}

	formattedResults := map[string]string{
		TrafficGeneratorTxPacketsKey:       fmt.Sprintf("%d", checkupStatus.Results.TrafficGeneratorTxPackets),
		TrafficGeneratorOutErrorPacketsKey: fmt.Sprintf("%d", checkupStatus.Results.TrafficGeneratorOutErrorPackets),
		TrafficGeneratorInErrorPacketsKey:  fmt.Sprintf("%d", checkupStatus.Results.TrafficGeneratorInErrorPackets),
		DPDKRxTestPacketsKey:               fmt.Sprintf("%d", checkupStatus.Results.DPDKRxTestPackets),
		DPDKRxDropsKey:                     fmt.Sprintf("%d", checkupStatus.Results.DPDKPacketsRxDropped),
		DPDKTxDropsKey:                     fmt.Sprintf("%d", checkupStatus.Results.DPDKPacketsTxDropped),
		TrafficGeneratorNodeKey:            checkupStatus.Results.TrafficGeneratorNode,
		DPDKVMNodeKey:                      checkupStatus.Results.DPDKVMNode,
	}

	return formattedResults
}
