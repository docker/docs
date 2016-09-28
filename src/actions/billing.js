const request = require('superagent-promise')(require('superagent'), Promise);
import { merge } from 'lodash';
import isStaging from 'lib/utils/isStaging';
import isDev from 'lib/utils/isDevelopment';
import encodeForm from 'lib/utils/encodeForm';
import createBearer from 'lib/utils/create-bearer';
import { bearer } from 'lib/utils/authHeaders';
import bugsnagNotify from 'lib/utils/metrics';
import { isProductBundle } from 'lib/utils/product-utils';
import {
  marketplaceFetchBundleDetail,
  marketplaceFetchRepositoryDetail,
} from 'actions/marketplace';

/* eslint-disable max-len */
export const BILLING_CREATE_PAYMENT_METHOD = 'BILLING_CREATE_PAYMENT_METHOD';
export const BILLING_CREATE_PAYMENT_TOKEN = 'BILLING_CREATE_PAYMENT_TOKEN';
export const BILLING_CREATE_PRODUCT = 'BILLING_CREATE_PRODUCT';
export const BILLING_CREATE_PROFILE = 'BILLING_CREATE_PROFILE';
export const BILLING_CREATE_SUBSCRIPTION = 'BILLING_CREATE_SUBSCRIPTION';
export const BILLING_DELETE_PAYMENT_METHOD = 'BILLING_DELETE_PAYMENT_METHOD';
export const BILLING_DELETE_SUBSCRIPTION = 'BILLING_DELETE_SUBSCRIPTION';
export const BILLING_FETCH_INVOICE_PDF = 'BILLING_FETCH_INVOICE_PDF';
export const BILLING_FETCH_INVOICES = 'BILLING_FETCH_INVOICES';
export const BILLING_FETCH_LICENSE_DETAIL = 'BILLING_FETCH_LICENSE_DETAIL';
export const BILLING_FETCH_LICENSE_FILE = 'BILLING_FETCH_LICENSE_FILE';
export const BILLING_FETCH_PAYMENT_METHODS = 'BILLING_FETCH_PAYMENT_METHODS';
export const BILLING_FETCH_PRODUCT = 'BILLING_FETCH_PRODUCT';
export const BILLING_FETCH_PROFILE_SUBSCRIPTIONS = 'BILLING_FETCH_PROFILE_SUBSCRIPTIONS';
export const BILLING_FETCH_PROFILE_SUBSCRIPTIONS_AND_PRODUCTS = 'BILLING_FETCH_PROFILE_SUBSCRIPTIONS_AND_PRODUCTS';
export const BILLING_FETCH_PROFILE = 'BILLING_FETCH_PROFILE';
export const BILLING_SET_DEFAULT_PAYMENT_METHOD = 'BILLING_SET_DEFAULT_PAYMENT_METHOD';
export const BILLING_UPDATE_PROFILE = 'BILLING_UPDATE_PROFILE';
export const BILLING_UPDATE_SUBSCRIPTION = 'BILLING_UPDATE_SUBSCRIPTION';

const BILLING_HOST = '';
export const BILLING_BASE_URL = `${BILLING_HOST}/api/billing/v4`;

/*
 * =============================================================================
 * BILLING PRODUCT CALLS
 * -  Calls associated with the product catalog & plans
 * =============================================================================
 */
export function billingFetchProduct({ id }) {
  const url = `${BILLING_BASE_URL}/products/${id}/`;
  return {
    type: BILLING_FETCH_PRODUCT,
    meta: { id },
    payload: {
      promise:
        request
          .get(url)
          .set(bearer())
          .end()
          .then((res) => res.body),
    },
  };
}

export const billingCreateProduct = ({ id, body }) => dispatch => {
  const url = `${BILLING_BASE_URL}/products/${id}/`;
  return dispatch({
    type: BILLING_CREATE_PRODUCT,
    meta: { id },
    payload: {
      promise:
        request
          .put(url)
          .accept('application/json')
          .set(bearer())
          .send(body)
          .end()
          .then((res) => res.body),
    },
  });
};


/*
 * =============================================================================
 * USER BILLING CALLS
 * - Calls associated with a User's billing/subscription information
 * - Must include JWT and docker_id
 * =============================================================================
 */

//------------------------------------------------------------------------------
// PAYMENTS & PAYMENT METHODS
//------------------------------------------------------------------------------

export function billingFetchPaymentMethods({ docker_id }) {
  const url = `${BILLING_BASE_URL}/accounts/${docker_id}/payment-methods/`;
  return {
    type: BILLING_FETCH_PAYMENT_METHODS,
    meta: {
      docker_id,
      shouldRedirectToLogin: true,
    },
    payload: {
      promise:
        request
          .get(url)
          .set(bearer())
          .end()
          .then((res) => res.body),
    },
  };
}

/*
  NOTE: Stripe's api only accepts x-www-form-urlencoded data
  NOTE: 'billforwardCreatePayment' is chained to this function and requires
    `billforward-id`. The `billforward-id` is being passed through to the
    create payment function via the META here.
*/
function createCardToken({
  cvc,
  exp_month,
  exp_year,
  name_first,
  name_last,
  number,
}, meta) {
  // TODO ENV VAR - nathan - move const variables to environment
  const stripeUrl = 'https://api.stripe.com/v1/tokens';
  let stripeToken = 'pk_live_89IjovLdwh2MTzV7JsGJK3qk';
  if (isStaging() || isDev()) {
    stripeToken = 'pk_test_DMJYisAqHlWvFPgRfkKayAcF';
  }
  const card = {
    name: `${name_first} ${name_last}`,
    cvc,
    number,
    exp_month,
    exp_year,
  };
  const encoded = encodeForm({ card });
  return {
    type: BILLING_CREATE_PAYMENT_TOKEN,
    payload: {
      promise: request.post(stripeUrl)
        .accept('application/json')
        .type('application/x-www-form-urlencoded')
        .set('Authorization', createBearer(stripeToken))
        .send(encoded)
        .then((res) => {
          return {
            body: res.body,
            meta,
          };
        }, (res) => {
          /*
          402 PAYMENT REQUIRED error occurs when user has incorrectly input
          their card information. We don't care about user error - only if
          something goes wrong on our end.
          */
          if (res.status !== 402 || !isDev()) {
            bugsnagNotify('STRIPE TOKEN ERR', res.message);
          }
        }),
    },
  };
}

/*
  NOTE: This is the call to billforward that actually adds a payment method
  to a user's billing profile.
  This call requires a billforward-id (accountID) which is DIFFERENT than the
  docker_id - Hence why 'billingCreatePaymentMethod' is NOT being wrapped by
  the getAccountFromNamespace decorator.
*/
function billforwardCreatePayment({
  '@type': type,
  accountID,
  cardID,
  defaultPaymentMethod,
  gateway,
  stripeToken,
}) {
  // TODO ENV VAR - nathan - move const variables to environment
  let billforwardUrl =
    'https://api.billforward.net/v1/tokenization/auth-capture';
  let billforwardToken = '650cbe35-4aca-4820-a7d1-accec8a7083a';
  if (isStaging() || isDev()) {
    billforwardUrl =
      'https://api-sandbox.billforward.net:443/v1/tokenization/auth-capture';
    billforwardToken = 'ec687f76-c1b6-4d71-b919-4fe99202ca13';
  }
  return {
    type: BILLING_CREATE_PAYMENT_METHOD,
    payload: {
      promise: request.post(billforwardUrl)
        .accept('application/json')
        .type('application/json')
        .set('Authorization', createBearer(billforwardToken))
        .send({
          '@type': type,
          accountID,
          cardID,
          defaultPaymentMethod,
          gateway,
          stripeToken,
        })
        .then((res) => res.body.results, (res) => {
          if (!isDev()) {
            bugsnagNotify('BF AUTH CAPTURE ERR', res.message);
          }
        }),
    },
  };
}

export const billingCreatePaymentMethod = ({
  billforward_id,
  cvc,
  exp_month,
  exp_year,
  name_first,
  name_last,
  number,
}) => dispatch => {
  const cardData = {
    cvc,
    exp_month,
    exp_year,
    name_first,
    name_last,
    number,
  };
  /*
  NOTE:
  Creating a payment method requires 2 parts
  1 - Generating a token from Stripe's api
  2 - Sending generated token to Billforward's api to attach payment method to
      a relevant billing profile.
  */
  return dispatch(createCardToken(cardData, { billforward_id }))
    .then((res) => {
      const tokenObject = res.value.body;
      const stripeToken = tokenObject.id;
      const cardID = tokenObject.card.id;
      const accountID = res.value.meta.billforward_id;
      const bfData = {
        '@type': 'StripeAuthCaptureRequest',
        accountID,
        cardID,
        defaultPaymentMethod: true,
        gateway: 'Stripe',
        stripeToken,
      };
      return dispatch(billforwardCreatePayment(bfData));
    });
};

// Note: No response from this endpoint. Should refetch payment methods in then
export const billingSetDefaultPaymentMethod = ({ docker_id, card_id }) => {
  const url =
    `${BILLING_BASE_URL}/accounts/${docker_id}/payment-methods/${card_id}/`;
  return {
    type: BILLING_SET_DEFAULT_PAYMENT_METHOD,
    payload: {
      promise: request
        .patch(url)
        .set(bearer())
        .send({
          default: true,
        })
        .end()
        .then((res) => res.body),
    },
  };
};

// Note: No response from this endpoint. Should refetch payment methods in then
export const billingDeletePaymentMethod = ({ docker_id, card_id }) => {
  const url =
    `${BILLING_BASE_URL}/accounts/${docker_id}/payment-methods/${card_id}/`;
  return {
    type: BILLING_DELETE_PAYMENT_METHOD,
    payload: {
      promise: request
        .del(url)
        .set(bearer())
        .end()
        .then((res) => res.body),
    },
  };
};

//------------------------------------------------------------------------------
// PROFILES
//------------------------------------------------------------------------------
export const billingFetchProfile = ({ docker_id, isOrg }) => dispatch => {
  const url = `${BILLING_BASE_URL}/accounts/${docker_id}/`;
  return dispatch({
    type: BILLING_FETCH_PROFILE,
    meta: { docker_id, isOrg },
    payload: {
      promise:
        request
          .get(url)
          .set(bearer())
          .end()
          .then((res) => {
            /*
            NOTE: required billforward-account-id must be pulled from the header
            */
            return merge(
              {},
              res.body,
              { profile:
                { billforward_id: res.header['billforward-account-id'] },
              }
            );
          }),
    },
  });
};

/* NOTE: A sample profile includes information such as...
 * profile: {
 *   first_name,
 *   last_name,
 *   email,
 *   primary_phone,
 *   company_name,
 *   job_function,
 *   addresses,
 * }
 */
export function billingCreateProfile({
  docker_id,
  profile,
}) {
  const url = `${BILLING_BASE_URL}/accounts/${docker_id}/`;
  return {
    type: BILLING_CREATE_PROFILE,
    meta: { docker_id },
    payload: {
      promise:
        request
          .put(url)
          .set(bearer())
          .send({ profile })
          .end()
          .then((res) => {
            /*
            NOTE: required billforward-account-id must be pulled from the header
            */
            return merge(
              {},
              res.body,
              { profile:
                { billforward_id: res.header['billforward-account-id'] },
              }
            );
          }),
    },
  };
}

export function billingUpdateProfile({
  addresses,
  company_name,
  docker_id,
  email,
  first_name,
  job_function,
  last_name,
  phone_primary,
}) {
  const url = `${BILLING_BASE_URL}/accounts/${docker_id}/profile/`;
  return {
    type: BILLING_UPDATE_PROFILE,
    meta: { docker_id },
    payload: {
      promise:
        request
          .patch(url)
          .set(bearer())
          .send({
            addresses,
            company_name,
            email,
            first_name,
            last_name,
            job_function,
            phone_primary,
          })
          .end()
          .then((res) => res.body),
    },
  };
}

//------------------------------------------------------------------------------
// INVOICES
//------------------------------------------------------------------------------
export function billingFetchInvoices({ docker_id }) {
  const url = `${BILLING_BASE_URL}/invoices/`;
  return {
    type: BILLING_FETCH_INVOICES,
    meta: { shouldRedirectToLogin: true },
    payload: {
      promise:
        request
        .get(url)
        .query({ docker_id })
        .set(bearer())
        .end()
        .then((res) => res.body),
    },
  };
}

export function billingFetchInvoicePDF({ docker_id, invoice_id }) {
  const url = `${BILLING_BASE_URL}/invoices/${invoice_id}`;
  // Undocumented responseType property
  // https://github.com/visionmedia/superagent/pull/888/files
  return {
    type: BILLING_FETCH_INVOICE_PDF,
    payload: {
      promise:
        request
        .get(url)
        .query({ docker_id })
        .set(bearer())
        .accept('application/pdf')
        .responseType('blob')
        .end()
        .then(({ xhr }) => {
          // xhr.response is a Blob
          return xhr.response;
        }),
    },
  };
}

//------------------------------------------------------------------------------
// LICENSES
//------------------------------------------------------------------------------
export const billingFetchLicenseDetail = ({ subscription_id }) => {
  const url = `${BILLING_BASE_URL}/subscriptions/${subscription_id}/license-detail/`;
  return {
    type: BILLING_FETCH_LICENSE_DETAIL,
    meta: { subscription_id },
    payload: {
      promise:
        request
          .get(url)
          .set(bearer())
          .end()
          .then((res) => res.body),
    },
  };
};

export const billingFetchLicenseFile = ({ subscription_id }) => {
  const url = `${BILLING_BASE_URL}/subscriptions/${subscription_id}/license-file/`;
  return {
    type: BILLING_FETCH_LICENSE_FILE,
    meta: { subscription_id },
    payload: {
      promise:
        request
          .get(url)
          .set(bearer())
          .end()
          .then((res) => res.body),
    },
  };
};

//------------------------------------------------------------------------------
// SUBSCRIPTIONS
//------------------------------------------------------------------------------
// TODO Kristie Replace this with the composition API when it is available
export const billingFetchProfileSubscriptionsAndProducts = ({ docker_id }) => dispatch => {
  const url = `${BILLING_BASE_URL}/subscriptions/`;
  return dispatch({
    type: BILLING_FETCH_PROFILE_SUBSCRIPTIONS_AND_PRODUCTS,
    payload: {
      promise:
        request
          .get(url)
          .query({ docker_id })
          .set(bearer())
          .end()
          .then((res) => {
            const results = res.body || [];
            // For each subscription, we must fetch the product details
            // from billing AND from product catalog
            const promises = results.map((sub) => {
              const { product_id: id, subscription_id } = sub;
              const getProductDetails = [
                // Do not make this call bubble up the err
                dispatch(billingFetchProduct({ id })).catch(() => {}),
              ];
              if (isProductBundle(id)) {
                getProductDetails.push(
                  dispatch(marketplaceFetchBundleDetail({ id }))
                );
                // Fetch license details for DDC
                getProductDetails.push(
                  dispatch(billingFetchLicenseDetail({ subscription_id }))
                    .catch(() => {}) // Do not make this call bubble up the err
                );
              } else {
                getProductDetails.push(
                  dispatch(marketplaceFetchRepositoryDetail({ id }))
                    .catch(() => {})
                );
              }
              return Promise.when(getProductDetails);
            });
            return Promise.when(promises).then(() => results);
          }),
    },
  });
};

export const billingFetchProfileSubscriptions = ({ docker_id }) => dispatch => {
  const url = `${BILLING_BASE_URL}/subscriptions/`;
  return dispatch({
    type: BILLING_FETCH_PROFILE_SUBSCRIPTIONS,
    payload: {
      promise:
        request
          .get(url)
          .query({ docker_id })
          .set(bearer())
          .end()
          .then((res) => res.body),
    },
  });
};

export function billingCreateSubscription({
  docker_id,
  eusa,
  name,
  pricing_components, // Array of objects
  product_id,
  product_rate_plan,
}) {
  const params = {
    docker_id,
    eusa,
    name,
    pricing_components,
    product_id,
    product_rate_plan,
  };
  const url = `${BILLING_BASE_URL}/subscriptions/`;
  return {
    type: BILLING_CREATE_SUBSCRIPTION,
    payload: {
      promise:
        request
          .post(url)
          .query({ docker_id })
          .set(bearer())
          .send(params)
          .end()
          .then((res) => res.body),
    },
  };
}

// body is the request body
export const billingUpdateSubscription =
  ({ subscription_id, body }) => dispatch => {
    const url = `${BILLING_BASE_URL}/subscriptions/${subscription_id}/`;
    return dispatch({
      type: BILLING_UPDATE_SUBSCRIPTION,
      meta: { subscription_id },
      payload: {
        promise:
          request
            .patch(url)
            .accept('application/json')
            .set(bearer())
            .send(body)
            .end()
            .then((res) => res.body),
      },
    });
  };

export const billingDeleteSubscription = ({ subscription_id }) => dispatch => {
  const url = `${BILLING_BASE_URL}/subscriptions/${subscription_id}/`;
  return dispatch({
    type: BILLING_DELETE_SUBSCRIPTION,
    meta: { subscription_id },
    payload: {
      promise:
        request
          .del(url)
          .accept('application/json')
          .set(bearer())
          .send()
          .end()
          .then((res) => res.body),
    },
  });
};
