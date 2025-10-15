.PHONY: all build test clean examples fmt vet test-quick test-gates build-rust test-rust test-nocache test-gates-nocache test-rust-nocache test-gates-rust-nocache

all: build test

build:
	@echo "Building go-tfhe (pure Go)..."
	go build ./...

build-rust:
	@echo "Building Rust FFT bridge..."
	cd fft-bridge && cargo build --release
	@echo "Building go-tfhe with Rust FFT..."
	go build -tags rust ./...

test:
	@echo "Running tests (pure Go)..."
	go test -v ./...

test-rust:
	@echo "Running tests with Rust FFT..."
	go test -tags rust -v ./...

test-quick:
	@echo "Running quick tests (non-gate tests)..."
	go test -v ./params ./utils ./bitutils ./tlwe ./trlwe ./key ./cloudkey ./fft

test-gates:
	@echo "Running gate tests (this will take several minutes)..."
	@echo "Each gate test takes ~400ms, batch tests take longer..."
	go test -v -timeout 30m ./gates

test-gates-rust:
	@echo "Running gate tests with Rust FFT (should be 4-5x faster)..."
	go test -tags rust -v -timeout 10m ./gates

test-nocache:
	@echo "Running tests without cache (pure Go)..."
	go test -count=1 -v ./...

test-rust-nocache:
	@echo "Running tests without cache (Rust FFT)..."
	go test -count=1 -tags rust -v ./...

test-gates-nocache:
	@echo "Running gate tests without cache (pure Go)..."
	go test -count=1 -v -timeout 30m ./gates

test-gates-rust-nocache:
	@echo "Running gate tests without cache (Rust FFT)..."
	go test -count=1 -tags rust -v -timeout 10m ./gates

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
	@echo "Cleaning Rust FFT bridge..."
	cd fft-bridge && cargo clean

install-deps:
	@echo "Installing dependencies..."
	go mod download
	go mod verify

benchmark:
	@echo "Running FFT benchmarks..."
	go test -bench=. -benchmem ./fft

benchmark-rust:
	@echo "Running FFT benchmarks with Rust backend..."
	go test -tags rust -bench=. -benchmem ./fft

help:
	@echo "Available targets:"
	@echo ""
	@echo "Building:"
	@echo "  all                      - Build and test (pure Go)"
	@echo "  build                    - Build all packages (pure Go)"
	@echo "  build-rust               - Build with Rust FFT backend"
	@echo ""
	@echo "Testing:"
	@echo "  test                     - Run all tests (pure Go)"
	@echo "  test-rust                - Run all tests with Rust FFT"
	@echo "  test-nocache             - Run all tests without cache (pure Go)"
	@echo "  test-rust-nocache        - Run all tests without cache (Rust FFT)"
	@echo "  test-quick               - Run quick tests (no gate tests)"
	@echo "  test-gates               - Run gate tests only (pure Go, slow)"
	@echo "  test-gates-rust          - Run gate tests with Rust FFT (4-5x faster)"
	@echo "  test-gates-nocache       - Run gate tests without cache (pure Go)"
	@echo "  test-gates-rust-nocache  - Run gate tests without cache (Rust FFT)"
	@echo ""
	@echo "Benchmarking:"
	@echo "  benchmark                - Benchmark FFT (pure Go)"
	@echo "  benchmark-rust           - Benchmark FFT (Rust backend)"
	@echo ""
	@echo "Examples:"
	@echo "  examples                 - Build all examples"
	@echo "  run-add                  - Run add_two_numbers example"
	@echo "  run-gates                - Run simple_gates example"
	@echo ""
	@echo "Utilities:"
	@echo "  fmt                      - Format code"
	@echo "  vet                      - Run go vet"
	@echo "  clean                    - Remove build artifacts"
	@echo "  install-deps             - Install/verify dependencies"
