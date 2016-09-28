import {
  BILLING_CREATE_PAYMENT_METHOD,
  BILLING_CREATE_PAYMENT_TOKEN,
  BILLING_CREATE_PROFILE,
  BILLING_CREATE_SUBSCRIPTION,
  BILLING_DELETE_SUBSCRIPTION,
  BILLING_FETCH_INVOICES,
  BILLING_FETCH_LICENSE_DETAIL,
  BILLING_FETCH_PAYMENT_METHODS,
  BILLING_FETCH_PRODUCT,
  BILLING_FETCH_PROFILE_SUBSCRIPTIONS,
  BILLING_FETCH_PROFILE_SUBSCRIPTIONS_AND_PRODUCTS,
  BILLING_FETCH_PROFILE,
  BILLING_UPDATE_PROFILE,
  BILLING_UPDATE_SUBSCRIPTION,
} from 'actions/billing';
import cloneDeep from 'lodash/cloneDeep';
import get from 'lodash/get';
import merge from 'lodash/merge';

export const DEFAULT_STATE = {
  invoices: {
    isFetching: false,
    results: [],
    error: '',
  },
  paymentMethods: {
    isFetching: false,
    isSubmitting: false,
    results: [],
    error: {},
    fetchingError: '',
  },
  products: {},
  profiles: {
    isSubmitting: false,
    isFetching: false,
    results: {},
  },
  subscriptions: {
    // metadata about deleting a subscription
    delete: {},
    error: '',
    isFetching: false,
    isSubmitting: false,
    licenses: {
      error: '',
      isFetching: false,
      // licenses are keyed by subscription_id
      results: {},
    },
    results: {},
    // metadata about updating a subscription, keyed on sub_id
    // sub_id: { error, isUpdating }
    update: {},
  },
};

const reduceSubscriptions = (payload) => {
  const results = {};
  payload.forEach((subscription) => {
    const { subscription_id } = subscription;
    results[subscription_id] = subscription;
  });
  return results;
};

// https://visionmedia.github.io/superagent/#error-handling
const reduceError = (payload) => {
  // Status will be a code like 404, etc.
  const status = get(payload, ['response', 'error', 'status']);
  // Message may come from the API
  const message = get(payload, ['response', 'error', 'message']) ||
    get(payload, ['response', 'body', 'message']);
  return message || status;
};

export default function billing(state = DEFAULT_STATE, action) {
  const nextState = cloneDeep(state);
  const { payload, type, meta } = action;
  /* eslint-disable max-len*/
  let errorObject = {
    message: 'There was an issue adding your payment information. Please try again or contact support.',
  };
  if (type === `${BILLING_CREATE_PAYMENT_TOKEN}_ERR`) {
    const error = payload.response.body.error;
    errorObject = error.type === 'card_error' ? error : errorObject;
  }

  switch (type) {
    //--------------------------------------------------------------------------
    // BILLING_CREATE_PAYMENT_METHOD
    //--------------------------------------------------------------------------
    case `${BILLING_CREATE_PAYMENT_TOKEN}_REQ`:
      nextState.paymentMethods.isSubmitting = true;
      nextState.paymentMethods.error = {};
      return nextState;
    case `${BILLING_CREATE_PAYMENT_TOKEN}_ACK`:
      nextState.paymentMethods.error = {};
      return nextState;
    case `${BILLING_CREATE_PAYMENT_TOKEN}_ERR`:
      nextState.paymentMethods.isSubmitting = false;
      nextState.paymentMethods.error = errorObject;
      return nextState;
    case `${BILLING_CREATE_PAYMENT_METHOD}_ACK`:
      nextState.paymentMethods.isSubmitting = false;
      return nextState;
    case `${BILLING_CREATE_PAYMENT_METHOD}_ERR`:
      nextState.paymentMethods.isSubmitting = false;
      nextState.paymentMethods.error = errorObject;
      return nextState;

    //--------------------------------------------------------------------------
    // BILLING_CREATE_SUBSCRIPTION
    //--------------------------------------------------------------------------
    case `${BILLING_CREATE_SUBSCRIPTION}_REQ`:
      nextState.subscriptions.isFetching = true;
      nextState.subscriptions.isSubmitting = true;
      nextState.subscriptions.error = '';
      return nextState;
    case `${BILLING_CREATE_SUBSCRIPTION}_ACK`:
      nextState.subscriptions.isFetching = false;
      nextState.subscriptions.isSubmitting = false;
      nextState.subscriptions.results = {
        [payload.subscription_id]: payload,
      };
      nextState.subscriptions.error = '';
      return nextState;
    case `${BILLING_CREATE_SUBSCRIPTION}_ERR`:
      nextState.subscriptions.isFetching = false;
      nextState.subscriptions.isSubmitting = false;
      nextState.subscriptions.error = 'There was an issue creating your subscription. Please try again or contact support.';
      return nextState;

    //--------------------------------------------------------------------------
    // BILLING_DELETE_SUBSCRIPTION
    //--------------------------------------------------------------------------
    case `${BILLING_DELETE_SUBSCRIPTION}_REQ`:
      nextState.subscriptions.delete.isDeleting = true;
      nextState.subscriptions.delete.subscription_id = meta.subscription_id;
      nextState.subscriptions.delete.error = '';
      return nextState;
    case `${BILLING_DELETE_SUBSCRIPTION}_ACK`:
      nextState.subscriptions.delete.isDeleting = false;
      nextState.subscriptions.delete.subscription_id = '';
      nextState.subscriptions.delete.error = '';
      return nextState;
    case `${BILLING_DELETE_SUBSCRIPTION}_ERR`:
      nextState.subscriptions.delete.isDeleting = false;
      nextState.subscriptions.delete.error =
        'Sorry, we are unable to delete this subscription';
      return nextState;

    //--------------------------------------------------------------------------
    // BILLING_UPDATE_SUBSCRIPTION
    //--------------------------------------------------------------------------
    case `${BILLING_UPDATE_SUBSCRIPTION}_REQ`:
      nextState.subscriptions.update[meta.subscription_id] = {
        isUpdating: true,
        error: '',
      };
      return nextState;
    case `${BILLING_UPDATE_SUBSCRIPTION}_ACK`:
      nextState.subscriptions.update[meta.subscription_id] = {
        isUpdating: false,
        error: '',
      };
      return nextState;
    case `${BILLING_UPDATE_SUBSCRIPTION}_ERR`:
      nextState.subscriptions.update[meta.subscription_id] = {
        isUpdating: false,
        error: 'Sorry, we are unable to update this subscription',
      };
      return nextState;

    //--------------------------------------------------------------------------
    // BILLING_FETCH_INVOICES
    //--------------------------------------------------------------------------
    case `${BILLING_FETCH_INVOICES}_REQ`:
      nextState.invoices.isFetching = true;
      nextState.invoices.error = '';
      return nextState;
    case `${BILLING_FETCH_INVOICES}_ACK`:
      nextState.invoices.isFetching = false;
      nextState.invoices.results = payload;
      nextState.invoices.error = '';
      return nextState;
    case `${BILLING_FETCH_INVOICES}_ERR`:
      nextState.invoices.isFetching = false;
      nextState.invoices.error = 'There was an issue when fetching your Invoices';
      return nextState;

    //--------------------------------------------------------------------------
    // BILLING_FETCH_PAYMENT_METHODS
    //--------------------------------------------------------------------------
    case `${BILLING_FETCH_PAYMENT_METHODS}_REQ`:
      nextState.paymentMethods.isFetching = true;
      nextState.paymentMethods.fetchingError = '';
      return nextState;
    case `${BILLING_FETCH_PAYMENT_METHODS}_ACK`:
      nextState.paymentMethods.isFetching = false;
      nextState.paymentMethods.results = payload;
      nextState.paymentMethods.fetchingError = '';
      return nextState;
    case `${BILLING_FETCH_PAYMENT_METHODS}_ERR`:
      nextState.paymentMethods.isFetching = false;
      nextState.paymentMethods.fetchingError = 'There was an issue when fetching your Payment Methods';
      return nextState;

    //--------------------------------------------------------------------------
    // BILLING_FETCH_PROFILE
    //--------------------------------------------------------------------------
    case `${BILLING_FETCH_PROFILE}_REQ`:
      nextState.profiles.isFetching = true;
      nextState.paymentMethods.isFetching = false;
      nextState.paymentMethods.results = [];
      nextState.paymentMethods.error = {};
      return nextState;
    case `${BILLING_FETCH_PROFILE}_ACK`:
      nextState.profiles.isFetching = false;
      nextState.profiles.error = '';
      nextState.profiles.results[meta.docker_id] = payload.profile;
      return nextState;
    case `${BILLING_FETCH_PROFILE}_ERR`:
      nextState.profiles.isFetching = false;
      nextState.profiles.error = '';
      return nextState;

    //--------------------------------------------------------------------------
    // BILLING_CREATE_PROFILE
    //--------------------------------------------------------------------------
    case `${BILLING_CREATE_PROFILE}_REQ`:
      nextState.profiles.isFetching = true;
      nextState.profiles.isSubmitting = true;
      return nextState;
    case `${BILLING_CREATE_PROFILE}_ACK`:
      nextState.profiles.isFetching = false;
      nextState.profiles.isSubmitting = false;
      nextState.profiles.error = '';
      nextState.profiles.results[meta.docker_id] = payload.profile;
      return nextState;
    case `${BILLING_CREATE_PROFILE}_ERR`:
      nextState.profiles.isFetching = false;
      nextState.profiles.isSubmitting = false;
      nextState.profiles.error = 'There was an issue creating your billing profile. Please try again or contact support.';
      return nextState;

    //--------------------------------------------------------------------------
    // BILLING_UPDATE_PROFILE
    //--------------------------------------------------------------------------
    case `${BILLING_UPDATE_PROFILE}_REQ`:
      nextState.profiles.isFetching = true;
      nextState.profiles.isSubmitting = true;
      return nextState;
    case `${BILLING_UPDATE_PROFILE}_ACK`:
      nextState.profiles.isFetching = false;
      nextState.profiles.isSubmitting = false;
      nextState.profiles.error = '';
      merge(nextState.profiles.results[meta.docker_id], payload);
      return nextState;
    case `${BILLING_UPDATE_PROFILE}_ERR`:
      nextState.profiles.isFetching = false;
      nextState.profiles.isSubmitting = false;
      nextState.profiles.error = 'There was an issue updating your billing profile. Please try again or contact support.';
      return nextState;

    //--------------------------------------------------------------------------
    // BILLING_FETCH_PRODUCT
    //--------------------------------------------------------------------------
    case `${BILLING_FETCH_PRODUCT}_REQ`:
      nextState.products[meta.id] = { isFetching: true };
      return nextState;
    case `${BILLING_FETCH_PRODUCT}_ACK`:
      nextState.products[meta.id] = { ...payload, isFetching: false };
      return nextState;
    case `${BILLING_FETCH_PRODUCT}_ERR`:
      nextState.products[meta.id].isFetching = false;
      nextState.products[meta.id].error = reduceError(payload);
      return nextState;

    //--------------------------------------------------------------------------
    // BILLING_FETCH_PROFILE_SUBSCRIPTIONS
    //--------------------------------------------------------------------------
    case `${BILLING_FETCH_PROFILE_SUBSCRIPTIONS_AND_PRODUCTS}_REQ`:
    case `${BILLING_FETCH_PROFILE_SUBSCRIPTIONS}_REQ`:
      nextState.subscriptions.isFetching = true;
      nextState.subscriptions.error = '';
      return nextState;
    case `${BILLING_FETCH_PROFILE_SUBSCRIPTIONS_AND_PRODUCTS}_ACK`:
    case `${BILLING_FETCH_PROFILE_SUBSCRIPTIONS}_ACK`:
      // We must create a new results object each time so that we don't
      // accidently persist old subscriptions from another namespace
      nextState.subscriptions.isFetching = false;
      nextState.subscriptions.results = reduceSubscriptions(payload);
      nextState.subscriptions.error = '';
      return nextState;
    case `${BILLING_FETCH_PROFILE_SUBSCRIPTIONS_AND_PRODUCTS}_ERR`:
    case `${BILLING_FETCH_PROFILE_SUBSCRIPTIONS}_ERR`:
      nextState.subscriptions.isFetching = false;
      nextState.subscriptions.results = {};
      nextState.subscriptions.error = 'Unable to fetch subscriptions';
      return nextState;

    //--------------------------------------------------------------------------
    // BILLING_FETCH_LICENSE_DETAIL
    //--------------------------------------------------------------------------
    case `${BILLING_FETCH_LICENSE_DETAIL}_REQ`:
      nextState.subscriptions.licenses.isFetching = true;
      nextState.subscriptions.licenses.error = '';
      return nextState;
    case `${BILLING_FETCH_LICENSE_DETAIL}_ACK`:
      nextState.subscriptions.licenses.isFetching = false;
      nextState.subscriptions.licenses.results[meta.subscription_id] = payload;
      return nextState;
    case `${BILLING_FETCH_LICENSE_DETAIL}_ERR`:
      nextState.subscriptions.licenses.isFetching = false;
      nextState.subscriptions.licenses.error = reduceError(payload);
      return nextState;

    default:
      return state;
  }
}
