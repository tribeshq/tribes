-include .env.develop

START_LOG = @echo "================================================= START OF LOG ==================================================="
END_LOG = @echo "================================================== END OF LOG ===================================================="

.PHONY: dev
dev:
	$(START_LOG)
	@nonodo -- air

.PHONY: verifier
verifier:
	$(START_LOG)
	@cd ./tools/tlsnotary/verifier && cargo build --release
	@cp ./tools/tlsnotary/verifier/target/release/libverifier.a ./internal/usecase/crowdfunding_usecase/
	$(END_LOG)
	
.PHONY: generate
generate:
	$(START_LOG)
	@go run ./pkg/rollups-contracts/generate
	$(END_LOG)

.PHONY: test
test:
	@cd ./tools/tlsnotary/verifier && cargo build --release
	@cp ./tools/tlsnotary/verifier/target/release/libverifier.a ./internal/usecase/crowdfunding_usecase/
	@go test -p=1 ./... -coverprofile=./coverage.md -v

.PHONY: coverage
coverage: test
	@go tool cover -html=./coverage.md

.PHONY: state
state:
	@chmod +x ./tools/state.sh
	@./tools/state.sh