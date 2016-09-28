const fs = require('fs');
const path = require('path');

module.exports = function fileExists(filepath, options = {}) {
  if (!filepath) return false;

  const fullpath = options.root ? path.join(options.root, filepath) : filepath;

  try {
    return fs.statSync(fullpath).isFile();
  } catch (err) {
    return false;
  }
};
