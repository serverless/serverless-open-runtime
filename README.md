# The Serverless Open Runtime for AWS Lambda (Proof of Concept)

This uses AWS's Runtime API for Lambda to implement a universal runtime.

When done it will feature (exact list TBD):
 * [CloudEvent](https://cloudevents.io/) based function signature
 * [Middlewares](#middlewares)
 * [Language Agnostic](#language-agnostic)
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
./build.sh
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
*Note:* it's possible to make a layer imply a middleware by searching for the file of the same name
as the layer, but it requires an API call to get the lambda's configuration at runtime startup
(slower cold start) and might be more AWS specific than is desireable.

### Middleware specifics example
If you set `SLSMIDDLEWARES=test` when an event, say `{"foo": "bar"}` is recieved, the equivalent
following will be executed:
```
echo '{"foo": "bar"}' | /opt/middlewares/test before
```
and the STDOUT of the that execution will be read and replace the original event.

Then, when the user handler returns a response, say: `{"body": "hello"}` the middleware will
similarlly be invoked as the equivalent of:
```
echo '{"body": "hello"}' | /opt/middlewares/test after
```
and the STDOUT will replace the original response.

## Language Agnostic
The current Proof of Concept is implemented in NodeJS 10, but the final open runtime
will likely be implemented in Go. A language runtime will be specified to the open
runtime to support specific languages.

This allows the open runtime to have a single implementation  while supporting many languages.
There are two main ideas for how to implement a language runtime:
 * The language runtime is invoked all middlewares have processed the event
   * pro: simplicity
   * con: startup cost associated with languages (eg: Java & Python)
 * The runtime starts the language runtime at startup and communicates with it
   * pro: no startup cost per request
   * con: more complicated
