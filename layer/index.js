const fetch = require('node-fetch')
const AWS = require('aws-sdk')

// BETA ENV STUFF
const { Agent } = require('https')
AWS.config.update({ region: 'us-west-2' })
AWS.NodeHttpClient.sslAgent = new Agent({ rejectUnauthorized: false })

const lambda = new AWS.Lambda({
  // BETA ENV STUFF
  endpoint: 'https://beta-04-2014-fe-alb-1663246439.us-west-2.elb.amazonaws.com'
})


module.exports = async function() {
  console.log("Starting node sls runtime.")

  // Add relevant paths to node's search path
  module.paths.unshift('/opt/node_modules', '/var/task/node_modules', '/var/task')

  // Runtime API URL
  const runtimeAPI = `http://${process.env.AWS_LAMBDA_RUNTIME_API}/2018-06-01/runtime`

  const lambdaInfo = await lambda.getFunction({FunctionName: process.env.AWS_LAMBDA_FUNCTION_NAME}).promise()

  // import the user's function
  const [modName, funcName] = process.env._HANDLER.split('.')
  const mod = require(modName)
  // wrap with a Promise just in case the user defined a syncronous function
  let func = (event) => Promise.resolve(mod[funcName](event))

  // Apply any middlewares
  for (const {Arn} of lambdaInfo.Configuration.Layers) {
    const middlewareName = Arn.split(':')[6]
    if (middlewareName === 'testRuntime') {
      continue
    }
    const middleware = require(middlewareName)
    func = middleware(func)
  }


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
