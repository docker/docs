const request = require('superagent-promise')(require('superagent'), Promise);
import { createAction } from 'redux-actions';
import { jwt } from 'lib/utils/authHeaders';
import { post } from 'superagent';
import { readCookie } from 'lib/utils/cookie-handler';

const DOCKERCLOUD_HOST = '';
const ACCOUNT_BASE_URL = `${DOCKERCLOUD_HOST}/v2`;

// eslint-disable-next-line
export const ACCOUNT_FETCH_CURRENT_USER_INFORMATION = 'ACCOUNT_FETCH_CURRENT_USER_INFORMATION';
export const ACCOUNT_FETCH_USER_EMAILS = 'ACCOUNT_FETCH_USER_EMAILS';
export const ACCOUNT_FETCH_USER_INFORMATION = 'ACCOUNT_FETCH_USER_INFORMATION';
export const ACCOUNT_FETCH_USER_NAMESPACES = 'ACCOUNT_FETCH_USER_NAMESPACES';
export const ACCOUNT_FETCH_USER_ORGS = 'ACCOUNT_FETCH_USER_ORGS';
export const ACCOUNT_LOGOUT = 'ACCOUNT_LOGOUT';
export const ACCOUNT_TOGGLE_MAGIC_CARPET = 'ACCOUNT_TOGGLE_MAGIC_CARPET';
export const ACCOUNT_SELECT_NAMESPACE = 'ACCOUNT_SELECT_NAMESPACE';

export const login = (username, password) => {
  return new Promise((resolve, reject) => {
    const endpoint = '/v2/users/login/';
    const req = post(endpoint)
      .set('Content-Type', 'application/json')
      .set('Accept', 'application/json')
      .set('X-CSRFToken', readCookie('csrftoken'));

    req.send({ username, password }).end((err, res) => {
      if (err) {
        const errors = {};
        let body = {};

        try {
          body = JSON.parse(res.text);
        } catch (e) {
          errors._error = res.text;
          reject(errors);
        }

        if (body.detail) {
          errors._error = body.detail;
        }

        if (body.username) {
          errors.username = body.username[0];
        }

        if (body.password) {
          errors.password = body.password[0];
        }

        reject(errors);
        return;
      }

      resolve();
    });
  });
};

export const signup = (email, username, password, redirect_value) => {
  return new Promise((resolve, reject) => {
    const endpoint = '/v2/users/signup/';
    const req = post(endpoint)
      .set('Content-Type', 'application/json')
      .set('Accept', 'application/json')
      .set('X-CSRFToken', readCookie('csrftoken'));

    req.send({ email, username, password, redirect_value }).end((err, res) => {
      if (err) {
        const errors = {};
        let body = {};

        try {
          body = JSON.parse(res.text);
        } catch (e) {
          errors._error = res.text;
          reject(errors);
        }

        if (body.detail) {
          errors._error = body.detail;
        }

        if (body.email) {
          errors.email = body.email[0];
        }

        if (body.username) {
          errors.username = body.username[0];
        }

        if (body.password) {
          errors.password = body.password[0];
        }

        reject(errors);
        return;
      }

      resolve();
    });
  });
};


export function accountLogout() {
  const accountUrl = `${ACCOUNT_BASE_URL}/user/logout`;
  return {
    type: ACCOUNT_LOGOUT,
    payload: {
      promise:
        request
          .post(accountUrl)
          .set(jwt())
          .set('Accept', '*/*')
          .end()
          .then((res) => res.body),
    },
  };
}

/*
Sample JSON for user fetch
  {
    "id": "4663b07ca74111e492090242ac110143",
    "username": "test1",
    "full_name": "asdfasdfasdf",
    "location": "asdf",
    "company": "stuff",
    "gravatar_email": "",
    "is_staff": false,
    "is_admin": false,
    "profile_url": "",
    "date_joined": "2014-09-23T19:42:13Z",
    "gravatar_url": "https://secure.gravatar.com/avatar/88d62a9d7579193eea16d4f5ddee3f62.jpg?s=80&r=g&d=mm",
    "type": "User"
  }
*/
export function accountFetchUser({ namespace, isOrg }) {
  const userOrOrg = isOrg ? 'orgs' : 'users';
  const accountUrl = `${ACCOUNT_BASE_URL}/${userOrOrg}/${namespace}/`;
  return {
    type: ACCOUNT_FETCH_USER_INFORMATION,
    payload: {
      promise: request
                .get(accountUrl)
                .set(jwt())
                .end()
                .then((res) => res.body),
    },
  };
}

export function accountFetchCurrentUser({ shouldRedirectToLogin } = {}) {
  const url = `${ACCOUNT_BASE_URL}/user/`;
  return {
    type: ACCOUNT_FETCH_CURRENT_USER_INFORMATION,
    meta: {
      shouldRedirectToLogin,
    },
    payload: {
      promise:
        request
          .get(url)
          .set(jwt())
          .end()
          .then((res) => res.body),
    },
  };
}

export function accountFetchUserEmails({ user }) {
  const url = `${ACCOUNT_BASE_URL}/emailaddresses/`;
  return {
    type: ACCOUNT_FETCH_USER_EMAILS,
    payload: {
      promise:
        request
          .get(url)
          .set(jwt())
          .query({ user })
          .end()
          .then((res) => res.body),
    },
  };
}

export function accountFetchUserOrgs() {
  // Fetches user objects for all namespaces you have access to - not just own
  const url = `${ACCOUNT_BASE_URL}/user/orgs/`;
  return {
    type: ACCOUNT_FETCH_USER_ORGS,
    payload: {
      promise:
        request
          .get(url)
          .accept('application/json')
          .set(jwt())
          .query({ page_size: 100 })
          .end()
          .then((res) => res.body),
    },
  };
}

export const accountSelectNamespace = createAction(ACCOUNT_SELECT_NAMESPACE);
export const accountToggleMagicCarpet =
  createAction(ACCOUNT_TOGGLE_MAGIC_CARPET);
