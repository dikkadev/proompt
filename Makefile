# Proompt Project Makefile
.PHONY: help build test test-integration test-smoke clean docker-build docker-test docker-clean

# Default target
help:
	@echo "Proompt Project Commands:"
	@echo ""
	@echo "  build              Build the server binary"
	@echo "  test               Run all Go unit tests"
	@echo "  test-integration   Run integration tests in Docker"
	@echo "  test-smoke         Run smoke tests only"
	@echo "  clean              Clean build artifacts and test data"
	@echo ""
	@echo "  docker-build       Build Docker image"
	@echo "  docker-test        Run containerized integration tests"
	@echo "  docker-clean       Clean Docker resources"
	@echo ""

# Build the server binary
build:
	@echo "Building server..."
	cd server && go build -buildvcs=false -o proompt ./cmd/proompt

# Run Go unit tests
test:
	@echo "Running unit tests..."
	cd server && go test ./...

# Run integration tests in Docker
test-integration: docker-build
	@echo "Running integration tests..."
	cd tests && docker build -f docker/Dockerfile.test -t proompt:test-runner .
	cd tests && docker compose -f docker/compose.test.yaml up --abort-on-container-exit
	cd tests && docker compose -f docker/compose.test.yaml down -v

# Run smoke tests only
test-smoke: docker-build
	@echo "Running smoke tests..."
	cd tests && docker build -f docker/Dockerfile.test -t proompt:test-runner .
	cd tests && docker tag proompt:latest proompt:test
	cd tests && docker compose -f docker/compose.test.yaml up proompt-server -d
	@sleep 5
	cd tests && docker run --rm --network proompt_test_net \
		-v "$(PWD)/tests/integration:/tests:ro" \
		-e PROOMPT_API_URL=http://proompt-server:8080 \
		proompt:test-runner /tests/smoke_tests.sh
	cd tests && docker compose -f docker/compose.test.yaml down -v

# Clean build artifacts and test data
clean:
	@echo "Cleaning up..."
	cd server && rm -f proompt
	cd tests && ./scripts/cleanup.sh -a

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	cd server && docker build -t proompt:latest .

# Run containerized integration tests
docker-test: docker-build
	@echo "Running Docker integration tests..."
	cd tests && docker compose -f docker/compose.test.yaml up --build --abort-on-container-exit
	cd tests && docker compose -f docker/compose.test.yaml down -v

# Clean Docker resources
docker-clean:
	@echo "Cleaning Docker resources..."
	cd tests && ./scripts/cleanup.sh -a