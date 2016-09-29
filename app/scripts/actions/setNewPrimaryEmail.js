'use strict';

import { sortByOrder } from 'lodash';
import { series, each } from 'async';
import {
  Emails
  } from 'hub-js-sdk';
var debug = require('debug')('hub:actions:setNewPrimaryEmail');

function updateEmailSettings({ dispatch },
                             {
                               JWT, username, emailId
                             },
                             done) {
  dispatch('START_SAVE_ACTION');
  series([
      function update(callback) {
          debug('emailObject', emailId);
          Emails.updateEmailByID(JWT,
            emailId,
            {primary: true},
            (err, res) => callback(null, res));
      }
    ],
    function afterUpdatingEmails(err, results) {
      if (err) {
        return debug(err);
      }
      Emails.getEmailsForUser(JWT, username, function(emailErr, res){
        if (emailErr) {
          debug('emailErr', emailErr);
          dispatch('FINISH_SAVE_ACTION');
          return done();
        }
        var emails = res.body.results;
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
