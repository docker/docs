'use strict';
var keyMirror = require('keymirror');

// Component-Global form states
var STATUS = keyMirror({
  DEFAULT: null,
  REPO_ALREADY_EXISTS: null,
  PRIVATE_REPO_QUOTA_EXCEEDED: null,
  BAD_REQUEST: null,
  REPO_NOT_FOUND: null
});

module.exports = {
  STATUS
};
