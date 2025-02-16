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
PROFILE = release-with-panic-unwind

# Build all libraries
build-libs: Cargo.lock
	@for target in $(TARGETS); do \
		rustup target add $$target; \
		cargo build --target $$target --profile $(PROFILE); \
		mkdir -p $(LIBS_DIR)/$$target; \
		cp target/$$target/$(PROFILE)/*.a $(LIBS_DIR)/$$target/; \
	done

