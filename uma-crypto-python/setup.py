from setuptools import setup

setup(
    name="uma_crypto_python",
    version="0.1.0",
    description="The Python language bindings for UMA crypto operations",
    long_description=open("README.md").read(),
    long_description_content_type="text/markdown",
    include_package_data = True,
    zip_safe=False,
    packages=["uma_crypto"],
    package_dir={"uma_crypto": "./src/uma_crypto"},
    url="https://github.com/uma-universal-money-address/uma-crypto-uniffi",
    author="Lightspark Group, Inc. <info@lightspark.com>",
    license="Apache 2.0",
    has_ext_modules=lambda: True,
)
