const spawn = require('../../lib/spawn');

module.exports = function jsLint() {
  return spawn('npm', ['run', 'lint:js']);
};
