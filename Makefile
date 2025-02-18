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
	docker run --rm -v $$PWD:/wd -w /wd --platform=linux/amd64 rust:1.84.1-bullseye /bin/bash -c '\
		for target in $(TARGETS); do \
			cargo build --profile $(PROFILE) --target $$target; \
			mkdir -p $(LIBS_DIR)/$$target; \
			cp $(BUILD_DIR)/$$target/$(PROFILE)/*.a $(LIBS_DIR)/$$target/; \
		done; \
		shasum -a 256 $(LIBS_DIR)/**/*.a; \
		'

dist-clean:
	@rm -rf $(BUILD_DIR) $(LIBS_DIR)
