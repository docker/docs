// Tasks are described by the name of the task file (./tasks/<name>)
// Tasks in arrays are run in parallel
// Tasks separated by ! are run sequentially

function run(filename, options) {
  const log = require('./lib/log');
  const moment = require('moment');

  const start = new Date();

  let fn;

  if (typeof filename === 'function') {
    fn = filename;
  } else {
    fn = require(`./tasks/${filename}`);
  }

  log(`Starting '${log.chalk.bold.green(fn.name)}'...`);
  return Promise.all([fn(options)]).then(() => {
    const end = new Date();
    const duration = end.getTime() - start.getTime();

    log([
      `Finished '${log.chalk.bold.green(fn.name)}' after`,
      log.chalk.bold.green(`${moment.duration(duration).as('seconds')}s`),
    ].join(' '));
  });
}

function runSequentialTasks(tasks) {
  const tasksArray = Array.isArray(tasks) ? tasks : [tasks];
  return tasksArray.reduce((promise, filename) => {
    // eslint-disable-next-line no-use-before-define
    const task = Array.isArray(filename) ? processTask : run;
    return promise.then(task.bind(null, filename));
  }, Promise.resolve());
}

function processTask(task) {
  let promise;

  if (Array.isArray(task)) {
    // Run all of these tasks in parallel
    promise = Promise.all(task.map(processTask));
  } else {
    // Run all tasks with '!' sequentially
    promise = runSequentialTasks(task.split('!'));
  }

  return promise;
}

function runner(tasks) {
  const log = require('./lib/log');

  return runSequentialTasks(tasks).catch((err) => {
    log.error(err);
    process.exit(1);
  });
}

module.exports = {
  run,
  runner,
};

if (process.mainModule.children.length === 0 && process.argv.length > 2) {
  delete require.cache[__filename];
  const module = process.argv[2];
  run(module).catch(err => require('./lib/log').error(err.stack));
}
