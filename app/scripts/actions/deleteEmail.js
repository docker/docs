'use strict';

import sortByOrder from 'lodash/collection/sortByOrder';
import { series, each } from 'async';
import {
  Emails
  } from 'hub-js-sdk';
var debug = require('debug')('hub:actions:deleteEmail');

function updateEmailSettings({ dispatch },
                             {JWT, delEmailID, username},
                             done) {
  dispatch('START_SAVE_ACTION');
  Emails.deleteEmailByID(JWT,
    delEmailID,
    (err, res) => {
      if (err) {
        return debug('deleteEmailByID error', err);
      }
      Emails.getEmailsForUser(JWT, username, function(emailErr, emailRes){
        if (emailErr) {
          debug('getEmailsForUser error', emailErr);
          dispatch('FINISH_SAVE_ACTION');
          return done();
        }
        var emails = emailRes.body.results;
        var sortedEmails = sortByOrder(emails,
          ['primary', 'verified'],
          [false, false]);
        dispatch('RECEIVE_EMAILS', {
          emails: sortedEmails
        });
        dispatch('FINISH_SAVE_ACTION');
      });
    });
}

module.exports = updateEmailSettings;
