'use strict';

import consts from 'consts';
import { Map } from 'immutable';

const defaultState = new Map({
  // Data from Normalizr
  orgs: new Map(),
  // Data from Normalizr for a single org
  org: new Map(),
  // Map of org names to normalizr data for members of a single ord
  orgMembers: new Map()
});

/**
 * This replaces all state within state.organizations with action.payload
 */
const replaceOrganizations = (state, action) => {
  if (!action.ready || action.error) {
    return state;
  }

  return state.set('orgs', action.payload);
};

const actions = {
  [consts.organizations.LIST_USER_ORGANIZATIONS]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    return replaceOrganizations(state, action.payload.orgs);
  },
  [consts.organizations.LIST_ORGANIZATIONS]: replaceOrganizations,
  [consts.organizations.GET_ORGANIZATION]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    return state.set('org', action.payload);
  },

  [consts.organizations.CREATE_ORGANIZATION]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    const { org } = action.payload;
    return state
      .setIn(['orgs', 'entities', 'org', org.name], org)
      .updateIn(['orgs', 'result'], list => list.push(org.name));
  },

  [consts.organizations.DELETE_ORGANIZATION]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    // TODO
    return state;
  },

  [consts.organizations.LIST_ORGANIZATION_MEMBERS]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    return state.setIn(['orgMembers', action.payload.orgName, 'members'], action.payload.data.getIn(['entities', 'user'], new Map()));
  },

  [consts.organizations.DELETE_MEMBER]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    const { orgName, memberName } = action.payload;
    return state.withMutations(s => {
      // Remove the member from the org's list of members
      s.deleteIn(['orgMembers', orgName, 'members', memberName]);
    });
  },

  [consts.organizations.ADD_MEMBER]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    const {
      payload: {
        member,
        orgName
      }
    } = action;
    return state.setIn(['orgMembers', orgName, 'members', member.get('name')], member);
  },

  [consts.organizations.UPDATE_MEMBER]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    const {
      payload: {
        member,
        orgName
      }
    } = action;

    let memberData = member.getIn(['entities', 'user']).toJS();
    let memberName = Object.keys(memberData)[0];

    return state.setIn(['orgMembers', orgName, 'members', memberName], memberData[memberName]);
  },

  [consts.organizations.CREATE_ADD_MEMBER]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    const {
      payload: {
        member,
        orgName
      }
    } = action;
    return state.setIn(['orgMembers', orgName, 'members', member.get('name')], member);
  }
};

export default function organizations(state = defaultState, action) {
  if (typeof actions[action.type] === 'function') {
    return actions[action.type](state, action);
  }
  return state;
}
