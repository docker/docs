'use strict';

module.exports = require('./webpack.config.builder')({
  devServer: true,
  devtool: 'eval',
  debug: true,
  failOnError: false
});
