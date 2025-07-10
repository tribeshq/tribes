-include .env

START_LOG = @echo "======================= START OF LOG ======================="
END_LOG = @echo "======================== END OF LOG ======================="

define deploy_assets
	$(START_LOG)
	@forge clean --root ./contracts
	@forge script ./contracts/script/DeployAssets.s.sol \
		--root ./contracts \
		--rpc-url $(BLOCKCHAIN_HTTP_ENDPOINT) \
		--private-key $(PRIVATE_KEY) \
		--broadcast \
		-vvv
	$(END_LOG)
endef

define deploy_chainlink
	$(START_LOG)\
	@forge clean --root ./contracts
	@forge script ./contracts/script/CrossChainNFT.s.sol:CrossChainNFTSourceMinter \
		--root ./contracts \
		--rpc-url $(BLOCKCHAIN_HTTP_ENDPOINT) \
		--private-key $(PRIVATE_KEY) \
		--broadcast \
		-vvv
	@forge script ./contracts/script/CrossChainNFT.s.sol:CrossChainNFTDestinationMinter \
		--root ./contracts \
		--rpc-url $(ARBITRUM_SEPOLIA_RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		--broadcast \
		-vvv
	$(END_LOG)
endef

define deploy_vlayer
	$(START_LOG)
	@forge clean --root ./contracts
	@forge script ./contracts/script/DeployVlayer.s.sol \
		--root ./contracts \
		--rpc-url $(BLOCKCHAIN_HTTP_ENDPOINT) \
		--private-key $(PRIVATE_KEY) \
		--broadcast \
		-vvv
	$(END_LOG)
endef

define deploy_delegatecall
	$(START_LOG)
	@forge clean --root ./contracts
	@forge script ./contracts/script/DeployDelegatecall.s.sol \
		--root ./contracts \
		--rpc-url $(BLOCKCHAIN_HTTP_ENDPOINT) \
		--private-key $(PRIVATE_KEY) \
		--broadcast \
		-vvv
	$(END_LOG)
endef

define setup
	$(START_LOG)
	@forge clean --root ./contracts
	@forge script ./contracts/script/CrossChainNFT.s.sol:SetupApplication \
		--root ./contracts \
		--rpc-url $(BLOCKCHAIN_HTTP_ENDPOINT) \
		--private-key $(PRIVATE_KEY) \
		--broadcast \
		-vvv
	$(END_LOG)
endef

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
	@go generate ./...
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

.PHONY: contracts
contracts: assets chainlink vlayer delegatecall ## Deploy the contracts

.PHONY: chainlink
chainlink: ## Deploy the chainlink contracts
	@$(deploy_chainlink)

.PHONY: vlayer
vlayer: ## Deploy the vlayer contracts
	@$(deploy_vlayer)

.PHONY: delegatecall
delegatecall: ## Deploy the delegatecall contracts
	@$(deploy_delegatecall)

.PHONY: assets
assets: ## Deploy the assets contracts
	@$(deploy_assets)

.PHONY: setup
setup: ## Transfers ownership of a deployed SourceMinter contract to the tribes application address
	@$(setup)

.PHONY: help
help: ## Show help for each of the Makefile recipes
	@grep "##" $(MAKEFILE_LIST) | grep -v grep | sed -e 's/:.*##/:\t/'