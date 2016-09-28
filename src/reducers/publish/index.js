import {
  REPOSITORY_FETCH_IMAGE_TAGS,
  REPOSITORY_FETCH_REPOSITORIES_FOR_NAMESPACE,
} from 'actions/repository';

import {
  PUBLISH_FETCH_PRODUCT_DETAILS,
  PUBLISH_FETCH_PRODUCT_LIST,
  PUBLISH_FETCH_PRODUCT_TIERS,
  PUBLISH_GET_PUBLISHERS,
  PUBLISH_GET_SIGNUP,
  PUBLISH_GET_VENDOR_AGREEMENT,
  PUBLISH_SUBSCRIBE,
} from 'actions/publish';

import cloneDeep from 'lodash/cloneDeep';
import get from 'lodash/get';

export const DEFAULT_STATE = {
  signup: {
    error: '',
    results: {},
  },
  productList: {
    error: '',
    results: [],
  },
  publishers: {
    error: '',
    results: {},
  },
  submit: {
    product_name: '',
    repositories: {},
    tags: {},
    agreement: {},
  },
  currentProductDetails: {
    error: '',
    results: {},
  },
  currentProductTiers: {
    error: '',
    results: [],
  },
};

export default function publish(state = DEFAULT_STATE, action) {
  const nextState = cloneDeep(state);
  const { payload, meta, type } = action;
  let metaNamespace;
  let metaReponame;
  if (meta) {
    metaNamespace = meta.namespace;
    metaReponame = meta.reponame;
  }

  switch (type) {
    //--------------------------------------------------------------------------
    // PUBLISH_SUBSCRIBE
    //--------------------------------------------------------------------------
    case `${PUBLISH_SUBSCRIBE}_ACK`:
      nextState.signup.results = payload || {};
      return nextState;
    case `${PUBLISH_SUBSCRIBE}_ERR`:
      nextState.signup.error = get(payload, ['response', 'error', 'message'],
        'Error subscribing user');
      return nextState;

    //--------------------------------------------------------------------------
    // PUBLISH_GET_SIGNUP
    //--------------------------------------------------------------------------
    case `${PUBLISH_GET_SIGNUP}_ACK`:
      nextState.signup.results = payload;
      nextState.signup.error = '';
      return nextState;

    case `${PUBLISH_GET_SIGNUP}_ERR`:
      nextState.signup.results = {};
      nextState.signup.error = 'User has not signed up to be a publisher';
      return nextState;

    //--------------------------------------------------------------------------
    // PUBLISH_GET_PUBLISHERS
    //--------------------------------------------------------------------------
    case `${PUBLISH_GET_PUBLISHERS}_ACK`:
      nextState.publishers.results = payload;
      nextState.publishers.error = '';
      return nextState;

    case `${PUBLISH_GET_PUBLISHERS}_ERR`:
      nextState.publishers.results = {};
      nextState.publishers.error = 'Publisher information unavailable';
      return nextState;

    //--------------------------------------------------------------------------
    // PUBLISH_FETCH_PRODUCT_LIST
    //--------------------------------------------------------------------------
    case `${PUBLISH_FETCH_PRODUCT_LIST}_ACK`:
      nextState.productList.results = payload;
      nextState.productList.error = '';
      return nextState;

    case `${PUBLISH_FETCH_PRODUCT_LIST}_ERR`:
      nextState.productList.results = [];
      nextState.productList.error = 'Product list unavailable';
      return nextState;

    //--------------------------------------------------------------------------
    // PUBLISH_FETCH_PRODUCT_DETAILS
    //--------------------------------------------------------------------------
    case `${PUBLISH_FETCH_PRODUCT_DETAILS}_ACK`:
      nextState.currentProductDetails.error = '';
      nextState.currentProductDetails.results = payload;
      return nextState;

    case `${PUBLISH_FETCH_PRODUCT_DETAILS}_ERR`:
      nextState.currentProductDetails.error =
        get(payload, ['response', 'error', 'message']);
      nextState.currentProductDetails.results = {};
      return nextState;

    //--------------------------------------------------------------------------
    // PUBLISH_GET_VENDOR_AGREEMENT
    //--------------------------------------------------------------------------
    case `${PUBLISH_GET_VENDOR_AGREEMENT}_ACK`:
      nextState.submit.agreement.results = payload.html;
      return nextState;

    case `${PUBLISH_GET_VENDOR_AGREEMENT}_ERR`:
      nextState.submit.agreement.error =
        get(payload, ['response', 'error', 'message']);
      return nextState;

    //--------------------------------------------------------------------------
    // REPOSITORY_FETCH_IMAGE_TAGS
    //--------------------------------------------------------------------------
    case `${REPOSITORY_FETCH_IMAGE_TAGS}_REQ`:
      if (!nextState.submit.tags[metaNamespace]) {
        nextState.submit.tags[metaNamespace] = {};
      }
      if (!nextState.submit.tags[metaNamespace][metaReponame]) {
        nextState.submit.tags[metaNamespace][metaReponame] = {};
      }
      nextState.submit.tags[metaNamespace][metaReponame].isFetching = true;
      return nextState;

    case `${REPOSITORY_FETCH_IMAGE_TAGS}_ACK`:
      if (!nextState.submit.tags[metaNamespace]) {
        nextState.submit.tags[metaNamespace] = {};
      }
      if (!nextState.submit.tags[metaNamespace][metaReponame]) {
        nextState.submit.tags[metaNamespace][metaReponame] = {};
      }
      nextState.submit.tags[metaNamespace][metaReponame].isFetching = false;
      nextState.submit.tags[metaNamespace][metaReponame].results =
        payload.results.map(t => t.name);
      return nextState;

    case `${REPOSITORY_FETCH_IMAGE_TAGS}_ERR`:
      if (!nextState.submit.tags[metaNamespace]) {
        nextState.submit.tags[metaNamespace] = {};
      }
      if (!nextState.submit.tags[metaNamespace][metaReponame]) {
        nextState.submit.tags[metaNamespace][metaReponame] = {};
      }
      nextState.submit.tags[metaNamespace][metaReponame].isFetching = false;
      nextState.submit.tags[metaNamespace][metaReponame].error =
        get(payload, ['response', 'error', 'message']);
      return nextState;

    //--------------------------------------------------------------------------
    // REPOSITORY_FETCH_REPOSITORIES_FOR_NAMESPACE
    //--------------------------------------------------------------------------
    case `${REPOSITORY_FETCH_REPOSITORIES_FOR_NAMESPACE}_REQ`:
      nextState.submit.repositories[metaNamespace] = { isFetching: true };
      return nextState;

    case `${REPOSITORY_FETCH_REPOSITORIES_FOR_NAMESPACE}_ACK`:
      if (!nextState.submit.repositories[metaNamespace]) {
        nextState.submit.repositories[metaNamespace] = {};
      }
      nextState.submit.repositories[metaNamespace].isFetching = false;
      nextState.submit.repositories[metaNamespace].results =
        payload.results.map(r => r.name);
      return nextState;

    case `${REPOSITORY_FETCH_REPOSITORIES_FOR_NAMESPACE}_ERR`:
      if (!nextState.submit.repositories[metaNamespace]) {
        nextState.submit.repositories[metaNamespace] = {};
      }
      nextState.submit.repositories[metaNamespace].isFetching = false;
      nextState.submit.repositories[metaNamespace].results = {};
      nextState.submit.repositories[metaNamespace].error =
        get(payload, ['response', 'error', 'message']);
      return nextState;

    //--------------------------------------------------------------------------
    // PUBLISH_FETCH_PRODUCT_TIERS
    //--------------------------------------------------------------------------
    case `${PUBLISH_FETCH_PRODUCT_TIERS}_ACK`:
      nextState.currentProductTiers.results = payload;
      nextState.currentProductTiers.error = '';
      return nextState;

    case `${PUBLISH_FETCH_PRODUCT_TIERS}_ERR`:
      nextState.currentProductTiers.results = [];
      nextState.currentProductTiers.error =
      get(payload, ['response', 'error', 'message'], 'Error subscribing user');
      return nextState;

    default:
      return state;
  }
}
