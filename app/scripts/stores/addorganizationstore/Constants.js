'use strict';
var keyMirror = require('keymirror');

// Component-Global form states
var STATUS = keyMirror({
  ATTEMPTING: null,
  DEFAULT: null,
  FACEPALM: null,
  BAD_REQUEST: null,
  SUCCESSFUL: null
});

module.exports = {
  STATUS
};
