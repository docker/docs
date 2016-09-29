'use strict';
import {
  Users,
  JWT
} from 'hub-js-sdk';
const debug = require('debug')('hub:actions:saveSettingsData');

var saveSettingsData = function(actionContext, payload) {
  actionContext.dispatch('ACCOUNT_INFO_ATTEMPT_START');
  Users.updateUser(payload.JWT, payload.username, payload.updateData, function(err, res) {
    if (err) {
      debug(err);
      actionContext.dispatch('ACCOUNT_INFO_BAD_REQUEST', err.response.body);
    } else {
      actionContext.dispatch('ACCOUNT_INFO_SUCCESS');
      setTimeout(() => {actionContext.dispatch('ACCOUNT_INFO_STATUS_CLEAR');}, 4000);
      actionContext.dispatch('RECEIVE_USER', res.body);
    }
  });
};

module.exports = saveSettingsData;
