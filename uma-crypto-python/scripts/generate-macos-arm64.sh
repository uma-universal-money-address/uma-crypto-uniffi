#!/usr/bin/env bash

set -euo pipefail
python3 --version
pip install --user -r requirements.txt

echo "Generating python file..."
cd ..
cargo run --bin uniffi-bindgen generate src/uma_crypto.udl --language python --out-dir uma-crypto-python/src/uma_crypto/ --no-format

echo "Generating native binaries..."
rustup target add aarch64-apple-darwin
cargo build --profile release-smaller --target aarch64-apple-darwin

echo "Copying libraries dylib..."
cp target/aarch64-apple-darwin/release-smaller/libuma_crypto.dylib uma-crypto-python/src/uma_crypto/libuniffi_uma_crypto.dylib

echo "All done!"

