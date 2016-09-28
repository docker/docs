'use strict';

import { Map, fromJS } from 'immutable';
import consts from 'consts';

const { settings } = consts;

const actions = {
    [settings.LDAP_SETTINGS]: (state, action) => {
      if (!action.ready || action.error) {
        return state;
      }
      return state.set('ldap', action.payload);
    },
    [settings.ALL_SETTINGS]: (state, action) => {
        if (!action.ready || action.error) {
            return state;
        }
        return state.merge(action.payload);
    },
    [settings.AUTH_SETTINGS]: (state, action) => {
        if (!action.ready || action.error) {
            return state;
        }
        return state.set('enzi', action.payload);
    },
    [settings.GET_STORAGE_SETTINGS]: (state, action) => {
      if (!action.ready || action.error) {
        return state;
      }
      return state.set('registry', action.payload);
    },
    [settings.GET_GC_SCHEDULE]: (state, action) => {
        if (!action.ready || action.error) {
            return state;
        }
        return state.setIn(['gc', 'schedule'], action.payload.data.schedule).setIn(['gc', 'timeout'], '' + action.payload.data.timeout);
    },
    [settings.GET_LAST_GC_SAVINGS]: (state, action) => {
        if (!action.ready || action.error) {
            return state;
        }
        return state.setIn(['gc', 'lastRun'], fromJS(action.payload));
    },
    [settings.GET_GC_STATUS]: (state, action) => {
        if (!action.ready || action.error) {
            return state;
        }
        return state.setIn(['gc', 'running'], action.payload.data.registryGC.status === 'running');
    },
    [settings.GET_LICENSE]: (state, action) => {
      if (!action.ready || action.error) {
        return state;
      }
      return state.set('license', fromJS(action.payload.data));
    },
    [settings.LICENSE_SETTINGS_SAVE]: (state, action) => {
      if (!action.ready || action.error) {
        return state;
      }
      return state.set('license', fromJS(action.payload.data));
    }
};

let defaultState = new Map({
  // stores all settings from /v0/meta/settings - DTR settings
  settings: new Map(),
  // stores all settings for enzi
  enzi: new Map(),
  // stores enzi ldap settings
  ldap: new Map(),
  registry: new Map({
    storage: new Map()
  }),
  gc: new Map({
    schedule: '',
    timeout: '0',
    lastRun: new Map({
      time: '',
      layers: new Map()
    }),
    running: false
  })
});

export default function(state = defaultState, data) {
  if (typeof actions[data.type] === 'function') {
    return actions[data.type](state, data);
  }
  return state;
}
