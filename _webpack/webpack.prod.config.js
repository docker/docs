var ExtractTextPlugin = require("extract-text-webpack-plugin");
var ENV_CONFIG = require('./_envConfig.js');
var HASH_CLIENT = require('./hashAndReplaceClient');
var _ = require('lodash');
var commonConfig = require('./_common.webpack.config');
var webpack = require('webpack');
var debug = require('debug')('webpack--client');

var clientConfig = {
  entry: '/opt/hub/app/scripts/client.js',
  devtool: 'source-map',
  output: {
    path: '.build-prod/public/',
    filename: 'js/client.[hash].js'
  },
  plugins: [
    ENV_CONFIG,
    new webpack.optimize.DedupePlugin(),
    new webpack.optimize.UglifyJsPlugin({
      compress: {
        warnings: false
      }
    }),
    new ExtractTextPlugin('styles/[name]-[id]-[hash].css', { allChunks: true }),
    HASH_CLIENT
  ]
};
var clientBundle = _.assign({}, commonConfig, clientConfig);

module.exports = [
  clientBundle
];
