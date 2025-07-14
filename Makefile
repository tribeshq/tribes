-include .env

START_LOG = @echo "======================= START OF LOG ======================="
END_LOG   = @echo "======================== END OF LOG ======================="

# Generic forge script runner
define FORGE_SCRIPT
	$(START_LOG)
	@forge clean --root ./contracts
	@forge script $(1) \
		--root ./contracts \
		--rpc-url $(BLOCKCHAIN_HTTP_ENDPOINT) \
		--private-key defaultKey \
		--broadcast \
		-vvv
	$(END_LOG)
endef

define deploy_creat_deployer_proxy
	$(call FORGE_SCRIPT,./contracts/script/DeployCREATDeployerProxy.s.sol)
endef

define deploy_creat2_deployer_proxy
	$(call FORGE_SCRIPT,./contracts/script/DeployCREAT2DeployerProxy.s.sol)
endef

define deploy_assets
	$(call FORGE_SCRIPT,./contracts/script/DeployAssets.s.sol)
endef

define deploy_badge
	$(call FORGE_SCRIPT,./contracts/script/DeployBadge.s.sol)
endef

define deploy_vlayer
	$(call FORGE_SCRIPT,./contracts/script/DeployVlayer.s.sol)
endef

define deploy_delegatecall
	$(call FORGE_SCRIPT,./contracts/script/DeployDelegatecall.s.sol)
endef

.PHONY: env
env: ## Create the environment variables file
	$(START_LOG)
	@cp .env.tmpl .env
	@echo "Environment variables file created"
	$(END_LOG)

.PHONY: generate
generate: ## Generate bytecode and Go bindings
	$(START_LOG)
	@forge clean --root ./contracts
	@forge script --root ./contracts ./contracts/script/GenerateBytecode.s.sol:GenerateBytecode
	@go generate ./...
	$(END_LOG)

.PHONY: test
test: ## Run the application tests (Contracts + Backend)
	$(START_LOG)
	@forge clean --root ./contracts
	@forge test --root ./contracts
	@go generate ./...
	@go test ./... -coverprofile=./coverage.md -v
	$(END_LOG)

.PHONY: lint
lint: ## Run code linting and formatting checks
	$(START_LOG)
	@go vet ./...
	@forge fmt --check --root ./contracts
	@echo "Linting completed"
	$(END_LOG)

.PHONY: fmt
fmt: ## Format all code (Contracts + Backend)
	$(START_LOG)
	@gofmt -w .
	@forge fmt --root ./contracts
	@echo "Formatting completed"
	$(END_LOG)

.PHONY: coverage
coverage: ## Open HTML coverage report
	$(START_LOG)
	@go tool cover -html=./coverage.md
	@echo "Coverage report opened"
	$(END_LOG)

.PHONY: contracts
contracts: ## Deploy all contracts
	@$(deploy_assets)
	@$(deploy_badge)
	@$(deploy_vlayer)
	@$(deploy_delegatecall)
	@$(deploy_creat_deployer_proxy)
	@$(deploy_creat2_deployer_proxy)

.PHONY: deploy-creat-deployer-proxy
deploy-creat-deployer-proxy: ## Deploy CREAT deployer proxy
	@$(deploy_creat_deployer_proxy)

.PHONY: deploy-creat2-deployer-proxy
deploy-creat2-deployer-proxy: ## Deploy CREAT2 deployer proxy
	@$(deploy_creat2_deployer_proxy)

.PHONY: deploy-badge
deploy-badge: ## Deploy Badge contract
	@$(deploy_badge)

.PHONY: deploy-vlayer
deploy-vlayer: ## Deploy Vlayer contract
	@$(deploy_vlayer)

.PHONY: deploy-delegatecall
deploy-delegatecall: ## Deploy Delegatecall contract
	@$(deploy_delegatecall)

.PHONY: deploy-assets
deploy-assets: ## Deploy Assets contract
	@$(deploy_assets)

.PHONY: help
help: ## Show help for each of the Makefile recipes
	@echo "Available commands:"
	@awk '/^[a-zA-Z0-9_-]+:.*##/ { \
		split($$0, parts, "##"); \
		split(parts[1], target, ":"); \
		printf "  \033[36m%-30s\033[0m %s\n", target[1], parts[2] \
	}' $(MAKEFILE_LIST)