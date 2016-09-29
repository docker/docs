/* @flow */
'use strict';

import type FluxibleActionContext from '../../../flow-libs/fluxible';

import { Auth } from 'hub-js-sdk';
import request from 'superagent';
import async from 'async';
const debug = require('debug')('hub:actions:logout');

module.exports = function(actionContext: FluxibleActionContext, jwt) {
  async.parallel([
    function(callback) {
      Auth.logout(jwt, function(err, res) {
        if (err) {
          debug('error', err);
          actionContext.dispatch('LOGOUT_ERROR', err);
        } else if (res.ok) {
          actionContext.dispatch('LOGOUT');
          actionContext.history.push('/');
          }
      });
    },
    function(callback) {
      request.post('/attempt-logout/')
      .end(callback);
    }],
                 function(err, results) {});
};
