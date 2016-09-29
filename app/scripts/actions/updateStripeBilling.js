'use strict';
/*global UpdateBillingInfoPayload */
const request = require('superagent');
import { Billing } from 'hub-js-sdk';
import has from 'lodash/object/has';
import merge from 'lodash/object/merge';
import map from 'lodash/collection/map';
import find from 'lodash/collection/find';
import isString from 'lodash/lang/isString';
import async from 'async';
var debug = require('debug')('hub:actions:updateBillingInformation');

import _encodeForm from '../components/utils/encodeForm.js';
import {
  ACCOUNT,
  BILLING,
  BILLFORWARD_ACCOUNT_ID,
  STRIPE_URL,
  STRIPE_STAGE_TOKEN,
  STRIPE_PROD_TOKEN,
  BF_STAGE_URL,
  BF_PROD_URL,
  BF_STAGE_TOKEN,
  BF_PROD_TOKEN,
  v4BillingProfile
} from 'stores/common/Constants.js';

const handleResponse = ({ callback, dispatch, type }) => (err, res) => {
  if (err) {
    let message;
    if (type === ACCOUNT) {
      message = 'There was an error updating your account profile information. Please check your information and try again.';
    } else if (type === BILLING) {
      message = 'There was an error updating your billing information. Please check your information and try again.';
    }
    if (res && res.body) {
      message = isString(res.body.detail) ? res.body.detail : res.body.message;
    }
    dispatch('BILLING_SUBMIT_ERROR', message);
    callback(err, res);
  } else if (res.ok) {
    if (type === ACCOUNT) {
      dispatch('BILLING_ACCOUNT_EXISTS');
    } else if (type === BILLING) {
      dispatch('BILLING_INFO_EXISTS');
    }
    const billforward_id = res.header[BILLFORWARD_ACCOUNT_ID];
    callback(null, { ...res.body, billforward_id });
  }
};

/*
  NOTE: Stripe's api only accepts x-www-form-urlencoded data
  NOTE: 'billforwardCreatePayment' is chained to this function and requires
    `billforward-id`. The `billforward-id` is being passed through to the
    create payment function via the META here.
*/
function _createCardToken({
  cvc,
  exp_month,
  exp_year,
  name_first,
  name_last,
  number
}, cb) {
  const stripeToken = process.env.ENV === 'production' ?
    STRIPE_PROD_TOKEN : STRIPE_STAGE_TOKEN;
  const card = {
    name: `${name_first} ${name_last}`,
    cvc,
    number,
    exp_month,
    exp_year
  };
  const encoded = _encodeForm({ card });
  request.post(STRIPE_URL)
    .accept('application/json')
    .type('application/x-www-form-urlencoded')
    .set('Authorization', 'Bearer ' + stripeToken)
    .send(encoded)
    .end(cb);
}

/*
  NOTE: This is the call to billforward that actually adds a payment method
  to a user's billing profile.
  This call requires a billforward-id (accountID) which is DIFFERENT than the
  docker_id - Hence why 'billingCreatePaymentMethod' is NOT being wrapped by
  the getAccountFromNamespace decorator.
*/
function _billforwardAddPaymentMethod({
  '@type': type,
  accountID,
  cardID,
  defaultPaymentMethod,
  gateway,
  stripeToken
}, cb) {
  const billforwardUrl = process.env.ENV === 'production' ?
    BF_PROD_URL : BF_STAGE_URL;
  const billforwardToken = process.env.ENV === 'production' ?
    BF_PROD_TOKEN : BF_STAGE_TOKEN;

  request.post(billforwardUrl)
        .accept('application/json')
        .type('application/json')
        .set('Authorization', 'Bearer ' + billforwardToken)
        .send({
          '@type': type,
          accountID,
          cardID,
          defaultPaymentMethod,
          gateway,
          stripeToken
        })
        .end(cb);
}

//--------------------------------------------------------------------------
// CREATE A BILLING PAYMENT METHOD ON BILLFORWAD WITH STRIPE
// DUPLICATED FROM ./createStripeSubscription.js
//--------------------------------------------------------------------------
function billingCreatePaymentMethod(dispatch, {
  billforwardId,
  cvc,
  exp_month,
  exp_year,
  name_first,
  name_last,
  number
}, cb) {
  const cardData = {
    cvc,
    exp_month,
    exp_year,
    name_first,
    name_last,
    number
  };
  /*
  NOTE:
  Creating a payment method requires 2 parts
  1 - Generating a token from Stripe's api
  2 - Sending generated token to Billforward's api to attach payment method to
      a relevant billing profile.
  */
  _createCardToken(cardData, (stripeErr, stripeRes) => {
    if (!stripeRes.ok) {
      let message = 'There was an error submitting your card information. Please check your information and try again.';
      if (stripeRes.error && stripeRes.error.message) {
        message = stripeRes.error.message;
      }
      dispatch('BILLING_SUBMIT_ERROR', message);
      cb(message);
    } else {
      const tokenObject = stripeRes && stripeRes.body;
      const stripeToken = tokenObject.id;
      const cardID = tokenObject.card.id;
      const accountID = billforwardId;
      const bfData = {
        '@type': 'StripeAuthCaptureRequest',
        accountID,
        cardID,
        defaultPaymentMethod: true,
        gateway: 'Stripe',
        stripeToken
      };
      _billforwardAddPaymentMethod(bfData, handleResponse({callback: cb, dispatch, type: BILLING}));
    }
  });
}
//--------------------------------------------------------------------------
// SAVE ADDRESS INFORMATION IN BILLFORWARD
//--------------------------------------------------------------------------
const updateV4BillingProfile = ({ JWT, profileData, docker_id }, cb) => {
  const v4BillingAPI = v4BillingProfile(docker_id);
  request
    .patch(v4BillingAPI)
    .accept('application/json')
    .type('application/json')
    .set('Authorization', 'Bearer ' + JWT)
    .send(profileData)
    .end(cb);
};

//--------------------------------------------------------------------------
// UPDATE STRIPE BILLING INFORMATION
// NOTE: This will just add a payment method - creating multiple cards
//--------------------------------------------------------------------------
function updateBillingInformation(actionContext, {
  JWT,
  user,
  billingInfo,
  accountInfo,
  card,
  billforwardId
}, done) {
  actionContext.dispatch('BILLING_SUBMIT_START');
  let isOrg = false;
  let username = '';
  if (has(user, 'username')) {
    username = user.username;
  } else if (has(user, 'orgname')) {
    username = user.orgname;
    isOrg = true;
  }
  const {
    email,
    first_name,
    last_name,
    company_name
  } = accountInfo;
  const {
    address1,
    address2,
    country,
    state,
    zip,
    city,
    first_name: billing_first,
    last_name: billing_last
  } = billingInfo;
  const {
    number,
    cvv,
    month,
    year
    } = card;

  async.series([
    function(callback){
      const account = merge({}, accountInfo, {user: username});
      Billing.updateBillingAccount(JWT, username, account, handleResponse({callback, dispatch: actionContext.dispatch, type: ACCOUNT}));
    },
    function(callback) {
      /*
      NOTE:
      THIS WILL CREATE A NEW PAYMENT METHOD ON EVERY CALL
      - In billforward this could potentially save multiple versions of the same card
      - Solution would be to check whether the card is the same and decide to create or not
      */
      billingCreatePaymentMethod(actionContext.dispatch, {
        billforwardId,
        cvc: cvv,
        exp_month: month,
        exp_year: year,
        name_first: billing_first,
        name_last: billing_last,
        number
      }, callback);
    },
    function(callback) {
      const profileData = {
        first_name,
        last_name,
        addresses: [
          {
            address_line_1: address1,
            address_line_2: address2,
            city,
            province: state,
            country,
            post_code: zip,
            primary_address: true
          }
        ]
      };
      updateV4BillingProfile({
        JWT,
        profileData,
        docker_id: user.id
      }, handleResponse({callback, dispatch: actionContext.dispatch, type: ACCOUNT}));
    }
  ],
  function(err, res) {
    if (err) {
      done();
    } else {
      if (isOrg) {
        actionContext.history.push(`/u/${username}/dashboard/billing/`);
      } else {
        actionContext.history.push('/account/billing-plans/');
      }
      actionContext.dispatch('BILLING_SUBMIT_SUCCESS');
      done();
    }
  });


}

module.exports = updateBillingInformation;
