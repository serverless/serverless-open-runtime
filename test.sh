#!/bin/bash -x

set -e

. ./env

LAYER="$(sls layers info -l testRuntime --latest-arn-only) $(sls layers info -l testMiddleware --latest-arn-only)"
ROLE=arn:aws:iam::377024778620:role/Test-Role

rm -f testRuntime-lambda.zip
cd example && zip -r ../testRuntime-lambda.zip handler.js node_modules && cd -
aws s3 cp ./testRuntime-lambda.zip s3://dschep-byol/testRuntime-lambda.zip
aws $beta lambda create-function \
    --function-name=testRuntime-test \
    --runtime=byol \
    --code=S3Bucket=dschep-byol,S3Key=testRuntime-lambda.zip \
    --role=$ROLE \
    --timeout=10 \
    --handler=handler.hello \
    --layers $LAYER

aws $beta lambda invoke --function-name=testRuntime-test --log-type Tail out > resp
jq -r .LogResult resp | base64 --decode
cat out #python -m json.tool out
aws $beta lambda delete-function --function-name=testRuntime-test
