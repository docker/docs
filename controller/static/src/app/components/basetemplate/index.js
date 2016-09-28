'use strict';

var angular = require('angular');

require('./layout.css');

module.exports = angular.module('ducp.layout', [])
  .controller('BaseController', require('./base.controller'))
  .controller('FooterController', require('./footer.controller'))
  .controller('BannerController', require('./banner.controller.js'))
  .controller('SupportController', require('./support.controller.js'))
  .config(require('./config.routes'));
