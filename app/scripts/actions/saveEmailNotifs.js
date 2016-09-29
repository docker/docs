'use strict';

var debug = require('debug')('hub:actions:saveEmailNotifs');
import async from 'async';

import { Notifications } from 'hub-js-sdk';

var saveEmailNotifs = function(actionContext, {jwt, notification}) {
  actionContext.dispatch('START_SAVE_ACTION');
  var _saveNotificationSettings = function(cb) {
    Notifications.setNotificationSubscription(jwt, notification, function(err, res) {
      if (err) {
        debug('error', err);
        cb(err);
      } else if (res.body && res.ok) {
        return cb(null, res.body);
      }
    });
  };

  var _getNotificationSettings = function(cb) {
    //all good, do the next call to get the notifications
    Notifications.getNotificationSubscriptions(jwt, function(err, res) {
      if (err) {
        debug('getNotificationSubscriptions error', err);
        cb(err);
      } else {
        cb(null, res.body.results);
      }
    });
  };

  async.series([
    function (callback) {
      _saveNotificationSettings(callback);
    },
    function (callback) {
      _getNotificationSettings(callback);
    }
  ], function (err, results) {
    if (err) {
      debug('async.series callback error', err);
      actionContext.dispatch('SAVE_NOTIFICATIONS_ERROR');
    } else if (results[1]) {
      actionContext.dispatch('SAVE_NOTIFICATIONS_SUCCESS');
      actionContext.dispatch('RECEIVE_NOTIFICATIONS', results[1]);
    }
  });
};

module.exports = saveEmailNotifs;
