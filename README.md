# The Serverless Open Runtime for AWS Lambda

This uses AWS's Runtime API for Lambda to implement a universal runtime.

When done it will feature (exact list TBD):
 * [CloudEvent](https://cloudevents.io/) based function signature
 * [Middlewares](#middlewares)
 * [Language Agnostic](#language-agnostic)
 * Graceful timeout handling - [poc branch](https://github.com/serverless/open-runtime-poc/tree/timeout-test)
 * more!

## Try it out


```shell
npm i -g serverless # make sure you have serverless framework verison 1.34.0 or greater
make
sls deploy
cd example
sls deploy
sls invoke -f hello
```

## NodeJS debugging!
The open runtime implementation Node JS supports remote debugging. To use it:
1. ensure you have a high timeout on your lambda (or else it'll be killed while debugging)
2. on a publicly accessible server, install `socat` and run:
   ```
   socat -v TCP-LISTEN:9999,reuseaddr,fork TCP-LISTEN:4444,reuseaddr
   ```
3. invoke your lambda with the `_debugProxy` key in the event set to `1.1.1.1:4444` but changing
   `1.1.1.1` to the IP address of the server you stared `socat` on:
   ```
   sls invoke -f hello -d '{"_debugProxy": "1.1.1.1:4444"}'
   ```
4. open the chrome debugger to [`chrome-devtools://devtools/bundled/inspector.html?ws=1.1.1.1:4444`](chrome-devtools://devtools/bundled/inspector.html?ws=1.1.1.1:4444) (again changing 1.1.1.1 to your server's IP)

## Middlewares
The current plan for middlewares allows them to be written in any language
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

### Potential downsides of this middleware architecture
How will can this be leveraged to implement APM style monitoring that usually
monkeypatches language constructs.

## Language Agnostic
The current Proof of Concept is implemented in NodeJS 10, but the final open runtime
will likely be implemented in Go. A language runtime will be specified to the open
runtime to support specific languages.

This allows the open runtime to have a single implementation  while supporting many languages.
A language-runtime will be started when the open-runtime starts and communicate with it via STDIN &
STDOUT as it receives events. There is currently an implementation of this for nodejs 10
