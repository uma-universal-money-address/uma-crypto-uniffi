setup-go-targets:
	rustup target add x86_64-apple-darwin aarch64-apple-darwin 

setup-jvm-targets:
	rustup target add x86_64-apple-darwin aarch64-apple-darwin

code-gen-kotlin:
	cargo run --bin uniffi-bindgen generate src/uma_crypto.udl --language kotlin --out-dir uma-crypto-kotlin
	mv uma-crypto-kotlin/uniffi/uma_crypto/uma_crypto.kt uma-crypto-kotlin/uniffi/uma_crypto/UmaCrypto.kt
	sed -i '' 's/package uniffi.uma_crypto/package me.uma.crypto.internal/g' uma-crypto-kotlin/uniffi/uma_crypto/UmaCrypto.kt

build-darwin-amd64:
	cargo build --profile release-smaller --target x86_64-apple-darwin

build-darwin-arm64:
	cargo build --profile release-smaller --target aarch64-apple-darwin

code-gen-go:
	mkdir -p uma-crypto-go/internal
	cargo install uniffi-bindgen-go --git https://github.com/NordSecurity/uniffi-bindgen-go
	uniffi-bindgen-go src/uma_crypto.udl --out-dir uma-crypto-go
	mv uma-crypto-go/uniffi/uma_crypto/uma_crypto.go uma-crypto-go/internal
	sed -i '' 's/package uma_crypto/package internal/g' uma-crypto-go/internal/uma_crypto.go
	rm -rf uma-crypto-go/uniffi

build-linux-amd64-static:
	docker buildx build -f build.Dockerfile --platform linux/amd64 -o docker-out .

build-linux-arm64-static:
	docker buildx build -f build.Dockerfile --platform linux/arm64 -o docker-out .

build-linux-amd64-shared:
	docker buildx build -f build.Dockerfile --build-arg CDYLIB=true --platform linux/amd64 -o docker-out .

build-linux-arm64-shared:
	docker buildx build -f build.Dockerfile  --build-arg CDYLIB=true --platform linux/arm64 -o docker-out .

go-libs: build-darwin-amd64 build-darwin-arm64 build-linux-amd64-static build-linux-arm64-static
	mkdir -p uma-crypto-go/libs/darwin/amd64
	mkdir -p uma-crypto-go/libs/darwin/arm64
	mkdir -p uma-crypto-go/libs/linux/amd64
	mkdir -p uma-crypto-go/libs/linux/arm64
	cp target/x86_64-apple-darwin/release-smaller/libuma_crypto.a uma-crypto-go/libs/darwin/amd64
	cp target/aarch64-apple-darwin/release-smaller/libuma_crypto.a uma-crypto-go/libs/darwin/arm64
	cp docker-out/target/x86_64-unknown-linux-gnu/release-smaller/libuma_crypto.a uma-crypto-go/libs/linux/amd64
	cp docker-out/target/aarch64-unknown-linux-gnu/release-smaller/libuma_crypto.a uma-crypto-go/libs/linux/arm64

build-go: setup-go-targets code-gen-go go-libs

build-jvm-targets: setup-jvm-targets build-darwin-amd64 build-darwin-arm64 build-linux-amd64-shared build-linux-arm64-shared

jvm-libs: build-jvm-targets
	mkdir -p uma-crypto-kotlin/jniLibs/jvm/darwin-aarch64
	mkdir -p uma-crypto-kotlin/jniLibs/jvm/darwin-x86-64
	mkdir -p uma-crypto-kotlin/jniLibs/jvm/linux-x86-64
	mkdir -p uma-crypto-kotlin/jniLibs/jvm/linux-aarch64
	cp -r target/aarch64-apple-darwin/release-smaller/libuma_crypto.dylib uma-crypto-kotlin/jniLibs/jvm/darwin-aarch64/libuniffi_uma_crypto.dylib
	cp -r target/x86_64-apple-darwin/release-smaller/libuma_crypto.dylib uma-crypto-kotlin/jniLibs/jvm/darwin-x86-64/libuniffi_uma_crypto.dylib
	cp -r docker-out/target/x86_64-unknown-linux-gnu/release-smaller/libuma_crypto.so uma-crypto-kotlin/jniLibs/jvm/linux-x86-64/libuniffi_uma_crypto.so
	cp -r docker-out/target/aarch64-unknown-linux-gnu/release-smaller/libuma_crypto.so uma-crypto-kotlin/jniLibs/jvm/linux-aarch64/libuniffi_uma_crypto.so

build-jvm: setup-jvm-targets jvm-libs
