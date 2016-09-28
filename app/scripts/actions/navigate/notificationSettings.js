'use strict';
const debug = require('debug')('navigate::notificationSettings');
import _ from 'lodash';
import async from 'async';
import {
  Emails, Notifications
  } from 'hub-js-sdk';

export default function notificationSettings({actionContext, payload, done, maybeData}){
  if (_.has(maybeData, 'token')) {
    var {token, user} = maybeData;
    async.parallel({
      subscriptions: function(callback) {
        Emails.getEmailSubscriptions(token, user.username, function(err, res){
          if (err) {
            debug(err);
            callback();
          } else {
            callback(null, res.body);
          }
        });
      },
      notifications: function(callback) {
        Notifications.getNotificationSubscriptions(token, function(err, res) {
          if (err) {
            debug(err);
            callback();
          } else {
            callback(null, res.body.results);
          }
        });
      },
      emails: function(callback) {
        Emails.getEmailsForUser(token, user.username, function(err, res){
          if (err) {
            debug(err);
            callback();
          } else {
            let emails = res.body.results;
            let sortedEmails = _.sortByOrder(emails, ['primary', 'verified'], [false, false]);
            callback(null, sortedEmails);
          }
        });
      }
    }, function(err, res) {
      let {subscriptions, notifications, emails} = res;
      var weeklyDigest, betaGroup;
      if (subscriptions) {
        weeklyDigest = subscriptions.DockerNewsMailingList;
        betaGroup = subscriptions.DockerBetaGroupMailingList;
        actionContext.dispatch('RECEIVE_EMAIL_SUBSCRIPTIONS', {weeklyDigest: weeklyDigest, betaGroup: betaGroup});
      }
      if (notifications) {
        actionContext.dispatch('RECEIVE_NOTIFICATIONS', notifications);
      }
      actionContext.dispatch('RECEIVE_EMAILS', {emails: emails});
      done();
    });
  } else {
    done();
  }
}
