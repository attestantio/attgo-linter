# Copyright © 2026 Attestant Limited.
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

CUSTOM_GCL := ./custom-gcl-attgo

.PHONY: all
all: lint test

# Build custom golangci-lint binary with attgo-linter plugin
$(CUSTOM_GCL): .custom-gcl.yml go.mod go.sum $(shell find . -name '*.go' -not -path './testdata/*')
	golangci-lint custom

.PHONY: build
build: $(CUSTOM_GCL)

.PHONY: lint
lint: $(CUSTOM_GCL)
	$(CUSTOM_GCL) run
	rm $(CUSTOM_GCL)

.PHONY: lint-fix
lint-fix: $(CUSTOM_GCL)
	$(CUSTOM_GCL) run --fix
	rm $(CUSTOM_GCL)

.PHONY: test
test:
	go test ./...

.PHONY: test-race
test-race:
	go test -race ./...

.PHONY: test-verbose
test-verbose:
	go test -v ./...

.PHONY: clean
clean:
	rm -f $(CUSTOM_GCL)
	go clean -cache -testcache

.PHONY: clean-lint
clean-lint:
	rm -f $(CUSTOM_GCL)
	rm -rf ~/.cache/golangci-lint

.PHONY: help
help:
	@echo "attgo-linter Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make              Build and run lint + tests"
	@echo "  make build        Build custom golangci-lint binary"
	@echo "  make lint         Run linter"
	@echo "  make lint-fix     Run linter with auto-fix"
	@echo "  make test         Run tests"
	@echo "  make test-race    Run tests with race detector"
	@echo "  make clean        Remove build artifacts"
	@echo "  make clean-lint   Remove custom binary and lint cache"
	@echo "  make help         Show this help"
