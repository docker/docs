var assert = require('assert');
var ExtractTextPlugin = require("extract-text-webpack-plugin");
var ENV_CONFIG = require('./_envConfig.js');
var _ = require('lodash');
var commonConfig = require('./_common.webpack.config');
var debug = require('debug')('webpack--server');
var fs = require('fs');

process.env.CLIENT_JS_FILENAME = fs.readFileSync('/tmp/.client-js-hash', 'utf-8').split('js/')[1];
var CSS_FILENAME = _.filter(fs.readdirSync('/opt/hub/.build-prod/public/styles/'),
                                    function(str) {
                                      return !/.*\.map$/.test(str);
                                    });

assert.strictEqual(1, CSS_FILENAME.length, CSS_FILENAME);
process.env.CSS_FILENAME = CSS_FILENAME[0];
fs.writeFileSync('/tmp/.client-js-hash', CSS_FILENAME);

/**
 * blacklist this array from being included in `externals`.
 *
 * This has the effect of making any modules in this list be 
 * resolved at build time instead of runtime. This affects the
 * server bundle
 */
var blacklist = ['.bin', 'hub-js-sdk', 'dux'];
var node_modules = fs.readdirSync('node_modules').filter(function(x) {
  return !_.includes(blacklist, x);
});

debug('modules that will be runtime require dependencies of the server: ', node_modules);


var serverConfig = {
  entry: '/opt/hub/app/scripts/server.js',
  output: {
    path: '.build-prod/',
    filename: 'server.js',
    libraryTarget: 'commonjs2'
  },
  plugins: [
    ENV_CONFIG,
    new ExtractTextPlugin('/.ignore/whatever.css', {
      allChunks: true
    })
  ],
  target: 'node',
  externals: node_modules,
  node: {
    __dirname: '/opt/hub/'
  }
};

module.exports = _.assign({}, commonConfig, serverConfig);
