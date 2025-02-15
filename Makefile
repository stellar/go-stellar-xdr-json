# update the Cargo.lock every time the Cargo.toml changes.
Cargo.lock: Cargo.toml
	cargo update --workspace

build-libs: Cargo.lock
	rustup target add x86_64-pc-windows-gnu && \
	rustup target add x86_64-apple-darwin && \
	rustup target add aarch64-apple-darwin && \
	rustup target add x86_64-unknown-linux-gnu && \
	rustup target add aarch64-unknown-linux-gnu && \
	cargo build --target x86_64-pc-windows-gnu --profile release-with-panic-unwind && \
	cargo build --target x86_64-apple-darwin --profile release-with-panic-unwind && \
	cargo build --target aarch64-apple-darwin --profile release-with-panic-unwind && \
	cargo build --target x86_64-unknown-linux-gnu --profile release-with-panic-unwind && \
	cargo build --target aarch64-unknown-linux-gnu --profile release-with-panic-unwind && \
	mkdir -p xdr2json/libs
	mkdir -p xdr2json/libs/x86_64-pc-windows-gnu
	mkdir -p xdr2json/libs/x86_64-apple-darwin
	mkdir -p xdr2json/libs/aarch64-apple-darwin
	mkdir -p xdr2json/libs/x86_64-unknown-linux-gnu
	mkdir -p xdr2json/libs/aarch64-unknown-linux-gnu
	cp target/x86_64-pc-windows-gnu/release-with-panic-unwind/*.a xdr2json/libs/x86_64-pc-windows-gnu/
	cp target/x86_64-apple-darwin/release-with-panic-unwind/*.a xdr2json/libs/x86_64-apple-darwin/
	cp target/aarch64-apple-darwin/release-with-panic-unwind/*.a xdr2json/libs/aarch64-apple-darwin/
	cp target/x86_64-unknown-linux-gnu/release-with-panic-unwind/*.a xdr2json/libs/x86_64-unknown-linux-gnu/
	cp target/aarch64-unknown-linux-gnu/release-with-panic-unwind/*.a xdr2json/libs/aarch64-unknown-linux-gnu/


