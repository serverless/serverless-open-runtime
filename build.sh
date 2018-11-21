#!/bin/bash
VERSION=10.13.0

git clean -dfx layer

cd layer
curl -O https://nodejs.org/dist/v$VERSION/node-v$VERSION-linux-x64.tar.xz
tar -xf node-v$VERSION-linux-x64.tar.xz
mv node-v$VERSION-linux-x64 runtime
npm i
cp ../lambda-2015-03-31.normal.json node_modules/aws-sdk/apis/lambda-2015-03-31.min.json
