'use strict';
var debug = require('debug')('navigate::home');
import async from 'async';
import _ from 'lodash';
import {
  Users
} from 'hub-js-sdk';

export default function enterpriseTrial({actionContext, payload, done, maybeData}){

  if(payload.location.query.partnervalue) {
    actionContext.dispatch('ENTERPRISE_PARTNER_RECEIVE_CODE', {
      code: payload.location.query.partnervalue
    });
  }

  if(_.has(maybeData, 'token')){
    Users.getNamespacesForUser(maybeData.token, function(err, res) {
      if (err) {
        done();
      } else {
        if(res.body && res.body.namespaces) {
          actionContext.dispatch('ENTERPRISE_TRIAL_RECEIVE_ORGS', res.body.namespaces);
        }
        done();
      }
    });
  } else {
    done();
  }
}
