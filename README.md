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
  const resp = await handler(event)
}
```

Handlers must be async and can safely assume that the handler it is provided also is async.

It is preferable to package a middleware as a layer, but that is optional. The main requirement is
that it be `require`able and that the middleware be exported as default.

To specify the handler middlewares your function should use, set the `SLSMIDDLEWARES` environment
variable to a comma delimited list of middlewares by the name by which they are to be `require`d
