'use strict';
const debug = require('debug')('navigate::serverTrialSuccess');
import merge from 'lodash/object/merge';
import find from 'lodash/collection/find';
import { Billing } from 'hub-js-sdk';

export default function serverTrialSuccess({
  actionContext, payload, done, maybeData
}){
  const { namespace } = payload.location.query;
  if(maybeData.token) {
    Billing.getLicensesForNamespace(maybeData.token, { namespace }, (err, res) => {
      if (err) {
        debug(err);
        if(err.response.badRequest) {
          const { detail } = err.response.body;
          if(detail) {
            actionContext.dispatch('RECEIVE_TRIAL_LICENSE_BAD_REQUEST', detail);
          }
        } else {
          actionContext.dispatch('RECEIVE_TRIAL_LICENSE_FACEPALM', err);
        }
        done();
      } else {
        let license = find(res.body.licenses, (obj) => {
          return obj.tier === 'Trial' || obj.alias === 'Trial';
        });
        license = merge({}, license, { namespace });
        actionContext.dispatch('RECEIVE_TRIAL_LICENSE', license);
        done();
      }
    });
  } else {
    // user must be logged in; they aren't
    const err = `You must be logged in to access your trial license`;
    actionContext.dispatch('RECEIVE_TRIAL_LICENSE_BAD_REQUEST', err);
    done();
  }
}
