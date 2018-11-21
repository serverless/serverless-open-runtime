#!/bin/bash -x

set -e

. ./env

LAYER="$(sls layers info -l testRuntime --latest-arn-only) $(sls layers info -l testMiddleware --latest-arn-only)"
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
    --environment 'Variables={SLSMIDDLEWARES=test-middleware}' \
    --layers $LAYER

aws $beta lambda invoke --function-name=testRuntime-test --log-type Tail out > resp
jq -r .LogResult resp > log
# py = https://gist.github.com/dschep/4358be665537463b9271f782e77ff85f
py 'print(base64.b64decode(open("log").read()))'
cat out #python -m json.tool out
aws $beta lambda delete-function --function-name=testRuntime-test
