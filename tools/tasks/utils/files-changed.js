const path = require('path');
const { execSync } = require('child_process');

module.exports = function filesChanged(
  gitCmd = 'git diff-index --name-only HEAD',
  dirs = []
) {
  let cmd = gitCmd;

  if (dirs.length) {
    cmd += ` ${dirs.join(' ')}`;
  }

  return execSync(cmd, { cwd: path.resolve(__dirname, '../../../') })
    .toString()
    .split('\n')
    .map((value) => value.trim())
    .filter((value) => !!value);
};
