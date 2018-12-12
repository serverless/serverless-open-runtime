module.exports.hello = async (event) => {
  return {body: JSON.stringify(event)}
}
