const Promise = require('bluebird')
module.exports.hello = async (event) => {
  await Promise.delay(1000*15)
  return {body: JSON.stringify(event)}
}
