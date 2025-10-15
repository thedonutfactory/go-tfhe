#!/bin/bash
# Build script for the Rust FFT bridge library

set -e

echo "Building Rust FFT bridge..."
cd "$(dirname "$0")"

# Build release version
cargo build --release

echo "✅ Rust FFT bridge built successfully!"
echo "📦 Library location: target/release/libtfhe_fft_bridge.{a,dylib}"
echo ""
echo "To use in Go:"
echo "  go build -tags rust ./..."
echo ""
echo "To test:"
echo "  cargo test"


