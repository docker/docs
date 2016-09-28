'use strict';

import { Users } from 'hub-js-sdk';
var debug = require('debug')('hub:actions:convertToOrganization');

module.exports = function(actionContext, {jwt, username, newOwner}, done) {
  Users.convertToOrganization(jwt, username, newOwner, function(err, res) {
    if (err) {
      debug('error', err);
      if (res.badRequest){
        let message;
        if (res.body) {
          message = res.body;
        } else {
          try {
            message = JSON.parse(res.text);
          } catch(e) {
            message = {
              error: 'Your account could not be converted. Make sure you are not a member of another group and that the new owner username exists'
            };
          }
        }
        actionContext.dispatch('CONVERT_TO_ORG_BAD_REQUEST', message);
        done();
      } else {
        done();
      }
    } else {
      actionContext.history.push('/login/');
      actionContext.dispatch('LOGOUT', null);
    }
  });
};
