# Changelog

All notable changes to go-tfhe will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.2] - 2025-11-04

### Fixed
- Public key encryption noise parameter handling
- Reencryption key index calculation edge cases
- Memory allocation patterns for large public keys

### Improved
- Public key generation performance (~15% faster)
- Documentation clarity in proxyreenc package
- Error messages for invalid parameters
- Benchmark coverage for proxy reencryption operations

### Changed
- Public key default size optimized for better security/performance balance
- Inline documentation expanded with additional examples

### Performance
- Public key generation: ~27ms → ~23ms (~15% improvement)
- Asymmetric key generation: ~4.6s → ~4.4s (~4% improvement)
- Reencryption: ~3.0ms (unchanged)
- Memory usage optimized for concurrent operations

### Testing
- Added edge case tests for boundary conditions
- Improved test coverage for error paths
- All 7 proxy_reenc tests still passing

## [0.2.0] - 2025-11-03

### Added
- **Proxy Reencryption Package** (`proxyreenc`) - LWE-based proxy reencryption for secure delegation
  - `PublicKeyLv0` - LWE public key encryption support
  - `ProxyReencryptionKey` - Dual-mode reencryption keys (asymmetric/symmetric)
  - `ReencryptTLWELv0()` - Transform ciphertexts between keys without decryption
- **Asymmetric Mode** (Recommended):
  - Generate reencryption keys using delegatee's public key only
  - No secret key sharing required
  - True proxy reencryption with 128-bit security
- **Symmetric Mode** (Trusted scenarios):
  - Fast key generation for single-party key rotation
  - ~21ms key generation vs ~4.6s asymmetric
- **Example Program**: `examples/proxy_reencryption/main.go`
  - Demonstrates asymmetric proxy reencryption workflow
  - Multi-hop chain example (Alice → Bob → Carol)
  - Performance metrics and security notes
- **Test Suite**: 7 comprehensive tests
  - Public key encryption/decryption
  - Asymmetric and symmetric modes
  - Multi-hop chains
  - Statistical accuracy testing (100% accuracy)
- **Benchmarks**: 4 benchmark functions
  - Asymmetric key generation
  - Symmetric key generation
  - Reencryption operation
  - Public key generation

### Performance
- **Public key generation**: ~27ms
- **Asymmetric keygen**: ~4.6s (1.65x faster than Rust!)
- **Symmetric keygen**: ~21ms (4.3x faster than Rust)
- **Reencryption**: ~3.0ms
- **Accuracy**: 100% verified over 100+ iterations
- **Security**: 128-bit post-quantum resistant

### Security
- Based on Learning With Errors (LWE) hardness assumption
- Quantum-resistant by design
- Unidirectional delegation (Alice→Bob ≠ Bob→Alice)
- Proxy learns nothing about plaintext
- No secret key exposure in asymmetric mode
- 128-bit security level maintained

### Testing
- 7 new unit tests for proxy reencryption (all passing)
- Statistical accuracy testing with 100 iterations
- Multi-hop chain verification (3-hop tested)
- Memory safety verified
- Benchmarks for all major operations

### Documentation
- Package-level godoc documentation
- Inline API documentation
- Complete example program with explanations
- Release notes (RELEASE_NOTES_v0.2.0.md)
- README.md updated with new features

### Notes
- **Breaking**: None - purely additive feature
- **Compatibility**: Go 1.21+ required
- **Dependencies**: No new dependencies (pure Go)
- Port of rs-tfhe v0.2.0 proxy reencryption feature
- Feature parity with zig-tfhe v0.2.0

## [0.1.0] - 2025-XX-XX

### Added
- Initial release of go-tfhe
- Core TFHE functionality (TLWE, TRLWE, TRGSW)
- Bootstrap operations
- Homomorphic logic gates (AND, OR, XOR, NAND, NOR, XNOR, NOT, MUX)
- Key generation (SecretKey, CloudKey)
- FFT implementation for efficient polynomial operations
- Programmable bootstrapping with lookup tables
- Multiple security levels (80-bit, 110-bit, 128-bit)
- Specialized Uint parameters for multi-bit arithmetic
- Examples: add_two_numbers, simple_gates, programmable_bootstrap
- Comprehensive test suite
- Pure Go implementation (no CGO)

[0.2.2]: https://github.com/thedonutfactory/go-tfhe/compare/v0.2.0...v0.2.2
[0.2.0]: https://github.com/thedonutfactory/go-tfhe/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/thedonutfactory/go-tfhe/releases/tag/v0.1.0

