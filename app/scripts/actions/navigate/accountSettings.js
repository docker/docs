'use strict';
const debug = require('debug')('navigate::accountSettings');
import {
  Emails,
  Users
} from 'hub-js-sdk';
import async from 'async';
import _ from 'lodash';

function sortEmails(emails) {
  return _.sortByOrder(emails,
                       ['primary', 'verified'],
                       [false, false]);
}

function parseRes({res, actionContext}) {
  let emails = res.body.results;
  let sortedEmails = _.sortByOrder(emails,
                                  ['primary', 'verified'],
                                  [false, false]);
  return sortedEmails;
}

export default function accountSettings({
  actionContext, payload, done, maybeData
}){
  actionContext.dispatch('CHANGE_PASS_CLEAR');
  if (_.has(maybeData, 'token')) {
    var { token, user } = maybeData;
    async.parallel({
      emails: function(callback) {
        debug('ACCOUNT SETTINGS EMAILS');
        if (user && user.isAdmin) {
          Emails.getEmailsForUser(token, user.username, function(err, res){
            if (err) {
              callback();
            } else {
              let emails = parseRes({res, actionContext});
              actionContext.dispatch('RECEIVE_EMAILS', {emails: emails});
              callback(null, emails);
            }
          });
        } else {
          Emails.getEmailsJWT(token, function(err, res){
            if (err) {
              callback();
            } else {
              let emails = parseRes({res, actionContext});
              actionContext.dispatch('RECEIVE_EMAILS', {emails: emails});
              callback(null, emails);
            }
          });
        }
      },
      repoStats: function(callback) {
        debug('GET REPO STATS');
        Users.getUserSettings(token, user.username, function(err, res) {
          if (err) {
            callback();
          } else {
            actionContext.dispatch('RECEIVE_PRIVATE_REPOSTATS', res.body);
            callback(null, res.body);
          }
        });
      }
    }, function(err, res) {
      if (err) {
        debug('ERROR', err);
        done();
      } else {
        let { emails, repoStats } = res;
        debug('SUCCESS');
        done();
      }
    });
  } else {
    done();
  }
}
