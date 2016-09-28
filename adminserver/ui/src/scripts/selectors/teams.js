'use strict';

import { createSelector } from 'reselect';
import { Map } from 'immutable';

/**
 * Returns the raw state.teams map
 */
export const getRawTeamState = state => state.teams;
/**
 * This returns the raw teams.teamMembers immutable map
 *
 * @return Immutable.Map state.teams.teamMembers
 */
export const getRawTeamMembers = state => state.teams.getIn(['teamMembers']);

/**
 * Returns the raw teams.membershipCheck immutable map
 *
 * @return Immutable.Map state.teams.membershipCheck
 */
export const getRawMembershipCheck = state => state.teams.getIn(['membershipCheck']);

/**
 * getCurrentTeamMembers returns a map of all users within the current team
 * (as defined from router parameters).  Note that the array will contain
 * UserRecords.
 *
 * @return Immutable.Map Map of UserRecords
 */
export const getCurrentTeamMembers = createSelector(
  getRawTeamMembers,
  (state, props) => props.params.org,
  (state, props) => props.params.team,
  (teamMembers, orgName, teamName) => {
    return teamMembers.getIn([orgName, teamName], new Map()).toJS();
  }
);

export const getIsOrgOwner = (state, props) => state.teams.getIn(
    ['teamMembers', props.params.namespace, 'owners', window.user.name]
);

export const teamDetails = (state, props) => {
  return state.teams.getIn(['teamsById', props.org, props.id], new Map());
};

export const teamsForUser = (state, props) => {
  return state.teams.getIn(['teamsForUser', props.params.username], []);
};

export const teamDetailsByName = (state, props) => {
  return state.teams.getIn(['teamsByOrgName', props.params.org, props.params.team], new Map());
};
