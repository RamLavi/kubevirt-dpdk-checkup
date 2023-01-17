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

package config_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	assert "github.com/stretchr/testify/require"

	kconfig "github.com/kiagnose/kiagnose/kiagnose/config"

	"github.com/kiagnose/kubevirt-dpdk-checkup/pkg/internal/config"
)

const (
	testPodName                                = "my-pod"
	testPodUID                                 = "0123456789-0123456789"
	numaSocket                                 = 1
	networkAttachmentDefinitionName            = "intel-dpdk-network1"
	portBandwidthGB                            = 100
	trafficGeneratorNodeLabelSelector          = "node-role.kubernetes.io/worker-dpdk1"
	trafficGeneratorPacketsPerSecondInMillions = 6
	dpdkNodeLabelSelector                      = "node-role.kubernetes.io/worker-dpdk2"
	trafficGeneratorEastMacAddress             = "DE:AD:BE:EF:00:01"
	dpdkMacAddress                             = "DE:AD:BE:EF:00:02"
	testDuration                               = "30m"
)

func TestNewShouldApplyDefaultsWhenOptionalFieldsAreMissing(t *testing.T) {
	baseConfig := kconfig.Config{
		PodName: testPodName,
		PodUID:  testPodUID,
		Params: map[string]string{
			config.NUMASocketParamName:                      fmt.Sprintf("%d", numaSocket),
			config.NetworkAttachmentDefinitionNameParamName: networkAttachmentDefinitionName,
		},
	}

	actualConfig, err := config.New(baseConfig)
	assert.NoError(t, err)

	trafficGeneratorEastMacAddressDefault, _ := net.ParseMAC(config.TrafficGeneratorEastMacAddressDefault)
	dpdkMacAddressDefault, _ := net.ParseMAC(config.DPDKMacAddressDefault)
	expectedConfig := config.Config{
		PodName:                         testPodName,
		PodUID:                          testPodUID,
		NUMASocket:                      numaSocket,
		NetworkAttachmentDefinitionName: networkAttachmentDefinitionName,
		TrafficGeneratorPacketsPerSecondInMillions: config.TrafficGeneratorPacketsPerSecondInMillionsDefault,
		PortBandwidthGB:                config.PortBandwidthGBDefault,
		TrafficGeneratorEastMacAddress: trafficGeneratorEastMacAddressDefault,
		DPDKMacAddress:                 dpdkMacAddressDefault,
		TestDuration:                   config.TestDurationDefault,
	}
	assert.Equal(t, expectedConfig, actualConfig)
}

func TestNewShouldApplyUserConfig(t *testing.T) {
	baseConfig := kconfig.Config{
		PodName: testPodName,
		PodUID:  testPodUID,
		Params:  getValidUserParameters(),
	}

	actualConfig, err := config.New(baseConfig)
	assert.NoError(t, err)

	trafficGeneratorEastHWAddress, _ := net.ParseMAC(trafficGeneratorEastMacAddress)
	dpdkHWAddress, _ := net.ParseMAC(dpdkMacAddress)
	expectedConfig := config.Config{
		PodName:                         testPodName,
		PodUID:                          testPodUID,
		NUMASocket:                      numaSocket,
		PortBandwidthGB:                 portBandwidthGB,
		NetworkAttachmentDefinitionName: networkAttachmentDefinitionName,
		TrafficGeneratorPacketsPerSecondInMillions: trafficGeneratorPacketsPerSecondInMillions,
		TrafficGeneratorNodeLabelSelector:          trafficGeneratorNodeLabelSelector,
		DPDKNodeLabelSelector:                      dpdkNodeLabelSelector,
		TrafficGeneratorEastMacAddress:             trafficGeneratorEastHWAddress,
		DPDKMacAddress:                             dpdkHWAddress,
		TestDuration:                               30 * time.Minute,
	}
	assert.Equal(t, expectedConfig, actualConfig)
}

func TestNewShouldFailWhen(t *testing.T) {
	type failureTestCase struct {
		description    string
		key            string
		faultyKeyValue string
		expectedError  error
	}

	testCases := []failureTestCase{
		{
			description:    "NUMA Socket is empty string",
			key:            config.NUMASocketParamName,
			faultyKeyValue: "",
			expectedError:  config.ErrInvalidNUMASocket,
		},
		{
			description:    "NUMA Socket is invalid",
			key:            config.NUMASocketParamName,
			faultyKeyValue: "-1",
			expectedError:  config.ErrInvalidNUMASocket,
		},
		{
			description:    "NetworkAttachmentDefinitionName is invalid",
			key:            config.NetworkAttachmentDefinitionNameParamName,
			faultyKeyValue: "",
			expectedError:  config.ErrInvalidNetworkAttachmentDefinitionName,
		},
		{
			description:    "TrafficGeneratorPacketsPerSecondInMillions is invalid",
			key:            config.TrafficGeneratorPacketsPerSecondInMillionsParamName,
			faultyKeyValue: "-14",
			expectedError:  config.ErrInvalidTrafficGeneratorPacketsPerSecondInMillions,
		},
		{
			description:    "PortBandwidthGB is invalid",
			key:            config.PortBandwidthGBParamName,
			faultyKeyValue: "0",
			expectedError:  config.ErrInvalidPortBandwidthGB,
		},
		{
			description:    "TrafficGeneratorEastMacAddress is invalid",
			key:            config.TrafficGeneratorEastMacAddressParamName,
			faultyKeyValue: "AB:CD:EF:GH:IJ:KH",
			expectedError:  config.ErrInvalidTrafficGeneratorEastMacAddress,
		},
		{
			description:    "DPDKMacAddress is invalid",
			key:            config.DPDKMacAddressParamName,
			faultyKeyValue: "AB:CD:EF:GH:IJ:KH",
			expectedError:  config.ErrInvalidDPDKMacAddress,
		},
		{
			description:    "TestDuration is invalid",
			key:            config.TestDurationParamName,
			faultyKeyValue: "invalid value",
			expectedError:  config.ErrInvalidTestDuration,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			faultyUserParams := getValidUserParameters()
			faultyUserParams[testCase.key] = testCase.faultyKeyValue

			baseConfig := kconfig.Config{
				PodName: testPodName,
				PodUID:  testPodUID,
				Params:  faultyUserParams,
			}

			_, err := config.New(baseConfig)
			assert.ErrorIs(t, err, testCase.expectedError)
		})
	}
}

func getValidUserParameters() map[string]string {
	return map[string]string{
		config.NUMASocketParamName:                                 fmt.Sprintf("%d", numaSocket),
		config.NetworkAttachmentDefinitionNameParamName:            networkAttachmentDefinitionName,
		config.PortBandwidthGBParamName:                            fmt.Sprintf("%d", portBandwidthGB),
		config.TrafficGeneratorNodeLabelSelectorParamName:          trafficGeneratorNodeLabelSelector,
		config.TrafficGeneratorPacketsPerSecondInMillionsParamName: fmt.Sprintf("%d", trafficGeneratorPacketsPerSecondInMillions),
		config.DPDKNodeLabelSelectorParamName:                      dpdkNodeLabelSelector,
		config.TrafficGeneratorEastMacAddressParamName:             trafficGeneratorEastMacAddress,
		config.DPDKMacAddressParamName:                             dpdkMacAddress,
		config.TestDurationParamName:                               testDuration,
	}
}
