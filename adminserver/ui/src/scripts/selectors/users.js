'use strict';

import _ from 'lodash';
import { createSelector } from 'reselect';
import { getCurrentTeamMembers } from './teams';
import { Map } from 'immutable';
import { getAuthMethod } from './settings';

const baseUserSelector = state => state.users;

/**
 * Returns a user record for the current user logged in
 *
 * @return Record
 */
export const currentUserSelector = createSelector(
    baseUserSelector,
    users => ({...users.get('currentUser')})
);

/**
 * Returns a boolean for whether the current user is signed in.  Note this
 * returns true if there is no authentication method set up
 *
 * @return bool
 */
export const isLoggedInSelector = createSelector(
    baseUserSelector,
    getAuthMethod,
    (users, auth) => {
        return (users.get('currentUser').id !== 0)
            || auth === 'none'
            || users.get('isEmergencyAccess');
    }
);

/**
 * isAdminSelector returns whether the current user is an admin
 *
 * @return bool
 */
export const isAdminSelector = state => (state.users ? state.users.get('currentUser').isAdmin === true : undefined);

// User selectors
export const usersSelector = state => state.users.get('users', new Map()).toJS();
export const listMaxIdSelector = state => state.users.get('listMaxId');
export const searchTermSelector = state => state.ui.getIn(['selectMember', 'searchTerm']);

const visibleUsersSelector = createSelector(
  [usersSelector, searchTermSelector, listMaxIdSelector],
  (users, searchTerm, listMaxId) => {
    const searchTermLower = searchTerm.toLowerCase();
    return _.sortBy(Object.keys(users), (userName) => users[userName].id) // Need to sort by id since we paginate by id
      .filter((userName) => {
        if (userName.toLowerCase().indexOf(searchTermLower) === -1) {
          // Filter usernames not matching search term
          return false;
        }
        if (searchTerm === '') {
          // If not searching, we have to only show users with less than our max fetched list id
          // Otherwise when we load more then won't be a contiguous chunk at the bottom
          return users[userName].id <= listMaxId;
        }
        return true;
      })
      .map((userName) => {
        return users[userName];
      });
  }
);


/**
 * visibleNonTeamMembersSelector returns all users within visibleUsersSelector
 * that **aren't** part of the current team as defined in getCurrentTeamMembers
 *
 */
export const visibleNonTeamMembersSelector = createSelector(
  [visibleUsersSelector, getCurrentTeamMembers],
  (users, teamMembers) => {
    // Filter existing team members
    return users.filter((user) => !teamMembers.get(user.name));
  }
);

export const getUserOrgs = (state, props) => state.users.getIn(['users', props.username, 'orgs', 'entities', 'org'], new Map()).toJS();

export const selectUser = (state, props) => state.users.getIn(['users', props.params.username], {});

