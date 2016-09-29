'use strict';

import { parallel } from 'async';
import isEmpty from 'lodash/lang/isEmpty';
import {
  Emails
  } from 'hub-js-sdk';
var debug = require('debug')('hub:actions:saveOutbound');

module.exports = function({ dispatch },
                          {
                            JWT,
                            username,
                            weeklyDigest,
                            betaGroup
                          }) {
  /*eslint-disable camelcase */
  var digestSubscribe = weeklyDigest.subscribed_emails;
  var digestUnsubscribe = weeklyDigest.unsubscribed_emails;
  var betaSubscribe = betaGroup.subscribed_emails;
  var betaUnsubscribe = betaGroup.unsubscribed_emails;
  /*eslint-enable camelcase */
  parallel({
    digestSubscribed: function (callback) {
      Emails.subscribeEmails(JWT, username, {
        /*eslint-disable camelcase */
        subscribed_emails: digestSubscribe,
        mailing_list: 'DockerNewsMailingList'
        /*eslint-enable camelcase*/
      }, function (err, res) {
        if (err) {
          callback(err);
        } else {
          callback(null, res);
        }
      });
    },
    betaSubscribed: function (callback) {
      Emails.subscribeEmails(JWT, username, {
        /*eslint-disable camelcase */
        subscribed_emails: betaSubscribe,
        mailing_list: 'DockerBetaGroupMailingList'
        /*eslint-enable camelcase */
      }, function (err, res) {
        if (err) {
          callback(err);
        } else {
          callback(null, res);
        }
      });
    },
    digestUnsubscribed: function (callback) {
      if (!isEmpty(digestUnsubscribe)) {
        Emails.unsubscribeEmails(JWT, username, {
          /*eslint-disable camelcase */
          unsubscribed_emails: digestUnsubscribe,
          mailing_list: 'DockerNewsMailingList'
          /*eslint-enable camelcase */
        }, function (err, res) {
          if (err) {
            callback(err);
          } else {
            callback(null, res.body.unsubscribed);
          }
        });
      } else {
        callback(null);
      }
    },
    betaUnsubscribed: function (callback) {
      if (!isEmpty(betaUnsubscribe)) {
        Emails.unsubscribeEmails(JWT, username, {
          /*eslint-disable camelcase */
          unsubscribed_emails: betaUnsubscribe,
          mailing_list: 'DockerBetaGroupMailingList'
          /*eslint-enable camelcase */
        }, function (err, res) {
          if (err) {
            callback(err);
          } else {
            callback(null, res.body.unsubscribed);
          }
        });
      } else {
       callback(null);
      }
    }
    },
    function(err, res) {
      if(err) {
        debug(err);
        dispatch('SAVE_OUTBOUND_ERROR');
      } else {
        Emails.getEmailSubscriptions(JWT, username, function(error, response){
          if (error) {
            return debug(error);
          }
          dispatch('RECEIVE_EMAIL_SUBSCRIPTIONS', {
            weeklyDigest: response.body.DockerNewsMailingList,
            betaGroup: response.body.DockerBetaGroupMailingList
          });
          dispatch('SAVE_OUTBOUND_SUCCESS');
        });
      }
    });
};
