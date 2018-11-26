const fetch = require('node-fetch')
const shim = require('./shim')

module.exports = async function() {
  console.log("Starting node sls runtime.")

  // Add relevant paths to node's search path
  module.paths.unshift('/opt/node_modules', '/var/task/node_modules', '/var/task')

  // Runtime API URL
  const runtimeAPI = `http://${process.env.AWS_LAMBDA_RUNTIME_API}/2018-06-01/runtime`

  // import the user's function
  const [modName, funcName] = process.env._HANDLER.split('.')
  const mod = require(modName)
  // wrap with a Promise just in case the user defined a syncronous function
  let func = (event) => Promise.resolve(mod[funcName](event))

  // Apply any middlewares
  for (const middlewareName of (process.env.SLSMIDDLEWARES||'').split(',')) {
    const middleware = require(middlewareName)
    func = middleware(func)
  }


  // lōōpz
  while (true) {
    // Request the next event from the Lambda Runtime
    const invocationResp = await fetch(`${runtimeAPI}/invocation/next`)

    // get the invocation ID & parse the event payload
    const invocationId = invocationResp.headers.get('lambda-runtime-aws-request-id')
    const eventPayload = await invocationResp.json()

    console.log(`Invoke received. Request ID: ${invocationId}`)

    // Invoke the users's function, converting event or events to CloudEvent form
    let body
    if (eventPayload.Records) {
      await Promise.all(eventPayload.Records.map(shim.transformAsyncEvent).map(func))
    }
    else {
      const handlerResp = await func(shim.transformSyncEvent(eventPayload))
      body = JSON.stringify(handlerResp)
    }

    // Send the response to Lambda Runtime
    await fetch(`${runtimeAPI}/invocation/${invocationId}/response`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body
    })
  }
}
