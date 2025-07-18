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

define deploy_badge_factory
	$(call FORGE_SCRIPT,./contracts/script/DeployBadgeFactory.s.sol:DeployBadgeFactory)
endef

define deploy_tokens
	$(call FORGE_SCRIPT,./contracts/script/DeployERC20.s.sol:DeployERC20)
endef

define deploy_vlayer
	$(call FORGE_SCRIPT,./contracts/script/DeployVLayer.s.sol:DeployVLayer)
endef

define deploy_emergency
	$(call FORGE_SCRIPT,./contracts/script/DeployEmergency.s.sol:DeployEmergency)
endef

define deploy_safe_erc1155_mint
	$(call FORGE_SCRIPT,./contracts/script/DeploySafeERC1155Mint.s.sol:DeploySafeERC1155Mint)
endef

.PHONY: env
env: ## Create the environment variables file
	$(START_LOG)
	@cp .env.tmpl .env
	@echo "Environment variables file created"
	$(END_LOG)

.PHONY: build
build: generate ## Build the contracts
	$(START_LOG)
	@cartesi build
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
	@cd contracts && forge soldeer install
	@forge test --root ./contracts -vvv
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

.PHONY: deploy-contracts
deploy-contracts: deploy-badge-factory deploy-erc20 deploy-vlayer deploy-emergency deploy-safe-erc1155-mint ## Deploy all contracts using modular deployment scripts
	$(START_LOG)
	@echo "All contracts deployed! Check ./deployments/ for individual deployment files"
	$(END_LOG)

.PHONY: deploy-contracts-simulate
deploy-contracts-simulate: ## Simulate deployment without broadcasting
	@$(call FORGE_SCRIPT_SIMULATE,./contracts/script/DeployBadgeFactory.s.sol:DeployBadgeFactory)
	@$(call FORGE_SCRIPT_SIMULATE,./contracts/script/DeployERC20.s.sol:DeployERC20)
	@$(call FORGE_SCRIPT_SIMULATE,./contracts/script/DeployVLayer.s.sol:DeployVLayer)
	@$(call FORGE_SCRIPT_SIMULATE,./contracts/script/DeployEmergency.s.sol:DeployEmergency)
	@$(call FORGE_SCRIPT_SIMULATE,./contracts/script/DeploySafeERC1155Mint.s.sol:DeploySafeERC1155Mint)

# =============================================================================
# MODULAR DEPLOYMENT COMMANDS
# =============================================================================

.PHONY: deploy-badge-factory
deploy-badge-factory: ## Deploy only BadgeFactory contract
	$(START_LOG)
	@$(deploy_badge_factory)
	@echo "BadgeFactory deployment completed! Check ./deployments/deployer.json for addresses"
	$(END_LOG)

.PHONY: deploy-erc20
deploy-erc20: ## Deploy only tokens (Collateral, Stablecoin)
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

.PHONY: deploy-safe-erc1155-mint
deploy-safe-erc1155-mint: ## Deploy only SafeERC1155Mint contract
	$(START_LOG)
	@$(deploy_safe_erc1155_mint)
	@echo "SafeERC1155Mint deployment completed! Check ./deployments/outputSafeCall.json for addresses"
	$(END_LOG)

# =============================================================================
# UTILITY COMMANDS
# =============================================================================

.PHONY: clean
clean: ## Clean build artifacts
	$(START_LOG)
	@forge clean --root ./contracts
	@rm -rf ./contracts/cache
	@rm -rf ./contracts/broadcast
	@rm -rf ./contracts/dependencies
	@rm -rf ./contracts/soldeer.lock
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