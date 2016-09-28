const childSpawn = require('../../lib/spawn');
const filesChanged = require('../utils/files-changed');

module.exports = function postmerge() {
  const files = filesChanged(
    'git diff-tree -r --name-only --no-commit-id ORIG_HEAD HEAD'
  );

  if (files.indexOf('package.json') === -1) return Promise.resolve();

  return childSpawn('npm', ['i'])
    .then(() => childSpawn('npm', ['prune']))
    .catch(() => process.exit(1));
};
