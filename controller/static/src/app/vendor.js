'use strict';

module.exports = function () {
  global.$ = global.jQuery = require('jquery');

  require('../../node_modules/codemirror/lib/codemirror.css');
  global.CodeMirror = require('codemirror');
  require('../../node_modules/codemirror/mode/yaml/yaml');

  require('../styles/main.css');
  require('../styles/nv.d3.css');

  require('../../semantic/dist/semantic.css');
  require('../../semantic/dist/semantic.js');


  require('jquery-tablesort');
  require('angular');
  require('angular-ui-router');
  require('angular-ui-codemirror');
  require('angular-jwt');
  require('angular-breadcrumb');
  require('angular-sanitize');
  require('angular-storage');
  require('angularjs-scroll-glue');
  require('ng-file-upload');
  require('ngclipboard');
  require('moment');
  require('d3');
  require('nvd3');
  require('oboe');

  require('ng-table');
};
