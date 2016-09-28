'use strict';
var keyMirror = require('keymirror');

// Component-Global form states
export const STATUS = keyMirror({
  ATTEMPTING: null,
  DEFAULT: null,
  FACEPALM: null,
  SUCCESSFUL: null,
  ERROR: null
});

export const ACCOUNT = 'account';
export const BILLING = 'billing';
export const STRIPE_URL = 'https://api.stripe.com/v1/tokens';
export const STRIPE_STAGE_TOKEN = 'pk_test_DMJYisAqHlWvFPgRfkKayAcF';
export const STRIPE_PROD_TOKEN = 'pk_live_89IjovLdwh2MTzV7JsGJK3qk';
export const BF_STAGE_URL = 'https://api-sandbox.billforward.net:443/v1/tokenization/auth-capture';
export const BF_PROD_URL = 'https://api.billforward.net/v1/tokenization/auth-capture';
export const BF_STAGE_TOKEN = 'ec687f76-c1b6-4d71-b919-4fe99202ca13';
export const BF_PROD_TOKEN = '650cbe35-4aca-4820-a7d1-accec8a7083a';
export const BILLFORWARD_ACCOUNT_ID = 'billforward-account-id';
export const v4BillingProfile = (docker_id) => {
  return `/api/billing/v4/accounts/${docker_id}/profile`;
};
