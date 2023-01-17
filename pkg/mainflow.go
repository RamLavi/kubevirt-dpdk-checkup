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

package pkg

import (
	"log"

	kconfig "github.com/kiagnose/kiagnose/kiagnose/config"

	"github.com/kiagnose/kubevirt-dpdk-checkup/pkg/internal/client"
	"github.com/kiagnose/kubevirt-dpdk-checkup/pkg/internal/config"
)

func Run(rawEnv map[string]string) error {
	c, err := client.New()
	if err != nil {
		return err
	}

	baseConfig, err := kconfig.Read(c, rawEnv)
	if err != nil {
		return err
	}

	cfg, err := config.New(baseConfig)
	if err != nil {
		return err
	}

	printConfig(cfg)

	return nil
}

func printConfig(checkupConfig config.Config) {
	log.Println("Using the following config:")
	log.Printf("\t%q: %q", config.NUMASocketParamName, checkupConfig.NUMASocket)
	log.Printf("\t%q: %q", config.NetworkAttachmentDefinitionNameParamName, checkupConfig.NetworkAttachmentDefinitionName)
	log.Printf("\t%q: %q", config.PortBandwidthGBParamName, checkupConfig.PortBandwidthGB)
	log.Printf("\t%q: %q", config.TrafficGeneratorNodeLabelSelectorParamName, checkupConfig.TrafficGeneratorNodeLabelSelector)
	log.Printf("\t%q: %q", config.TrafficGeneratorPacketsPerSecondInMillionsParamName,
		checkupConfig.TrafficGeneratorPacketsPerSecondInMillions)
	log.Printf("\t%q: %q", config.DPDKNodeLabelSelectorParamName, checkupConfig.DPDKNodeLabelSelector)
	log.Printf("\t%q: %q", config.TrafficGeneratorMacAddressParamName, checkupConfig.TrafficGeneratorMacAddress)
	log.Printf("\t%q: %q", config.DPDKMacAddressParamName, checkupConfig.DPDKMacAddress)
	log.Printf("\t%q: %q", config.TestDurationParamName, checkupConfig.TestDuration)
}
