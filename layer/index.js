const fetch = require('node-fetch')

module.exports = async function() {
  console.log("Starting node sls runtime.")

  // Add relevant paths to node's search path
  module.paths.unshift('/opt/node_modules', '/var/task/node_modules', '/var/task')

  // Runtime API URL
  const runtimeAPI = `http://${process.env.AWS_LAMBDA_RUNTIME_API}/2018-06-01/runtime`

  // import the user's function
  const [modName, funcName] = process.env._HANDLER.split('.')
  const mod = require(modName)
  const func = mod[funcName]

  // lōōpz
  while (true) {
    // Request the next event from the Lambda Runtime
    const invocationResp = await fetch(`${runtimeAPI}/invocation/next`)

    // get the invocation ID & parse the event payload
    const invocationId = invocationResp.headers.get('x-amz-aws-request-id')
    const eventPayload = await invocationResp.json()

    console.log(`Invoke received. Request ID: ${invocationId}`)

    // Invoke the users's function
    const handlerResp = await Promise.resolve(func(eventPayload))

    // Send the response to Lambda Runtime
    await fetch(`${runtimeAPI}/invocation/${invocationId}/response`, {
      method: 'POST',
      body: JSON.stringify(handlerResp), 
      headers: { 'Content-Type': 'application/json' }
    })
  }
}
