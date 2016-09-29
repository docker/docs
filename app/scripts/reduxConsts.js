'use strict';
//Consts used for redux actions
const keyMirror = require('keymirror');

export default keyMirror({
  //repos
  RECEIVE_REPO: null,

  //tags
  DELETE_REPO_TAG: null,
  RECEIVE_NAUTILUS_TAGS_FOR_REPOSITORY: null,
  RECEIVE_SCANNED_TAG_DATA: null,
  RECEIVE_TAGS_FOR_REPOSITORY: null,

  // statuses
  ATTEMPTING: null,
  ERROR: null,
  SUCCESS: null
});
