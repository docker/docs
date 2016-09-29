'use strict';
var keyMirror = require('keymirror');

// Component-Global form states
var STATUS = keyMirror({
  DEFAULT: null,
  MEMBER_UNAUTHORIZED: null,
  TEAM_UNAUTHORIZED: null,
  MEMBER_ERROR: null,
  TEAM_ERROR: null,
  MEMBER_BAD_REQUEST: null,
  TEAM_BAD_REQUEST: null,
  GENERAL_SERVER_ERROR: null,
  CREATE_TEAM_SUCCESS: null,
  CREATE_MEMBER_SUCCESS: null,
  UPDATE_TEAM_ERROR: null,
  UPDATE_TEAM_SUCCESS: null
});

module.exports = {
  STATUS
};
