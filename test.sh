#!/bin/bash -x

set -e

. ./env

# need to update the version in this ARN on every `sls layers publish`
LAYER=arn:aws:lambda:us-west-2:377024778620:layer:testRuntime:11
ROLE=arn:aws:iam::377024778620:role/Test-Role

rm -f testRuntime-lambda.zip
cd example && zip ../testRuntime-lambda.zip handler.js && cd -
aws s3 cp ./testRuntime-lambda.zip s3://dschep-byol/testRuntime-lambda.zip
aws $beta lambda create-function \
    --function-name=testRuntime-test \
    --runtime=byol \
    --code=S3Bucket=dschep-byol,S3Key=testRuntime-lambda.zip \
    --role=$ROLE \
    --timeout=5 \
    --handler=handler.hello \
    --layers=$LAYER

aws $beta lambda invoke --function-name=testRuntime-test --log-type Tail out > resp
jq -r .LogResult resp > log
# py = https://gist.github.com/dschep/4358be665537463b9271f782e77ff85f
py 'print(base64.b64decode(open("log").read()))'
cat out #python -m json.tool out
aws $beta lambda delete-function --function-name=testRuntime-test
