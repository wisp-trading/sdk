.PHONY: generate-proto
generate-proto:
	@echo "Generating gRPC code from proto files..."
	@mkdir -p grpc/gen/inference
	protoc --go_out=. --go_opt=module=github.com/wisp-trading/sdk \
		--go-grpc_out=. --go-grpc_opt=module=github.com/wisp-trading/sdk \
		grpc/proto/inference.proto
	@echo "Generated gRPC code successfully"

.PHONY: generate-mocks
generate-mocks:
	@echo "Generating mocks..."
	mockery
	@echo "Generated mocks successfully"

# Redundancy / dead-code sweep (structural + deadcode from product surface).
# Keeps the graph lean for AAA market DX: clone a domain shell, ship a bot fast.
.PHONY: redundancy
redundancy:
	@go run ./tools/redundancy

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  generate-proto  - Generate Go code from proto files"
	@echo "  generate-mocks  - Generate mocks from interfaces (mockery)"
	@echo "  redundancy      - Structural zombies + deadcode report"
