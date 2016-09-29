'use strict';
var keyMirror = require('keymirror');

// Component-Global form states
var STATUS = keyMirror({
  DEFAULT: null,
  ERROR: null,
  DELETING: null,
  CONFIRMING: null
});

module.exports = STATUS;
