# !/bin/bash

# this build file is used for single module builder
poetry build
poetry run pip install --upgrade -t bundle dist/*.whl
cd bundle ; zip -r ../lambda_bundle.zip . -x '*.pyc'
cd ..
rm -rf dist bundle
