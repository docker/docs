'use strict';

module.exports = require('./webpack.config.builder')({
  minify: !process.env.SKIP_MEDIA_MINIFY,
  failOnError: true
});
