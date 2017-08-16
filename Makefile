# File managed by pluginsync
# http://www.apache.org/licenses/LICENSE-2.0.txt
#
#
# Copyright 2017 Intel Corporation
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

default: 
	$(MAKE) deps
	$(MAKE) all

deps:
	bash -c "./scripts/deps.sh"
#test:
	# bash -c "./scripts/test.sh"
		# todo: test.sh first run make (to get binary). Know location of this binary, 
		# update location of binary in client_test.go ln 44 exec.Command(...)
		# call $go test -run TestClient

#test-all:
	# $(MAKE) test
all: 
	bash -c "./scripts/build.sh"