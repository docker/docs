'use strict';

import { createSelector } from 'reselect';
import { Map } from 'immutable';
import { OrganizationRecord } from 'records';

import {
  currentUserSelector,
  isAdminSelector
} from './users.js';


// INPUT SELECTORS
// ===============

// Returns the 'org' state map
const getOrgState = state => state.organizations.get('org', new Map());

// Returns a Map of organization records
const organizationsSelector = state =>
  state.organizations
    .getIn(['orgs', 'entities', 'org'], new Map())
    .filter((org) => org.name !== 'docker-datacenter');

export const getOrgMembers = (state, props) => state.organizations.getIn(['orgMembers', props.params.org, 'members'], new Map()).toJS();


// COMPOSED SELECTORS
// ==================

// Returns all organization names as an array
export const orgNamesSelector = createSelector(
  [organizationsSelector],
  (orgs) => orgs.map(o => o.name).toArray()
);

/**
 * Returns the currently viewed organization within the router props.
 * Used within `scenes/organization`
 *
 * @return OrganizationRecord
 */
export const getPageOrg = createSelector(
  getOrgState,
  (state, props) => props.params ? props.params.org : '',
  (orgState, name) => {
    return orgState.getIn(['entities', 'org', name], new OrganizationRecord());
  }
);

/**
 * getIsCurrentUserPageOrgAdmin returns whether the current logged in user has
 * admin rights over the currently visible page, as defined by the router.
 *
 * If the user is a global admin this always returns true.  Otherwise, we check
 * whether the user is in the 'owners' team of the organization.
 *
 * @return bool
 */
export const getIsCurrentUserPageOrgAdmin = createSelector(
  (state, props) => state.organizations.getIn(['orgMembers', props.params.org, 'members'], new Map()),
  currentUserSelector,
  isAdminSelector,
  (orgMembers, user, isGlobalAdmin) => orgMembers.getIn([user.name, 'isAdmin'], false) || isGlobalAdmin
);
