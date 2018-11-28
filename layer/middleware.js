const { spawnSync } = require("child_process");

module.exports = middlewareName => handler => event => {
  const options = {
    stdio: ["pipe", "pipe", "inherit"],
    env: {
      ...process.env,
      PATH: `/opt/runtime/bin:${process.env.path}`
    }
  }
  const beforeResult = spawnSync(
    `/opt/middlewares/${middlewareName}`,
    ["before"],
    { ...options, input: JSON.stringify(event) }
  );
  const processedEvent = JSON.parse(beforeResult.stdout.toString());

  const result = handler(processedEvent);

  const afterResult = spawnSync(
    `/opt/middlewares/${middlewareName}`,
    ["after"],
    { ...options, input: JSON.stringify(result) }
  );
  const processedResult = JSON.parse(afterResult.stdout.toString());
  return processedResult;
};
