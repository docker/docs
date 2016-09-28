'use strict';

import { createSelector } from 'reselect';
import { Map } from 'immutable';

/**
 * getOrgMembers uses the org name from the router and returns the normalized
 * org member data from state.organizations.get('orgMembers')
 *
 * @return Map normalized response of org members API call
 */
const getOrgMembers = createSelector(
  (state, props) => props.params.org,
  state => state.organizations.get('orgMembers'),

  (orgName, orgMemberMap) => {
    return orgMemberMap.getIn([orgName, 'members'], new Map());
  }
);

/**
 * getTeamMembers returns normalized API response listing team members for the
 * current team/org params in the router
 *
 * @return Map normalzied response of team members API call
 */
const getTeamMembers = createSelector(
  (state, props) => props.params.org,
  (state, props) => props.params.team,
  state => state.teams.get('teamMembers'),
  (orgName, teamName, teamMemberMap) => {
    return teamMemberMap.getIn([orgName, teamName], new Map());
  }
);

/**
 * getMembers returns the members of the current team. If no team is
 * selected, this returns members of the current organization.
 *
 * This is used in combination with the `listOrgOrTeamMembers` action which uses
 * the same team predicate.
 *
 * @return array  Array of UserRecords
 */
export const getMembers = createSelector(
  (state, props) => props.params.team,
  getOrgMembers,
  getTeamMembers,

  (teamName, orgMembers, teamMembers) => {
    if (teamName === undefined) {
      // Return the organization's members.
      return orgMembers.getIn(['entities', 'user'], new Map()).toArray();
    }
    return teamMembers.getIn(['entities', 'user'], new Map()).toArray();
  }
);
