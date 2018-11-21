# The Serverless NodeJS Runtime for AWS Lambda

This uses AWS's Runtime API for Lambda to implement a more universal runtime for Node JS.

It is deployed using the [serverless-lambda-layer plugin](https://github.com/serverless/lambda-layers-plugin)

When done it will feature (exact list TBD):
 * [CloudEvent]() based function signature
 * [Middlewares](#Middlewares)
 * Graceful timeout handling
 * more!

## Try it out
You need to use our whitelisted AWS account, msg @dschep for info.

The test script utilizes the awscli instead of Serverless Framework because AWS hasn't finialized
CloudFormation support for layers yet.

```shell
sls layers deploy # deploy the runtime & a test middleware layer
./test.sh # create a lambda with example/handler.js & execute it & delete it
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

Middlewares must be packaged and included as Lambda Layers, and the name of the layer must be the
same as the node module to `require`.

To specify the handler middlewares your function should use, include the layer a middleware is
included in with your lambda function.
