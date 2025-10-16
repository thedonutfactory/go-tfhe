.PHONY: all build test clean examples fmt vet test-quick test-gates test-nocache test-gates-nocache

all: build test

build:
	@echo "Building go-tfhe..."
	go build ./...

test:
	@echo "Running tests..."
	go test -v ./...

test-quick:
	@echo "Running quick tests (non-gate tests)..."
	go test -v ./params ./utils ./bitutils ./tlwe ./trlwe ./key ./cloudkey ./evaluator

test-gates:
	@echo "Running gate tests (this will take several minutes)..."
	@echo "Each gate test takes ~400ms, batch tests take longer..."
	go test -v -timeout 30m ./gates

test-nocache:
	@echo "Running tests without cache..."
	go test -count=1 -v ./...

test-gates-nocache:
	@echo "Running gate tests without cache..."
	go test -count=1 -v -timeout 30m ./gates

examples:
	@echo "Building examples..."
	cd examples/add_two_numbers && go build -o ../../bin/add_two_numbers
	cd examples/simple_gates && go build -o ../../bin/simple_gates
	cd examples/programmable_bootstrap && go build -o ../../bin/programmable_bootstrap

run-add:
	@echo "Running add_two_numbers example..."
	cd examples/add_two_numbers && go run main.go

run-gates:
	@echo "Running simple_gates example..."
	cd examples/simple_gates && go run main.go

run-pbs:
	@echo "Running programmable_bootstrap example..."
	cd examples/programmable_bootstrap && go run main.go

fmt:
	@echo "Formatting code..."
	go fmt ./...

vet:
	@echo "Running go vet..."
	go vet ./...

clean:
	@echo "Cleaning build artifacts..."
	go clean ./...
	rm -rf bin/
	rm -f examples/add_two_numbers/add_two_numbers
	rm -f examples/simple_gates/simple_gates
	rm -f examples/programmable_bootstrap/programmable_bootstrap

install-deps:
	@echo "Installing dependencies..."
	go mod download
	go mod verify

benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

help:
	@echo "Available targets:"
	@echo ""
	@echo "Building:"
	@echo "  all                      - Build and test"
	@echo "  build                    - Build all packages"
	@echo ""
	@echo "Testing:"
	@echo "  test                     - Run all tests"
	@echo "  test-nocache             - Run all tests without cache"
	@echo "  test-quick               - Run quick tests (no gate tests)"
	@echo "  test-gates               - Run gate tests only"
	@echo "  test-gates-nocache       - Run gate tests without cache"
	@echo ""
	@echo "Benchmarking:"
	@echo "  benchmark                - Benchmark FFT"
	@echo ""
	@echo "Examples:"
	@echo "  examples                 - Build all examples"
	@echo "  run-add                  - Run add_two_numbers example (8-bit ripple-carry)"
	@echo "  run-gates                - Run simple_gates example"
	@echo "  run-pbs                  - Run programmable_bootstrap example"
	@echo ""
	@echo "Utilities:"
	@echo "  fmt                      - Format code"
	@echo "  vet                      - Run go vet"
	@echo "  clean                    - Remove build artifacts"
	@echo "  install-deps             - Install/verify dependencies"
