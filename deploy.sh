#!/bin/bash
mkdir build  &&
GOOS=linux go build -o build/main cmd/main.go &&
zip build/deployment.zip build/main &&
aws cloudformation package --template-file ./sam.yaml --output-template-file build/new_sam_spec.yaml --s3-bucket yuva-lambda-test &&
aws cloudformation deploy --template-file build/new_sam_spec.yaml --stack-name samoauthgo --capabilities CAPABILITY_IAM &&
rm -rf build