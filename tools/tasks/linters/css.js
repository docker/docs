const spawn = require('../../lib/spawn');

module.exports = function cssLint() {
  return spawn('npm', ['run', 'lint:css']);
};
