const fs = require('fs');
const path = require('path');
const promisify = require('promisify-node');
const log = require('../../lib/log');

module.exports = function convert() {
  const modulesDir = path.resolve(__dirname, '../../../node_modules');
  const pkgPath = path.resolve(__dirname, '../../../package.json');
  const pkg = require(pkgPath);

  const convertToExact = (object) => {
    const semverRegex = /^[\^\d~]/;
    const response = object;

    Object.keys(object).forEach((module) => {
      const version = object[module];

      let file;

      if (!semverRegex.test(version)) return;

      try {
        file = require(path.join(modulesDir, module, 'package.json'));
      } catch (e) {
        log.error(`${module}/package.json not found`);
        return;
      }

      response[module] = file.version;
    });

    return response;
  };

  convertToExact(pkg.dependencies);
  convertToExact(pkg.devDependencies);
  return promisify(fs.writeFile(pkgPath, JSON.stringify(pkg, null, 2)));
};
