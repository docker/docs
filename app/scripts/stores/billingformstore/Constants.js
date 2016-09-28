'use strict';
var keyMirror = require('keymirror');

// Component-Global form states
var STATUS = keyMirror({
  DEFAULT: null,
  ATTEMPTING: null,
  SUCCESS: null,
  FORM_ERROR: null
});

module.exports = {
  STATUS
};
