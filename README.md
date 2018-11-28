# The Serverless Open Runtime for AWS Lambda (Proof of Concept)

This uses AWS's Runtime API for Lambda to implement a more universal runtime for Node JS.

It is deployed using the [serverless-lambda-layer plugin](https://github.com/serverless/lambda-layers-plugin)

When done it will feature (exact list TBD):
 * [CloudEvent](https://cloudevents.io/) based function signature
 * [Middlewares](#Middlewares)
 * Graceful timeout handling
 * more!

## Try it out
You need to use our whitelisted AWS account: 377024778620

First, you need to install non-public sls layers branch & configure aws-sdk for layers:
```shell
git clone git@github.com:serverless/nda-serverless
cd nda-serverless
git checkout layers
npm i -g .
cd ..
clone git@github.com:serverless/lambda-layers-plugin
cp lambda-layers-plugin/lambda-2015-03-31.normal.json $(dirname $(dirname $(which sls)))/lib/node_modules/serverless/node_modules/aws-sdk/apis/lambda-2015-03-31.min.json
```

```shell
cd example
sls deploy
sls invoke -f hello
```

## Middlewares
The current proof of concept for middlewares is built around the decorator pattern. As such, a
middleware should be a function that accepts a handler and returns a new handler. EG:
```javascript
const middleware = async (handler) => (event) => {
  // do something before invocation
  const resp = await handler(event)
  // do something after invocation
  return resp
}
```

Handlers must be async and can safely assume that the handler it is provided also is async.

It is preferable to package a middleware as a layer, but that is optional. The main requirement is
that it be `require`able and that the middleware be exported as default.

To specify the handler middlewares your function should use, set the `SLSMIDDLEWARES` environment
variable to a comma delimited list of middlewares by the name by which they are to be `require`d
