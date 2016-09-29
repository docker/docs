'use strict';

import { sortByOrder } from 'lodash';
import {
  Emails
  } from 'hub-js-sdk';
var debug = require('debug')('hub:actions:addEmail');

module.exports = function(actionContext,
                          {
                            JWT, newEmail, user
                          },
                          done) {
  Emails.addEmailsForUser(JWT, user.username, newEmail, function(err, res) {
    if (err) {
      debug('error', res.body);
      actionContext.dispatch('ADD_EMAIL_INVALID', res.body.email);
    } else if (res.ok) {
      actionContext.dispatch('ADD_EMAIL_SUCCESS', res.body.email);
      Emails.getEmailsForUser(JWT, user.username, function(emailErr, emailRes){
        if (emailErr) {
          return done();
        }
        var emails = emailRes.body.results;
        var sortedEmails = sortByOrder(emails,
                                       ['primary', 'verified'],
                                       [false, false]);
        actionContext.dispatch('RECEIVE_EMAILS', {
          emails: sortedEmails
        });
      });
    }
  });
};
