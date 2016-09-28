import {
  ACCOUNT_FETCH_CURRENT_USER_INFORMATION,
  ACCOUNT_FETCH_USER_EMAILS,
  ACCOUNT_FETCH_USER_INFORMATION,
  ACCOUNT_FETCH_USER_ORGS,
  ACCOUNT_TOGGLE_MAGIC_CARPET,
  ACCOUNT_LOGOUT,
  ACCOUNT_SELECT_NAMESPACE,
} from 'actions/account';
import {
  WHITELIST_FETCH_AUTHORIZATION,
  WHITELIST_AM_I_WAITING,
} from 'actions/whitelist';
import {
  REPOSITORY_FETCH_OWNED_NAMESPACES,
} from 'actions/repository';
import { cloneDeep } from 'lodash';

export const DEFAULT_STATE = {
  // Current User is the profile of the logged in user
  currentUser: {},
  isCurrentUserWhitelisted: false,
  isCurrentUserBetalisted: false,
  userEmails: {},
  // Array of strings (namespaces that the user is an owner of)
  ownedNamespaces: [],
  namespaceObjects: {
    isFetching: false,
    // Object of user (org) objects
    results: {},
    error: '',
  },
  selectedNamespace: '',
  magicCarpet: '',
};

export default function account(state = DEFAULT_STATE, action) {
  const nextState = cloneDeep(state);
  const { payload, type } = action;
  let accountName;
  if (type === `${ACCOUNT_FETCH_USER_INFORMATION}_ACK`) {
    accountName = payload.username || payload.orgname;
  }

  switch (type) {
    //--------------------------------------------------------------------------
    // ACCOUNT_FETCH_CURRENT_USER_INFORMATION
    //--------------------------------------------------------------------------
    case `${ACCOUNT_FETCH_CURRENT_USER_INFORMATION}_ACK`:
      nextState.currentUser = payload;
      if (!nextState.namespaceObjects.results[payload.username]) {
        nextState.namespaceObjects.results[payload.username] = payload;
      }
      return nextState;
    case `${ACCOUNT_FETCH_USER_INFORMATION}_ACK`:
      nextState.namespaceObjects.results[accountName] = payload;
      return nextState;

    //--------------------------------------------------------------------------
    // ACCOUNT_FETCH_USER_EMAILS
    //--------------------------------------------------------------------------
    case `${ACCOUNT_FETCH_USER_EMAILS}_ACK`:
      nextState.userEmails = payload;
      return nextState;

    //--------------------------------------------------------------------------
    // ACCOUNT_TOGGLE_MAGIC_CARPET
    //--------------------------------------------------------------------------
    case ACCOUNT_TOGGLE_MAGIC_CARPET:
      nextState.magicCarpet = payload.magicCarpet;
      return nextState;

    //--------------------------------------------------------------------------
    // ACCOUNT_SELECT_NAMESPACE
    //--------------------------------------------------------------------------
    case ACCOUNT_SELECT_NAMESPACE:
      nextState.selectedNamespace = payload.namespace;
      return nextState;

    //--------------------------------------------------------------------------
    // REPOSITORY_FETCH_OWNED_NAMESPACES
    //--------------------------------------------------------------------------
    case `${REPOSITORY_FETCH_OWNED_NAMESPACES}_ACK`:
      nextState.ownedNamespaces = payload.namespaces;
      return nextState;

    //--------------------------------------------------------------------------
    // ACCOUNT_FETCH_USER_ORGS
    //--------------------------------------------------------------------------
    case `${ACCOUNT_FETCH_USER_ORGS}_REQ`:
      nextState.namespaceObjects.isFetching = true;
      return nextState;
    case `${ACCOUNT_FETCH_USER_ORGS}_ACK`:
      nextState.namespaceObjects.isFetching = false;
      payload.results.forEach((org) => {
        const namespace = org.username || org.orgname;
        nextState.namespaceObjects.results[namespace] = org;
      });
      return nextState;
    case `${ACCOUNT_FETCH_USER_ORGS}_ERR`:
      nextState.namespaceObjects.isFetching = false;
      nextState.namespaceObjects.error = 'Unable to fetch user orgs';
      return nextState;

    //--------------------------------------------------------------------------
    // WHITELIST_FETCH_AUTHORIZATION
    //--------------------------------------------------------------------------
    case `${WHITELIST_FETCH_AUTHORIZATION}_ACK`:
      nextState.isCurrentUserWhitelisted = true;
      return nextState;
    case `${WHITELIST_FETCH_AUTHORIZATION}_ERR`:
      nextState.isCurrentUserWhitelisted = false;
      return nextState;

    //--------------------------------------------------------------------------
    // WHITELIST_AM_I_WAITING
    //--------------------------------------------------------------------------
    case `${WHITELIST_AM_I_WAITING}_ACK`:
      nextState.isCurrentUserBetalisted = true;
      return nextState;
    case `${WHITELIST_AM_I_WAITING}_ERR`:
      nextState.isCurrentUserBetalisted = false;
      return nextState;

    //--------------------------------------------------------------------------
    // ACCOUNT_LOGOUT
    //--------------------------------------------------------------------------
    case `${ACCOUNT_LOGOUT}_ACK`:
      // Clear current user so the store is refreshed
      return cloneDeep(DEFAULT_STATE);

    default:
      return state;
  }
}
