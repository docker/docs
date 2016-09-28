const moment = require('moment');
const chalk = require('chalk');

function timeStamp(date = new Date()) {
  return chalk.gray(`[${moment(date).format('HH:mm:ss')}]`);
}

function log(msg) {
  console.log([timeStamp(), msg].join(' '));
}

log.chalk = chalk;

log.error = function logError(msg) {
  log(chalk.yellow(msg));
};

log.success = function logSuccess(msg) {
  log(chalk.green(msg));
};

log.info = function logInfo(msg) {
  log(chalk.italic.gray(msg));
};

module.exports = log;
