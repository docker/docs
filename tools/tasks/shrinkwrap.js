const fs = require('fs');
const path = require('path');
const promisify = require('promisify-node');
const childSpawn = require('../lib/spawn');
const log = require('../lib/log');

const rootDir = path.resolve(__dirname, '../..');

function removeDeps(keys) {
  const shrinkwrapPath = path.join(rootDir, 'npm-shrinkwrap.json');
  const shrinkwrapFile = require(shrinkwrapPath);
  keys.forEach((key) => delete shrinkwrapFile.dependencies[key]);
  return promisify(
    fs.writeFile(shrinkwrapPath, JSON.stringify(shrinkwrapFile, null, 2))
  );
}

function runNpmPrune() {
  return childSpawn(
    'npm',
    ['prune']
  );
}

function runShrinkwrap() {
  return childSpawn(
    'npm',
    ['shrinkwrap']
  );
}

function removeOptionalDeps() {
  log('removing optional dependencies');
  const { optionalDependencies } = require(path.join(rootDir, 'package.json'));
  return removeDeps(Object.keys(optionalDependencies));
}

module.exports = function shrinkwrap() {
  return runNpmPrune().then(runShrinkwrap).then(removeOptionalDeps);
};
