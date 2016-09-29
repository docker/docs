'use strict';
var keyMirror = require('keymirror');

// Component-Global form states
var STATUS = keyMirror({
  ATTEMPTING: null,
  DEFAULT: null,
  FORM_ERROR: null,
  SHOWING_CONFIRM_BOX: null
});

module.exports = {
  STATUS
};
