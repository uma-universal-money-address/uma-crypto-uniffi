#!/usr/bin/env bash

set -euo pipefail
python --version
pip install -r requirements.txt

echo "Generating python file..."
cd ..
cargo run --bin uniffi-bindgen generate src/uma_crypto.udl --language python --out-dir uma-crypto-python/src/uma_crypto/ --no-format

echo "Generating native binaries..."
docker buildx build -f build.Dockerfile --platform linux/arm64 -o docker-out .

echo "Copying linux binary..."
cp docker-out/target/aarch64-unknown-linux-gnu/release-smaller/libuma_crypto.so uma-crypto-python/src/uma_crypto/libuniffi_uma_crypto.so

echo "All done!"
