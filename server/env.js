var path = require('path');
var fs = require('fs');

require('css-modules-require-hook')({
  generateScopedName: '[name]__[local]___[hash:base64:5]',
});

require('app-module-path').addPath(path.join(__dirname, '..', 'src'));
require('app-module-path')
  .addPath(path.join(__dirname, '..', 'src', 'components'));
require('app-module-path')
  .addPath(path.join(__dirname, '..', 'src', 'lib', 'css'));
require('babel-register', {
  plugins: ['add-module-exports'],
});

// Webpack loaders aren't used Server-side, so we need to handle .md files
require.extensions['.md'] = function (module, filename) {
  // eslint-disable-next-line
  module.exports = fs.readFileSync(filename, 'utf8');
};
