#!/usr/bin/env bash

set -euo pipefail
python3 --version
pip install --user -r requirements.txt

echo "Generating python file..."
cd ..
cargo run --bin uniffi-bindgen generate src/uma_crypto.udl --language python --out-dir uma-crypto-python/src/uma_crypto/ --no-format

echo "Generating native binaries..."
rustup target add x86_64-pc-windows-msvc
cargo build --profile release-smaller --target x86_64-pc-windows-msvc

echo "Copying libraries..."
cp target/x86_64-pc-windows-msvc/release-smaller/uma_crypto.dll uma-crypto-python/src/uma_crypto/uma_crypto.dll

echo "All done!"
