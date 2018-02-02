#!/bin/bash
echo 'Build started'
mkdir build  &&
GOOS=linux go build -o build/main cmd/main.go &&
zip -j build/deployment.zip build/main &&
echo 'Creating a package'
aws cloudformation package --template-file ./sam.yaml --output-template-file build/new_sam_spec.yaml --s3-bucket yuva-lambda-test &&
echo 'Deploying'
aws cloudformation deploy --template-file build/new_sam_spec.yaml --stack-name samoauthgo --capabilities CAPABILITY_IAM &&
rm -r build