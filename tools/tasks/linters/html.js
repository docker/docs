const spawn = require('../../lib/spawn');

module.exports = function htmlLint() {
  return spawn('npm', ['run', 'lint:html']);
};
