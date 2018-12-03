#!/usr/bin/env node

const split = require('split')

process.stdin.pipe(split()).on('data', processEvent)

module.paths.push(process.cwd())
const [module_name, func_name] = process.argv[2].split('.')
const package = require(module_name)
const handler = package[func_name]

function processEvent (event) {
  if (!event)
    return
  const parsedEvent = JSON.parse(event)
  const result = handler(event)
  console.log(JSON.stringify(result))
}

