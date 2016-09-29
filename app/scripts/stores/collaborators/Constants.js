'use strict';
var keyMirror = require('keymirror');

// Component-Global form states
var STATUS = keyMirror({
  DEFAULT: null,
  ATTEMPTING: null,
  SUCCESS: null,
  ERROR: null
});

module.exports = {
  STATUS
};
