/* eslint-disable max-len */
import accountReducer,
  { DEFAULT_STATE } from 'reducers/account';
import {
  ACCOUNT_FETCH_CURRENT_USER_INFORMATION,
  ACCOUNT_FETCH_USER_EMAILS,
  ACCOUNT_FETCH_USER_INFORMATION,
  ACCOUNT_FETCH_USER_ORGS,
  ACCOUNT_LOGOUT,
} from 'actions/account';
import {
  REPOSITORY_FETCH_OWNED_NAMESPACES,
} from 'actions/repository';
import {
  WHITELIST_AM_I_WAITING,
  WHITELIST_FETCH_AUTHORIZATION,
} from 'actions/whitelist';
import { expect, assert } from 'chai';
import filter from 'lodash/filter';

describe('account reducer', () => {
  it('should return the initial state', () => {
    expect(accountReducer(undefined, {})).to.deep.equal(DEFAULT_STATE);
  });

  //----------------------------------------------------------------------------
  // ACCOUNT_FETCH_CURRENT_USER_INFORMATION
  //----------------------------------------------------------------------------
  it('should handle ACCOUNT_FETCH_CURRENT_USER_INFORMATION_ACK', () => {
    const payload = {
      id: '4663b07ca74111e492090242ac110143',
      username: 'test1',
      full_name: 'asdfasdfasdf',
      location: 'asdf',
      company: 'stuff',
      gravatar_email: '',
      is_staff: false,
      is_admin: false,
      profile_url: '',
      date_joined: '2014-09-23T19:42:13Z',
      gravatar_url: 'https://secure.gravatar.com/avatar/88d62a9d7579193eea16d4f5ddee3f62.jpg?s=80&r=g&d=mm',
      type: 'User',
    };
    const action = {
      type: `${ACCOUNT_FETCH_CURRENT_USER_INFORMATION}_ACK`,
      payload,
    };
    const reducer = accountReducer(undefined, action);
    expect(reducer.currentUser).to.equal(payload);
    expect(reducer.namespaceObjects).to.have.property('isFetching');
    expect(reducer.namespaceObjects).to.have.property('error');
    expect(reducer.namespaceObjects).to.have.property('results');
    const results = reducer.namespaceObjects.results;
    expect(results).to.have.property(payload.username);
    expect(results[payload.username]).to.exist;
    expect(results[payload.username]).to.equal(payload);
  });

  it('should handle ACCOUNT_FETCH_USER_INFORMATION', () => {
    const payload = {
      id: '4663b07ca74111e492090242ac110143',
      orgname: 'test1',
      full_name: 'asdfasdfasdf',
      location: 'asdf',
      company: 'stuff',
      gravatar_email: '',
      is_staff: false,
      is_admin: false,
      profile_url: '',
      date_joined: '2014-09-23T19:42:13Z',
      gravatar_url: 'https://secure.gravatar.com/avatar/88d62a9d7579193eea16d4f5ddee3f62.jpg?s=80&r=g&d=mm',
      type: 'User',
    };
    const action = {
      type: `${ACCOUNT_FETCH_USER_INFORMATION}_ACK`,
      payload,
    };
    const reducer = accountReducer(undefined, action);
    expect(reducer.namespaceObjects).to.have.property('isFetching');
    expect(reducer.namespaceObjects).to.have.property('error');
    expect(reducer.namespaceObjects).to.have.property('results');
    const results = reducer.namespaceObjects.results;
    expect(results).to.have.property(payload.orgname);
    expect(results[payload.orgname]).to.exist;
    expect(results[payload.orgname]).to.equal(payload);
  });

  //----------------------------------------------------------------------------
  // ACCOUNT_FETCH_USER_EMAILS
  //----------------------------------------------------------------------------
  it('should handle ACCOUNT_FETCH_USER_EMAILS', () => {
    const payload = {
      count: 2,
      next: null,
      previous: null,
      results: [
        {
          id: 34517,
          user: 'test1',
          email: 'nathan.hsieh+123@docker.com',
          verified: false,
          primary: true,
        },
        {
          id: 22488,
          user: 'test1',
          email: 'hsieh.nathan+111@gmail.com',
          verified: true,
          primary: false,
        },
      ],
    };
    const action = {
      type: `${ACCOUNT_FETCH_USER_EMAILS}_ACK`,
      payload,
    };
    const reducer = accountReducer(undefined, action);
    expect(reducer.userEmails).to.equal(payload);
    expect(reducer.userEmails).to.have.property('results').with.length(payload.count);
    const primary = filter(payload.results, email => email.primary);
    expect(primary).to.have.length(1);
    assert.isObject(primary[0]);
  });

  //----------------------------------------------------------------------------
  // REPOSITORY_FETCH_OWNED_NAMESPACES
  //----------------------------------------------------------------------------
  it('should handle REPOSITORY_FETCH_OWNED_NAMESPACES', () => {
    const payload = {
      namespaces: ['test1', 'test2', 'test3'],
    };
    const action = {
      type: `${REPOSITORY_FETCH_OWNED_NAMESPACES}_ACK`,
      payload,
    };
    const reducer = accountReducer(undefined, action);
    expect(reducer.ownedNamespaces).to.equal(payload.namespaces);
  });

  //----------------------------------------------------------------------------
  // ACCOUNT_FETCH_USER_ORGS
  //----------------------------------------------------------------------------
  it('should handle ACCOUNT_FETCH_USER_ORGS_REQ', () => {
    const action = {
      type: `${ACCOUNT_FETCH_USER_ORGS}_REQ`,
    };
    const reducer = accountReducer(undefined, action);
    const { isFetching, results, error } = reducer.namespaceObjects;
    expect(isFetching).to.equal(true);
    expect(error).to.equal('');
    expect(Object.keys(results).length).to.equal(0);
  });

  it('should handle ACCOUNT_FETCH_USER_ORGS_ACK', () => {
    const org = { orgname: 'haha' };
    const action = {
      type: `${ACCOUNT_FETCH_USER_ORGS}_ACK`,
      payload: { results: [org] },
    };
    const reducer = accountReducer(undefined, action);
    const { isFetching, results, error } = reducer.namespaceObjects;
    expect(isFetching).to.equal(false);
    expect(error).to.equal('');
    expect(results).to.have.property(org.orgname);
    expect(results[org.orgname]).to.deep.equal(org);
  });

  it('should handle ACCOUNT_FETCH_USER_ORGS_ERR', () => {
    const action = {
      type: `${ACCOUNT_FETCH_USER_ORGS}_ERR`,
    };
    const reducer = accountReducer(undefined, action);
    const { isFetching, error } = reducer.namespaceObjects;
    expect(isFetching).to.equal(false);
    expect(error).to.exist;
  });

  //----------------------------------------------------------------------------
  // ACCOUNT_SELECT_NAMESPACE
  //----------------------------------------------------------------------------
  it('should handle ACCOUNT_SELECT_NAMESPACE', () => {
    const action = {
      type: 'ACCOUNT_SELECT_NAMESPACE',
      payload: { namespace: 'testing' },
    };
    const reducer = accountReducer(undefined, action);
    expect(reducer.selectedNamespace).to.exist;
    expect(reducer.selectedNamespace).to.equal(action.payload.namespace);
  });

  //----------------------------------------------------------------------------
  // ACCOUNT_TOGGLE_MAGIC_CARPET
  //----------------------------------------------------------------------------
  it('should handle ACCOUNT_TOGGLE_MAGIC_CARPET', () => {
    const action = {
      type: 'ACCOUNT_TOGGLE_MAGIC_CARPET',
      payload: { magicCarpet: 'billing' },
    };
    const reducer = accountReducer(undefined, action);
    expect(reducer.magicCarpet).to.exist;
    expect(reducer.magicCarpet).to.equal(action.payload.magicCarpet);
  });

  //----------------------------------------------------------------------------
  // WHITELIST_FETCH_AUTHORIZATION
  //----------------------------------------------------------------------------
  it('should handle WHITELIST_FETCH_AUTHORIZATION_ACK', () => {
    const action = {
      type: `${WHITELIST_FETCH_AUTHORIZATION}_ACK`,
    };
    const reducer = accountReducer(undefined, action);
    expect(reducer.isCurrentUserWhitelisted).to.equal(true);
  });

  it('should handle WHITELIST_FETCH_AUTHORIZATION_ERR', () => {
    const action = {
      type: `${WHITELIST_FETCH_AUTHORIZATION}_ERR`,
    };
    const reducer = accountReducer(undefined, action);
    expect(reducer.isCurrentUserWhitelisted).to.equal(false);
  });

  //----------------------------------------------------------------------------
  // WHITELIST_AM_I_WAITING
  //----------------------------------------------------------------------------
  it('should handle WHITELIST_AM_I_WAITING_ACK', () => {
    const action = {
      type: `${WHITELIST_AM_I_WAITING}_ACK`,
    };
    const reducer = accountReducer(undefined, action);
    expect(reducer.isCurrentUserBetalisted).to.equal(true);
  });

  it('should handle WHITELIST_AM_I_WAITING_ERR', () => {
    const action = {
      type: `${WHITELIST_AM_I_WAITING}_ERR`,
    };
    const reducer = accountReducer(undefined, action);
    expect(reducer.isCurrentUserBetalisted).to.equal(false);
  });

  //----------------------------------------------------------------------------
  // ACCOUNT_LOGOUT_ACK
  //----------------------------------------------------------------------------
  it('should handle ACCOUNT_LOGOUT_ACK', () => {
    const action = {
      type: `${ACCOUNT_LOGOUT}_ACK`,
    };
    // simulate a logged in user
    const REQ_STATE = { currentUser: { id: '1234567' } };
    const reducer = accountReducer(REQ_STATE, action);
    expect(reducer).to.deep.equal(DEFAULT_STATE);
  });
});
