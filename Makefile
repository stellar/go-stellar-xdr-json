# update the Cargo.lock every time the Cargo.toml changes.
Cargo.lock: Cargo.toml
	cargo update --workspace

# Define targets and output directories
TARGETS = \
    x86_64-pc-windows-gnu \
    x86_64-apple-darwin \
    aarch64-apple-darwin \
    x86_64-unknown-linux-gnu \
    aarch64-unknown-linux-gnu

LIBS_DIR = xdr2json/libs
BUILD_DIR = target
PROFILE = release-with-panic-unwind

# Build all libraries
build-libs: Cargo.lock
	docker run --rm -v $$PWD:/wd -w /wd rust:1.84.1-bullseye /bin/bash -c '\
		cargo build --profile release-with-panic-unwind --target x86_64-pc-windows-gnu; \
		cargo build --profile release-with-panic-unwind --target x86_64-apple-darwin; \
		cargo build --profile release-with-panic-unwind --target aarch64-apple-darwin; \
		cargo build --profile release-with-panic-unwind --target x86_64-unknown-linux-gnu; \
		cargo build --profile release-with-panic-unwind --target aarch64-unknown-linux-gnu; \
		mkdir -p xdr2json/libs/x86_64-pc-windows-gnu; \
		mkdir -p xdr2json/libs/x86_64-apple-darwin; \
		mkdir -p xdr2json/libs/aarch64-apple-darwin; \
		mkdir -p xdr2json/libs/x86_64-unknown-linux-gnu; \
		mkdir -p xdr2json/libs/aarch64-unknown-linux-gnu; \
		cp target/x86_64-pc-windows-gnu/release-with-panic-unwind/*.a xdr2json/libs/x86_64-pc-windows-gnu/; \
		cp target/x86_64-apple-darwin/release-with-panic-unwind/*.a xdr2json/libs/x86_64-apple-darwin/; \
		cp target/aarch64-apple-darwin/release-with-panic-unwind/*.a xdr2json/libs/aarch64-apple-darwin/; \
		cp target/x86_64-unknown-linux-gnu/release-with-panic-unwind/*.a xdr2json/libs/x86_64-unknown-linux-gnu/; \
		cp target/aarch64-unknown-linux-gnu/release-with-panic-unwind/*.a xdr2json/libs/aarch64-unknown-linux-gnu/; \
		shasum xdr2json/libs/**/*.a; \
		'

dist-clean:
	@rm -rf $(BUILD_DIR) $(LIBS_DIR)
