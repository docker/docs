'use strict';

import { Map } from 'immutable';
import consts from 'consts';
import { TeamRecord } from 'records';

const defaultState = new Map({
  // teamsByOrgName is a Map of org names to normalized listTeams responses.
  teamsByOrgName: new Map(),

  // teamMembers is a map of orgName => teamName => normalizr data for user records
  teamMembers: new Map(),

  // a map of teams by org for a specific user
  teamsForUser: new Map(),

  // membershipCheck is used from the `checkTeamMembership` call to store
  // whether the user is in a team. This is stored separately from teamMembers
  // because teamMembers stores data from the LIST_MEMBERS API call as returned
  // from normalizr.
  //
  // Map of orgName => teamName => memberName => {bool}
  membershipCheck: new Map()
});

const actions = {
  [consts.teams.LIST_TEAMS]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    const { orgName, data } = action.payload;
    return state.setIn(['teamsByOrgName', orgName], data);
  },
  [consts.teams.CREATE_TEAM]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }

    const {
      payload: {
        orgName,
        team
      }
    } = action;

    const entityPath = ['teamsByOrgName', orgName, 'entities', 'team', team.name];

    return state
      .setIn(entityPath, team)
      .updateIn(['teamsByOrgName', orgName, 'result'], list => list.push(team.name));
  },

  [consts.teams.GET_TEAM]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }

    const { orgName } = action.payload[0];

    // the second payload has a name
    // the name is set to the team's ID
    // the record picks up both, and uses the second to overwrite the first
    // this delete simply removes the second payload's name
    // so the resulting record has the correct datum
    delete action.payload[1].team.name;

    const record = new TeamRecord({
      ...action.payload[0].team,
      ...action.payload[1].team
    });

    return state
      // for display on the repo permissions page this is much easier
      // TODO Tectonic
      .setIn(['teamsById', orgName, record.id], record)
      .setIn(['teamsByOrgName', orgName, record.name], record);
  },

  [consts.teams.UPDATE_TEAM]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    const {
      orgName,
      team,
      team: {
        name: teamName
      }
    } = action.payload;
    return state
      .setIn(['teams', orgName, teamName], team)
      .setIn(['teamsByOrgName', orgName, teamName], team);
  },

  [consts.teams.DELETE_TEAM]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    const {
      orgName,
      teamName
    } = action.payload;
    return state
      .deleteIn(['teams', orgName, teamName])
      .deleteIn(['teamMembers', orgName, teamName]);
  },

  /** Member actions **/

  [consts.teams.LIST_MEMBERS]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    const { orgName, teamName, data } = action.payload;
    return state.setIn(['teamMembers', orgName, teamName], data.getIn(['entities', 'user'], new Map()));
  },

  [consts.teams.ADD_TEAM_MEMBER]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }

    const { name, team, member } = action.payload;
    return state.mergeIn(['teamMembers', name, team], {
      [member.name]: {
        isAdmin: member.Admin || false,
        member
      }
    });
  },

  [consts.teams.ADD_TEAM_MEMBERS]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    const { orgName, teamName, members } = action.payload;
    return state.mergeIn(['teamMembers', orgName, teamName], members);
  },

  [consts.teams.GET_TEAM_MEMBER]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    const { orgName, teamName, memberName, isMember } = action.payload;
    return state.setIn(['membershipCheck', orgName, teamName, memberName], isMember);
  },

  [consts.teams.DELETE_TEAM_MEMBER]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    const { orgName, teamName, memberName } = action.payload;
    return state.withMutations(s => {
      // Remove the member from the 'user' entity map within the teamMembers map
      s.setIn(['teamMembers', orgName, teamName], s.getIn(['teamMembers', orgName, teamName]).filter((userRecord) => {
        return userRecord.get('member').name !== memberName;
      }));
    });
  },

  // This is listened to in both reducers/organizations.js and here
  // This removes the deleted member from all teams
  [consts.organizations.DELETE_MEMBER]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    const { orgName, memberName } = action.payload;
    const teamMembers = state.getIn(['teamMembers', orgName]);
    if (!teamMembers) {
      return state;
    }
    return state.setIn(
      ['teamMembers', orgName],
      teamMembers.withMutations((mutableMembers) => {
        for (let teamName of mutableMembers.keys()) {
          mutableMembers.deleteIn([teamName, memberName]);
        }
      })
    );
  },
  [consts.teams.GET_TEAMS_FOR_USER]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    return state.setIn(['teamsForUser', action.payload.user], action.payload.teamsByOrg);
  }

};

export default function teamReducer(state = defaultState, action) {
  if (typeof actions[action.type] === 'function') {
    return actions[action.type](state, action);
  }
  return state;
}
