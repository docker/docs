'use strict';
var keyMirror = require('keymirror');

// Component-Global form states
var STATUS = keyMirror({
  DEFAULT: null,
  ATTEMPTING_LOGIN: null,
  ERROR_UNAUTHORIZED: null,
  GENERIC_ERROR: null
});

module.exports = {
  STATUS
};
