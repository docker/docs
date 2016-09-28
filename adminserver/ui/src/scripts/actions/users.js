'use strict';

import { Accounts, Index } from 'dtr-js-sdk';
import consts from 'consts';

const AccountsAPI = Accounts.default;

/**
 * Fetch an organization's information by its name
 */
export function getUser(name) {
  return {
    type: consts.users.GET_USER,
    meta: {
      promiseScope: [name],
      promise: AccountsAPI.getAccount({ name })
      .then(response => {
        return response.data;
      })
    }
  };
}

export function changeUserPassword({ name, oldPassword, newPassword }) {
  return {
    type: consts.users.CHANGE_USER_PASSWORD,
    meta: {
      promise: AccountsAPI.changeUserPassword({ name }, { oldPassword, newPassword }),
      notifications: {
        pending: 'Changing password',
        success: 'New password set',
        error: 'There was an error setting the new password.'
      }
    }
  };
}

export function listUsers(start = 0) {
  return {
    type: consts.users.LIST_USERS,
    meta: {
      promiseScope: [start],
      promise: AccountsAPI.listAccounts({}, { filter: 'users', limit: '200', start })
        .then(response => response.data.accounts)
    }
  };
}

export function searchUsers(searchTerm) {
  return {
    type: consts.users.SEARCH_USERS,
    searchTerm,
    meta: {
      promise: Index.autocomplete({
          query: searchTerm,
          includeRepositories: false,
          includeAccounts: true
        })
        .then(
          resp => {
            return (resp.data.accountResults.filter((account) => {
              return !account.isOrg;
            }) || []);
          }
        )
    }
  };
}

export function createUser(data) {
    return {
        type: consts.users.CREATE_USER,
        meta: {
            promise: AccountsAPI.createAccount({}, data),
            notifications: {
                pending: 'Creating user',
                success: 'Created new user',
                error: (resp) => {
                  return `Unable to create the user:\n ${resp.data.errors[0].message}`;
                }
            }
        }
    };
}

export function updateUser(name, data) {
  return {
    type: consts.users.UPDATE_ACCOUNT,
    meta: {
      promise: AccountsAPI.updateAccount({ name }, data),
      notifications: {
        pending: 'Updating user',
        success: 'Updated user',
        error: 'There was an error updating the user'
      }
    }
  };
}

export function deleteUser(name) {
  return {
    type: consts.users.DELETE_USER,
    meta: {
      promise: AccountsAPI.deleteAccount({ name }).then(response => {
        if (response.status === 204) {
          window.location.href = '/users';
        } else {
          throw new Error('Didn\'t delete the user.');
        }
      }),
      notifications: {
        pending: 'Deleting user',
        success: 'Deleted user',
        error: 'There was an error deleting the user'
      }
    }
  };
}
