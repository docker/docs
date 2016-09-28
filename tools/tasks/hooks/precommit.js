const path = require('path');
const childSpawn = require('../../lib/spawn');
const fileExists = require('../../lib/file-exists');
const filesChanged = require('../utils/files-changed');

module.exports = function precommit() {
  const files = {
    html: [],
    style: [],
    es: [],
  };

  // eslint-disable-next-line no-unused-vars
  const [node, runner, script, ...dirs] = process.argv;
  const changedFiles = filesChanged(undefined, dirs.concat['package.json']);
  const packageChanged = changedFiles.indexOf('package.json') !== -1;

  const root = path.resolve(__dirname, '../../../');

  changedFiles.forEach((file) => {
    // Cover files that have been deleted
    if (!fileExists(file, { root })) return;

    if (file.endsWith('.html')) {
      files.html.push(file);
    } else if (file.endsWith('.css')) {
      files.style.push(file);
    } else if (file.endsWith('.js')) {
      files.es.push(file);
    }
  });

  const processes = Object.keys(files)
    .filter((type) => !!files[type].length)
    .map((type) => childSpawn(`./node_modules/.bin/${type}lint`, files[type]));

  if (packageChanged) {
    processes.push(
      childSpawn('npm', ['run', 'shrinkwrap'])
        .then(() => childSpawn('git', ['add', 'npm-shrinkwrap.json']))
    );
  }

  return Promise.all(processes).catch(() => process.exit(1));
};
