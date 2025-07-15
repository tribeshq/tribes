-include .env

START_LOG = @echo "======================= START OF LOG ======================="
END_LOG   = @echo "======================== END OF LOG ======================="

# Generic forge script runner
define FORGE_SCRIPT
	$(START_LOG)
	@mkdir -p ./contracts/deployments
	@forge clean --root ./contracts
	@forge script $(1) \
		--root ./contracts \
		--rpc-url $(BLOCKCHAIN_HTTP_ENDPOINT) \
		--private-key defaultKey \
		--broadcast \
		-vvv
	$(END_LOG)
endef

# Generic forge script runner (simulation only)
define FORGE_SCRIPT_SIMULATE
	$(START_LOG)
	@mkdir -p ./contracts/deployments
	@forge clean --root ./contracts
	@forge script $(1) \
		--root ./contracts \
		-vvv
	$(END_LOG)
endef

define deploy_deployer
	$(call FORGE_SCRIPT,./contracts/script/DeployDeployer.s.sol:DeployDeployer)
endef

define deploy_tokens
	$(call FORGE_SCRIPT,./contracts/script/DeployTokens.s.sol:DeployTokens)
endef

define deploy_vlayer
	$(call FORGE_SCRIPT,./contracts/script/DeployVLayer.s.sol:DeployVLayer)
endef

define deploy_emergency
	$(call FORGE_SCRIPT,./contracts/script/DeployEmergency.s.sol:DeployEmergency)
endef

define deploy_delegatecall
	$(call FORGE_SCRIPT,./contracts/script/DeployDelegatecall.s.sol:DeployDelegatecall)
endef

define deploy_safe_call
	$(call FORGE_SCRIPT,./contracts/script/DeploySafeCall.s.sol:DeploySafeCall)
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
test: test-integration test-contracts ## Run all tests

.PHONY: test-integration
test-integration: ## Run the integration tests
	$(START_LOG)
	@go generate ./...
	@go test -p 1 ./test/integration/... -coverprofile=./coverage.md -v
	$(END_LOG)

.PHONY: test-contracts
test-contracts: ## Run the contracts tests
	$(START_LOG)
	@forge clean --root ./contracts
	@forge test --root ./contracts -vvvv
	$(END_LOG)

.PHONY: lint
lint: ## Run code linting and formatting checks
	$(START_LOG)
	@go vet ./...
	@forge fmt --root ./contracts
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

# =============================================================================
# DEPLOYMENT COMMANDS
# =============================================================================

.PHONY: contracts
contracts: deploy-deployer deploy-tokens deploy-vlayer deploy-emergency deploy-safe-call ## Deploy all contracts using modular deployment scripts
	$(START_LOG)
	@echo "All contracts deployed! Check ./deployments/ for individual deployment files"
	$(END_LOG)

.PHONY: deploy
deploy: contracts ## Alias for contracts deployment

.PHONY: deploy-simulate
deploy-simulate: ## Simulate deployment without broadcasting
	@$(call FORGE_SCRIPT_SIMULATE,./contracts/script/DeployDeployer.s.sol:DeployDeployer)
	@$(call FORGE_SCRIPT_SIMULATE,./contracts/script/DeployTokens.s.sol:DeployTokens)
	@$(call FORGE_SCRIPT_SIMULATE,./contracts/script/DeployVLayer.s.sol:DeployVLayer)
	@$(call FORGE_SCRIPT_SIMULATE,./contracts/script/DeployEmergency.s.sol:DeployEmergency)
	@$(call FORGE_SCRIPT_SIMULATE,./contracts/script/DeploySafeCall.s.sol:DeploySafeCall)

# =============================================================================
# MODULAR DEPLOYMENT COMMANDS
# =============================================================================

.PHONY: deploy-deployer
deploy-deployer: ## Deploy only Deployer contract
	$(START_LOG)
	@$(deploy_deployer)
	@echo "Deployer deployment completed! Check ./deployments/deployer.json for addresses"
	$(END_LOG)

.PHONY: deploy-tokens
deploy-tokens: ## Deploy only tokens (Collateral, Stablecoin)
	$(START_LOG)
	@$(deploy_tokens)
	@echo "Tokens deployment completed! Check ./deployments/tokens.json for addresses"
	$(END_LOG)

.PHONY: deploy-vlayer
deploy-vlayer: ## Deploy only VLayer contracts (Prover, Verifier)
	$(START_LOG)
	@$(deploy_vlayer)
	@echo "VLayer deployment completed! Check ./deployments/vlayer.json for addresses"
	$(END_LOG)

.PHONY: deploy-emergency
deploy-emergency: ## Deploy only emergency contracts (EmergencyWithdraw)
	$(START_LOG)
	@$(deploy_emergency)
	@echo "Emergency deployment completed! Check ./deployments/emergency.json for addresses"
	$(END_LOG)

.PHONY: deploy-delegatecall
deploy-delegatecall: ## Deploy only delegatecall contracts (EmergencyWithdraw)
	$(START_LOG)
	@$(deploy_delegatecall)
	@echo "Delegatecall deployment completed! Check ./deployments/emergency.json for addresses"
	$(END_LOG)

.PHONY: deploy-safe-call
deploy-safe-call: ## Deploy only SafeCall contract
	$(START_LOG)
	@$(deploy_safe_call)
	@echo "SafeCall deployment completed! Check ./deployments/safeCall.json for addresses"
	$(END_LOG)



# =============================================================================
# UTILITY COMMANDS
# =============================================================================

.PHONY: clean
clean: ## Clean build artifacts
	$(START_LOG)
	@forge clean --root ./contracts
	@rm -rf ./deployments
	@echo "Clean completed"
	$(END_LOG)

.PHONY: size
size: ## Check contract sizes
	$(START_LOG)
	@forge build --root ./contracts --sizes
	$(END_LOG)

.PHONY: gas
gas: ## Run gas reports
	$(START_LOG)
	@forge test --root ./contracts --gas-report
	$(END_LOG)

.PHONY: help
help: ## Show help for each of the Makefile recipes
	@echo "Available commands:"
	@awk '/^[a-zA-Z0-9_-]+:.*##/ { \
		split($$0, parts, "##"); \
		split(parts[1], target, ":"); \
		printf "  \033[36m%-30s\033[0m %s\n", target[1], parts[2] \
	}' $(MAKEFILE_LIST)