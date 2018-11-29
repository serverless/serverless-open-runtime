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

First, you need to install private sls layers branch & configure aws-sdk for layers:
```shell
git clone git@github.com:serverless/nda-serverless
cd nda-serverless
git checkout layers
npm i -g .
cd ..
```

```shell
sls deploy
# Update ARN in example/serverless.yml with the one that was just printed in the deploy
cd example
sls deploy
sls invoke -f hello
```

## Middlewares
The current proof of concept for middlewares allows them to be written in any language
by invoking the middleware as an executable with the event or response passed in via
standard in & out. The first argument specifies the hook that is being invoked.

Middlewares are best stored in layers and must be stored in the `middlewares` directory
(IE: `/opt/middlewares/` inside the lambda execution environment)

To specify the handler middlewares your function should use, set the `SLSMIDDLEWARES` environment
variable to a comma delimited list of middlewares by the filename of the middleware executable.
