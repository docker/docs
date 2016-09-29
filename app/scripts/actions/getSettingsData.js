'use strict';
import _ from 'lodash';
import {
  Users
  } from 'hub-js-sdk';
var debug = require('debug')('hub:actions:getSettingsData');

var getSettingsData = function(actionContext, {JWT, username, repoType}) {
  Users.getUserSettings(JWT, username, function(err, res) {
    if (err){
      debug('getUserSettings', err);
    } else {
      if (repoType === 'regular') {
        actionContext.dispatch('CREATE_REPO_UPDATE_FIELD_WITH_VALUE', {
          fieldKey: 'is_private',
          fieldValue: res.body.default_repo_visibility
        });
      } else if (repoType === 'autobuild') {
        actionContext.dispatch('AUTOBUILD_FORM_UPDATE_FIELD_WITH_VALUE', {
          fieldKey: 'isPrivate',
          fieldValue: res.body.default_repo_visibility
        });
      }
    }
  });
};

module.exports = getSettingsData;
