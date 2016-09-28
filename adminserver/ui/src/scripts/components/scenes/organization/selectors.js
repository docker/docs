'use strict';

import { createSelector } from 'reselect';
import { Map } from 'immutable';

/**
 * getTeams returns teams
 * for the current org (determined by router param
 * ordered alphabetically
 *
 * @return array Array of TeamRecords
 */
export const getTeams = createSelector(
  (state, props) => props.params.org,
  (state) => state.teams.get('teamsByOrgName'),

  (orgName, teamsByOrgName) => {
    return teamsByOrgName
      .getIn([orgName, 'entities', 'team'], new Map())
      .toArray()
      .sort((a, b) => a.name.charCodeAt(0) - b.name.charCodeAt(0));
  }
);

/**
 * Returns all teams for an org excluding the 'owners' team
 *
 * @return array Array of TeamRecords
 */
export const getTeamsWithoutOwners = createSelector(
  getTeams,
  (teams) => teams.filter(t => t.name !== 'owners')
);
