#
# Copyright 2024 Terminal Stream Inc.
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

WEASEL="licenseweasel/weasel:v0.4"
GOLANGCILINT="golangci/golangci-lint:v1.61.0"

lint:
	@echo "Running linters..."
	@docker run --rm -v $(shell pwd):/app -w /app ${GOLANGCILINT} golangci-lint run

license:
	@echo "Verifying license headers..."
	@docker run --rm -v $(shell pwd):/app -w /app ${WEASEL} weasel .

checks: lint license

