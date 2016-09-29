'use strict';

import { Billing } from 'hub-js-sdk';
const debug = require('debug')('hub:actions:createNewLicenseProduction');

export default function(actionContext, {
  address1,
  address2,
  city,
  country,
  cvv,
  first_name,
  JWT,
  last_name,
  month,
  number,
  package_name,
  partnervalue,
  postal_code,
  state,
  username,
  year
}, done) {
  actionContext.dispatch('ENTERPRISE_PAID_ATTEMPT_START');
  try {
    /**
     * This will only run where window is defined. ie: the browser
     * It throws an exception in node
     */
    window.recurly.configure(process.env.RECURLY_PUBLIC_KEY);
  } catch(e) {
    debug('error', e);
  }

  var recurlyData = {
    first_name,
    last_name,
    month,
    year,
    cvv,
    number,
    address1,
    address2,
    city,
    country,
    postal_code,
    state
  };

  window.recurly.token(recurlyData, function(recurlyErr, token) {
    if (recurlyErr) {
      debug('recurly error', recurlyErr.fields);
      actionContext.dispatch('ENTERPRISE_PAID_GET_RECURLY_ERROR', recurlyErr.fields);
      done();
    } else {
      debug('creating License');
      Billing.createLicense(JWT, {
        username,
        package: package_name,
        payment_token: token.id,
        first_name,
        last_name,
        postal_code,
        partner_value: partnervalue
      }, (err, results) => {
        if(err) {
          if(err.response.badRequest) {
            const { detail } = err.response.body;
            if(detail) {
              actionContext.dispatch('ENTERPRISE_PAID_BAD_REQUEST', detail);
            }
          } else {
            actionContext.dispatch('ENTERPRISE_PAID_API_ERROR');
          }
        } else {
          actionContext.dispatch('ENTERPRISE_PAID_CLEAR_FORM');
          actionContext.history.push('/account/licenses/');
        }
      });
    }
  });
}
