'use strict';
var keyMirror = require('keymirror');

// Component-Global form states
var STATUS = keyMirror({
  DEFAULT: null,
  SAVING: null
});

var EMAILSTATUS = keyMirror({
  SUCCESS: null,
  ATTEMPTING: null,
  FAILED: null
});

module.exports = {
  STATUS,
  EMAILSTATUS
};
