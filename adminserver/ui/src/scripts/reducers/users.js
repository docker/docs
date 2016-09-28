'use strict';

import consts from 'consts';
import { Map, fromJS } from 'immutable';

let defaultState = new Map({
  // Map from userName -> UserRecord
  users: new Map(),
  /*
  {
    userName: userData
  }
  */

  // `users` will contain every user less than `listMaxId`
  listMaxId: 0,

  // This is a static entry and will never change
  // TODO: convert to a record
  currentUser: window.user,
  isEmergencyAccess: window.isEmergencyAccess || false
});

const mapUsers = (users) => {
    const userdict = users.reduce((accum, user) => {
        accum[user.name] = user;
        return accum;
    }, {});
    return fromJS(userdict);
};

const actions = {
    [consts.users.LIST_USERS]: (state, action) => {
        if (!action.ready || action.error) {
            return state;
        }

        const users = action.payload;
        const listMaxId = Math.max(state.get('listMaxId'),
                                   Math.max.apply(null, users.map((user) => user.id)));

        state = state.set('listMaxId', listMaxId);
        state = state.set('users', mapUsers(users, state));
        return state;
    },
    [consts.users.SEARCH_USERS]: (state, action) => {
        if (!action.ready || action.error) {
            return state;
        }

        const users = action.payload;
        return state.set('users', mapUsers(users, state));
    },
    [consts.users.CREATE_USER]: (state, action) => {
        if (!action.ready || action.error) {
            return state;
        }
        const {
            payload: {
                data
            }
        } = action;
        return state.setIn(['users', data.name], fromJS(data));
    },
    [consts.users.UPDATE_ACCOUNT]: (state, action) => {
      if (!action.ready || action.error) {
        return state;
      }
      const {
        payload: {
          data
          }
        } = action;
      return state.setIn(['users', data.name], data);
    },
    [consts.users.DELETE_USER]: (state, action) => {
      if (!action.ready || action.error) {
        return state;
      }
      const {
        payload: {
          data
        }
      } = action;
      return state.deleteIn(['users', data.name]);
    },
    [consts.users.GET_USER]: (state, action) => {
      if (!action.ready || action.error) {
        return state;
      }
      return state.setIn(['users', action.payload.name], action.payload);
    },
    // We receive a list of organizations when a user lists the orgs they're in.
    // Add them to the correct user model.
    [consts.organizations.LIST_USER_ORGANIZATIONS]: (state, action) => {
      if (!action.ready || action.error) {
        return state;
      }
      // This sets a normalized result in user.orgs
      const {
        username,
        orgs
      } = action.payload;

      return state.setIn(['users', username, 'orgs'], orgs);
    }
};

export default function usersReducer(state = defaultState, action) {
    if (typeof actions[action.type] === 'function') {
        return actions[action.type](state, action);
    }
    return state;
}
