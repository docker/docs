'use strict';

var debug = require('debug')('hub:actions:deleteEmailNotifs');
import async from 'async';

import { Notifications } from 'hub-js-sdk';

var deleteEmailNotifs = function(actionContext, {jwt, notificationID}) {
  var _delNotificationSettings = function(cb) {
    Notifications.deleteNotificationSubscription(jwt, notificationID, function(err, res) {
      if (err) {
        debug('deleteNotificationSubscription error', err);
        cb(err);
      } else if (res.ok) {//TODO:Weird that 204 no content is sent on success
        return cb(null, res);
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
    _delNotificationSettings,
    _getNotificationSettings
  ], function (err, results) {
    if (err) {
      actionContext.dispatch('SAVE_NOTIFICATIONS_ERROR');
      debug('final callback error', err);
    } else if (results[1]) {
      actionContext.dispatch('SAVE_NOTIFICATIONS_SUCCESS');
      actionContext.dispatch('RECEIVE_NOTIFICATIONS', results[1]);
    }
  });
};

module.exports = deleteEmailNotifs;
