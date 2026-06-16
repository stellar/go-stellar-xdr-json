build:
	go build ./...

test:
	go test ./...

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

LIBS_DIR = xdrjson/libs
BUILD_DIR = target
PROFILE = release-with-panic-unwind
# Build image, pinned by digest (not just tag) so its OS layer can't drift, so
# that the binaries are reproducible.
IMAGE = rust:1.84.1-bullseye@sha256:1e3f7a9fd1f278cc4be02a830745f40fe4b22f0114b2464a452c50273cc1020d

# Build all libraries
build-libs:
	docker run --rm -v $$PWD:/wd -w /wd --platform=linux/amd64 $(IMAGE) /bin/bash -c '\
		rustc -vV > $(LIBS_DIR)/rust-version; \
		for target in $(TARGETS); do \
			cargo build --locked --profile $(PROFILE) --target $$target --package xdrjson; \
			mkdir -p $(LIBS_DIR)/$$target; \
			cp $(BUILD_DIR)/$$target/$(PROFILE)/*.a $(LIBS_DIR)/$$target/; \
		done; \
		shasum -a 256 $(LIBS_DIR)/**/*.a; \
		'

generate-types:
	cargo run --bin generate-types | gofmt > xdrjson/types.go

dist-clean:
	@rm -rf $(BUILD_DIR) $(LIBS_DIR)
