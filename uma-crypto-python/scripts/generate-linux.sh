#!/usr/bin/env bash

set -euo pipefail
${PYBIN}/python --version
${PYBIN}/pip install -r requirements.txt

echo "Generating python file..."
cd ..
cargo run --bin uniffi-bindgen generate src/uma_crypto.udl --language python --out-dir uma-crypto-python/src/uma_crypto/ --no-format

echo "Generating native binaries..."
cargo build --profile release-smaller

echo "Copying linux binary..."
cp target/release-smaller/libuma_crypto.so uma-crypto-python/src/uma_crypto/libuniffi_uma_crypto.so

echo "All done!"
