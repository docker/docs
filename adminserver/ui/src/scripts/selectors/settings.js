'use strict';
import { createSelector } from 'reselect';
import { Map } from 'immutable';

export const settingsMap = (state) => {
    return state.settings;
};

export const getSettings = (state) => state.settings.toJS();

export const getDomainName = (state) => {
  // Change in 2.0.0: domain name is now dtrHost and includes ports.
  // Derp.
  const host = state.settings.get('dtrHost', '');
  let [domain, port] = host.split(':');
  if (port === '80' || port === '443') {
    return domain;
  }
  return host;
};

export const getAuthMethod = createSelector(
    settingsMap,
    settings => settings.getIn(['enzi', 'backend'], 'managed')
);

// by default upgrades are disabled
export const getEnabledUpgrades = state => !state.settings.getIn(['http', 'disableUpgrades'], true);

export const getRegistry = state => state.settings.get('registry');

// when we save storage settings via a form we need these properties
export const getNonStorageRegistry = createSelector(
    getRegistry,
    (registrySettings) => {
        const searchKeys = ['maintenance', 'cache', 'delete', 'redirect'];
        return registrySettings.get('storage').filter((value, key) => {
            return searchKeys.indexOf(key) !== -1;
        }).toJS();
    }
);

export const getLdapStatus = state => state.status.get('LDAP_CHECK');

export const getAuthSettings = state => state.settings.get('auth');
export const getLdapSettings = state => state.settings.get('ldap', new Map()).toJS();

// since the storage settings key can be multiple values
// filter the opposite of getNonStorageRegistry
export const getStorageRegistry = createSelector(
    getRegistry,
    (registrySettings) => {
        const searchKeys = ['maintenance', 'cache', 'delete', 'redirect'];
        return registrySettings.get('storage', {}).filter((value, key) => {
            return searchKeys.indexOf(key) === -1;
        }).toJS();
    }
);

export const getStorageType = createSelector(
    getStorageRegistry,
    (storageSettings) => {
        return Object.keys(storageSettings)[0];
    }
);

export const getStorageYAML = createSelector(
  getRegistry,
  (reg) => reg.get('config', '')
);

export const getUpdates = state => state.updates;

export const getLicense = state => {
  const settings = state.settings.get('license', new Map());
  if (settings === null ) {
    return {};
  }
  return settings.toJS();
};

export const getGCSettings = state => state.settings.get('gc', new Map());
export const gcSchedule = createSelector(
    getGCSettings,
    (gcSettings) => {
        return gcSettings.get('schedule');
    }
);
export const gcTimeout = createSelector(
    getGCSettings,
    (gcSettings) => {
        return gcSettings.get('timeout');
    }
);
export const gcLastRun = createSelector(
    getGCSettings,
    (gcSettings) => {
        return gcSettings.get('lastRun', new Map()).toJS();
    }
);
export const gcRunning = createSelector(
    getGCSettings,
    (gcSettings) => {
        return gcSettings.get('running');
    }
);

export const getSSLCert = state => state.settings.get('SSL');
