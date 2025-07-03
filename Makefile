include .env.tmpl

START_LOG = @echo "======================= START OF LOG ======================="
END_LOG = @echo "======================== END OF LOG ======================="

.PHONY: env
env: ## Create the environment variables file
	@cp .env.tmpl .env

.PHONY: build
build: ## Build the application RISC-V image with cartesi cli
	$(START_LOG)
	@cartesi build
	$(END_LOG)
	
.PHONY: generate
generate: ## Generate the application code
	$(START_LOG)
	@go generate ./...
	$(END_LOG)

.PHONY: test
test: ## Run the application tests
	$(START_LOG)
	@forge test -vvv --root ./contracts
	@go test -p=1 ./... -coverprofile=./coverage.md -v
	$(END_LOG)

.PHONY: test-sol
test-sol: ## Run only Solidity tests
	$(START_LOG)
	@forge test -vvv --root ./contracts
	$(END_LOG)

.PHONY: test-go
test-go: ## Run only Go tests
	$(START_LOG)
	@go test -p=1 ./... -coverprofile=./coverage.md -v
	$(END_LOG)

.PHONY: lint
lint: ## Run linting and formatting checks
	$(START_LOG)
	@test -z "$(gofmt -l .)" || (echo "Go code is not formatted. Run 'gofmt -w .'" && exit 1)
	@go vet ./...
	@forge fmt --check --root ./contracts
	$(END_LOG)

.PHONY: fmt
fmt: ## Format code
	$(START_LOG)
	@gofmt -w .
	@forge fmt --root ./contracts
	$(END_LOG)

.PHONY: coverage
coverage: test ## Generate the application code coverage report
	$(START_LOG)
	@go tool cover -html=./coverage.md
	$(END_LOG)

.PHONY: state
state: ## Run the application state for devnet (demo)
	@chmod +x ./tools/state.sh
	@./tools/state.sh

.PHONY: help
help: ## Show help for each of the Makefile recipes
	@grep "##" $(MAKEFILE_LIST) | grep -v grep | sed -e 's/:.*##/:\t/'