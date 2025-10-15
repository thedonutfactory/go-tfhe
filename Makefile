.PHONY: all build test clean examples fmt vet test-quick test-gates

all: build test

build:
	@echo "Building go-tfhe..."
	go build ./...

test:
	@echo "Running tests..."
	go test -v ./...

test-quick:
	@echo "Running quick tests (non-gate tests)..."
	go test -v ./params ./utils ./bitutils ./tlwe ./trlwe ./key ./cloudkey ./fft

test-gates:
	@echo "Running gate tests (this will take several minutes)..."
	@echo "Each gate test takes ~400ms, batch tests take longer..."
	go test -v -timeout 30m ./gates

examples:
	@echo "Building examples..."
	cd examples/add_two_numbers && go build
	cd examples/simple_gates && go build

run-add:
	@echo "Running add_two_numbers example..."
	cd examples/add_two_numbers && go run main.go

run-gates:
	@echo "Running simple_gates example..."
	cd examples/simple_gates && go run main.go

fmt:
	@echo "Formatting code..."
	go fmt ./...

vet:
	@echo "Running go vet..."
	go vet ./...

clean:
	@echo "Cleaning build artifacts..."
	go clean ./...
	rm -f examples/add_two_numbers/add_two_numbers
	rm -f examples/simple_gates/simple_gates

install-deps:
	@echo "Installing dependencies..."
	go mod download
	go mod verify

help:
	@echo "Available targets:"
	@echo "  all          - Build and test"
	@echo "  build        - Build all packages"
	@echo "  test         - Run all tests"
	@echo "  test-quick   - Run quick tests (no gate tests)"
	@echo "  test-gates   - Run gate tests only (slow)"
	@echo "  examples     - Build all examples"
	@echo "  run-add      - Run add_two_numbers example"
	@echo "  run-gates    - Run simple_gates example"
	@echo "  fmt          - Format code"
	@echo "  vet          - Run go vet"
	@echo "  clean        - Remove build artifacts"
	@echo "  install-deps - Install/verify dependencies"
