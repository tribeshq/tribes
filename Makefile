-include .env.develop

START_LOG = @echo "================================================= START OF LOG ==================================================="
END_LOG = @echo "================================================== END OF LOG ===================================================="

.PHONY: verifier
verifier:
	$(START_LOG)
	@cd ./tools/tlsnotary/verifier && cargo build --release
	@cp ./tools/tlsnotary/verifier/target/release/libverifier.a ./internal/usecase/crowdfunding/
	$(END_LOG)
	
.PHONY: generate
generate:
	$(START_LOG)
	@go generate ./...
	$(END_LOG)

.PHONY: test
test:
	@cd ./tools/tlsnotary/verifier && cargo build --release
	@cp ./tools/tlsnotary/verifier/target/release/libverifier.a ./internal/usecase/crowdfunding/
	@go test -p=1 ./... -coverprofile=./coverage.md -v

.PHONY: coverage
coverage: test
	@go tool cover -html=./coverage.md

.PHONY: state
state:
	@chmod +x ./tools/state.sh
	@./tools/state.sh

.PHONY: state-dev
state-dev:
	@chmod +x ./tools/state-dev.sh
	@./tools/state-dev.sh

.PHONY: state-mcp
state-mcp:
	@chmod +x ./tools/state-mcp.sh
	@./tools/state-mcp.sh $(DAPP_ADDRESS)
