'use strict';
var keyMirror = require('keymirror');

// Component-Global form states
var STATUS = keyMirror({
  DEFAULT: null,
  ATTEMPTING_SIGNUP: null,
  SUCCESSFUL_SIGNUP: null
});

module.exports = {
  STATUS
};
