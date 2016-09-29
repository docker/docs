'use strict';

import _ from 'lodash';
import {
  Users
  } from 'hub-js-sdk';
var debug = require('debug')('navigate::orgSummary');

export default function orgs({actionContext, payload, done, maybeData}){
  debug('Organization Settings Navigate Token -> ' + maybeData.token);
  if (_.has(maybeData, 'token')) {
    Users.getOrgsForUser(maybeData.token, function(err, res) {
      if (err) {
        debug(err);
        done();
      } else {
        actionContext.dispatch('CURRENT_USER_ORGS', res.body);
        done();
      }
    });
  } else {
    done();
  }
}
