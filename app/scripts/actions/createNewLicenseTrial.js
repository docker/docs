'use strict';

import { Billing } from 'hub-js-sdk';
const debug = require('debug')('hub:actions:createNewLicenseTrial');

export default function(actionContext, {
  JWT,
  companyName,
  country,
  email,
  firstName,
  jobFunction,
  lastName,
  namespace,
  packageName,
  phoneNumber,
  state
}) {
  actionContext.dispatch('ENTERPRISE_TRIAL_ATTEMPT_START');
  const licenseData = {
    company_name: companyName,
    country,
    email,
    first_name: firstName,
    job_function: jobFunction,
    last_name: lastName,
    package: packageName,
    phone_number: phoneNumber,
    username: namespace
  };

  // state is an optional field
  if(state) {
      licenseData.state = state;
  }

  Billing.createLicense(JWT, licenseData, (err, results) => {
   if(err) {
     debug('error', err);
     if(err.response.badRequest) {
      const { detail } = err.response.body;
      if(detail) {
        actionContext.dispatch('ENTERPRISE_TRIAL_BAD_REQUEST', detail);
      }
     } else {
       actionContext.dispatch('ENTERPRISE_TRIAL_FACEPALM');
     }
   } else {
      actionContext.dispatch('ENTERPRISE_TRIAL_SUCCESS');
      actionContext.history.push(`/enterprise/trial/success/?namespace=${namespace}&step=1`);
   }
  });
}
