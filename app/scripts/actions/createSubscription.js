'use strict';
const request = require('superagent');
import { Billing } from 'hub-js-sdk';
/*global UpdateBillingInfoPayload */
import map from 'lodash/collection/map';
import has from 'lodash/object/has';
import merge from 'lodash/object/merge';
import isString from 'lodash/lang/isString';
import async from 'async';
var debug = require('debug')('hub:actions:createBillingSubscription');

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
    let message = 'There was an error creating your subscription. Please check your information and try again.';
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
    callback(null, { billforward_id });
  }
};

//--------------------------------------------------------------------------
// CREATE A BILLING PAYMENT METHOD ON BILLFORWAD WITH STRIPE
//--------------------------------------------------------------------------
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

function billingStripeCreatePaymentMethod(dispatch, {
  billforward_id,
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
      const accountID = billforward_id;
      const bfData = {
        '@type': 'StripeAuthCaptureRequest',
        accountID,
        cardID,
        defaultPaymentMethod: true,
        gateway: 'Stripe',
        stripeToken
      };
      _billforwardAddPaymentMethod(bfData, handleResponse({ callback: cb, dispatch, type: BILLING }));
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
// CREATE NEW BILLING ACCOUNT/PROFILE & SUBSCRIPTION
//--------------------------------------------------------------------------
function createSubscription(actionContext, {
  JWT,
  user,
  accountInfo,
  billingInfo,
  card,
  isNewBillingAccount,
  billforwardId: existingBfId,
  plan_code,
  package_code
}, done) {
  actionContext.dispatch('BILLING_SUBMIT_START');
  const {
    first_name,
    last_name,
    address1,
    address2,
    country,
    state,
    zip,
    city
    } = billingInfo;
  const {
    number,
    cvv,
    month,
    year,
    last_four,
    type,
    coupon_code,
    coupon
    } = card;

  let isOrg = false;
  let username = '';
  if (has(user, 'username')) {
    username = user.username;
  } else if (has(user, 'orgname')) {
    username = user.orgname;
    isOrg = true;
  }

  var subscriptionData = {
    first_name,
    last_name,
    email: accountInfo.email,
    username,
    plan: plan_code,
    package: package_code
  };
  if (coupon_code) {
    // coupon code is an optional field
    subscriptionData.coupon_code = coupon_code;
  }

  async.waterfall([
    function(callback) { // create billing profile
      const account = merge({}, accountInfo, {username});
      /*
      NOTE:
      Can only get to this action IF
      1) You don't have a billing profile account (Create the billing account)
      2) You don't have any billing payment information - but you have a billing profile account (Update the account)
      */
      if (isNewBillingAccount) {
        Billing.createBillingAccount(JWT, account, handleResponse({callback, dispatch: actionContext.dispatch, type: ACCOUNT}));
        // on success - callback(null, { billforward_id })
      } else {
        Billing.updateBillingAccount(JWT, username, account, handleResponse({callback, dispatch: actionContext.dispatch, type: ACCOUNT}));
      }
    },
    function({ billforward_id }, callback) {
      const bfId = billforward_id || existingBfId;
      if (bfId) {
        // NOTE: Save billing address information
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
        }, (err, res) => {
          if (err) {
            let message = 'There was an error saving your address information. Please check your information and try again.';
            if (res && res.body) {
              message = isString(res.body.detail) ? res.body.detail : res.body.message;
            }
            actionContext.dispatch('BILLING_SUBMIT_ERROR', message);
            callback(null);
            return;
          }
          callback(null, {...res.body, billforward_id });
        });
        return;
      }
      // NOTE: Unecessary for non-migrated accounts - Creation of billingprofile will update the account data
      callback(null, { billforward_id });
    },
    function({ billforward_id }, callback) {
      const bfId = billforward_id || existingBfId;
      if (bfId) {
        // NOTE: If we have a billforward id - create payment Method via stripe
        billingStripeCreatePaymentMethod(actionContext.dispatch, {
          billforward_id: bfId,
          cvc: cvv,
          exp_month: month,
          exp_year: year,
          name_first: first_name,
          name_last: last_name,
          number
        }, callback);
        return;
      }
      // NOTE: If we do NOT have a billforward id - tokenize via recurly
      const recurlyData = {
        first_name,
        last_name,
        address1,
        address2,
        country,
        state,
        zip,
        city,
        number,
        cvv,
        month,
        year,
        coupon_code
      };
      try {
        /**
         * This will only run where window is defined. ie: the browser
         * It throws an exception in node
         */
        window.recurly.configure(process.env.RECURLY_PUBLIC_KEY);
      } catch(e) {
        debug('error', e);
      }
      window.recurly.token(recurlyData, function(recurlyErr, token) {
        if (recurlyErr) {
          debug('recurly token error', recurlyErr.message);
          actionContext.dispatch('GET_RECURLY_ERROR', recurlyErr);
          callback(recurlyErr);
          return;
        }
        callback(null, { recurlyToken: token.id });
      });
    },
    function({ recurlyToken }, callback) {
      // NOTE: append the recurly token to the subscriptionData object;
      subscriptionData.payment_token = recurlyToken;
      /*
      NOTE: Creating billing subscription will either
      A) Use default card if billing account is migrated and we get a billforward id
      B) Subscription Data will have a recurly token if no billforward id
      IF Billing Info already exists:
      - Stripe: Adds a new default card
      - Recurly: we're posting with the recurly token which will update it
      */
      Billing.createBillingSubscription(JWT, subscriptionData, handleResponse({callback, dispatch: actionContext.dispatch, type: BILLING}));
    }
  ],
  function(err, res) {
    if (err) {
      debug(err);
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

module.exports = createSubscription;
