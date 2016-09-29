'use strict';

import { getToken, logout } from 'hub-js-sdk/src/Hub/SDK/Auth';
import { changePassword } from 'hub-js-sdk/src/Hub/SDK/Users';
import async from 'async';
import request from 'superagent';
var debug = require('debug')('hub:actions:changePassword');

var changePassAction = function({ dispatch, history },
                                {
                                  JWT, username, oldpassword, newpassword
                                }) {
  changePassword(
    JWT,
    {username, oldpassword, newpassword},
    (err, res) => {

      if (err) {
        debug('error', err);
        dispatch('RESET_PASSWORD_ERROR', res.body);
      } else if (res.ok) {
        dispatch('CHANGE_PASS_SUCCESS');
        async.parallel([
            function(callback) {
              logout(JWT, function(outErr, outRes) {
                if (outErr) {
                  debug('outErr', outErr);
                  dispatch('LOGOUT_ERROR', outErr);
                } else if (outRes.ok) {
                  dispatch('CHANGE_PASS_CLEAR', {});
                  history.pushState(null, '/account/password-reset-confirm/success/');
                  dispatch('LOGOUT');
                }
              });
            },
            function(callback) {
              request.post('/attempt-logout/')
                .end(callback);
            }],
          function(asyncErr, asyncRes) {

          });
      }
    }
  );
};

module.exports = changePassAction;
