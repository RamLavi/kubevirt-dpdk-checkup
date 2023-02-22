#!/usr/bin/env bash
#
# This file is part of the kiagnose project
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# Copyright 2023 Red Hat, Inc.
#

set -eu

if [ "${SET_VERBOSE}" == "TRUE" ]; then
	set -x
fi

print_params() {
	echo NUM_OF_TRAFFIC_CPUS="${NUM_OF_TRAFFIC_CPUS}"
}

print_params

./t-rex-64 --no-ofed-check --no-hw-flow-stat -i -v 3 -c "${NUM_OF_TRAFFIC_CPUS}" --iom 0
