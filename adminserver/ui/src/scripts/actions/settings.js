'use strict';

import { Settings, LDAP, Updates, Enzi, Jobs } from 'dtr-js-sdk';
import consts from 'consts';
import { fromJS } from 'immutable';

const {
    settings
} = consts;

/**
 * getAuthSettings returns enzi auth settings
 */
export function getAuthSettings() {
  return {
    type: settings.AUTH_SETTINGS,
    meta: {
      promise: Enzi.default.getAuthConfig().then((response) => {
        return fromJS(response.data);
      })
    }
  };
}

export function saveAuthSettings({ backend }) {
  return {
    type: settings.SAVE_AUTH_SETTINGS,
    meta: {
      promise: Enzi.default.updateAuthConfig({}, { backend }),
      notifications: {
        pending: 'Saving settings',
        success: 'Settings saved successfully',
        error: 'Unable to save settings'
      }
    }
  };
}

/**
 * getLdapSettings returns enzi auth settings
 */
export function getLdapSettings() {
  return {
    type: settings.LDAP_SETTINGS,
    meta: {
      promise: Enzi.default.getLdapConfig().then((response) => {
        return fromJS(response.data);
      })
    }
  };
}

/**
 * getLdapSettings returns enzi auth settings
 */
export function saveLdapSettings(data) {
  return (dispatch) => ({
    type: settings.LDAP_SETTINGS,
    meta: {
      promise: Enzi.default.updateLdapConfig({}, data).then(() => {
        dispatch(saveAuthSettings({ backend: 'ldap' }));
      }),
      notifications: {
        pending: 'Saving LDAP settings',
        success: 'LDAP settings saved successfully',
        error: 'Unable to save LDAP settings'
      }
    }
  });
}

export function getSettings() {
    return {
        type: settings.ALL_SETTINGS,
        meta: {
            promise: Settings.getAllSettings().then((response) => {
                return fromJS(response.data);
            })
        }
    };
}

export function getStorageSettings() {
    return {
        type: settings.GET_STORAGE_SETTINGS,
        meta: {
            promise: Settings.getStorageSettings().then((response) => {
                return fromJS(response.data);
            })
        }
    };
}

export function getUpdates() {
    return {
        type: settings.UPDATES,
        meta: {
            promise: Updates.getUpdates().then((response) => {
                return response.data;
            })
        }
    };
}

export function saveSettings(inputSettings) {
  return {
    type: settings.GENERAL_SETTINGS_SAVE,
    meta: {
      promise: Settings.updateSettings(inputSettings),
      notifications: {
        pending: 'Saving settings',
        success: 'Settings saved successfully',
        error: (resp) => {
          return `Unable to save settings: ${resp.data.errors[0].detail}`;
        }
      }
    }
  };
}

export function getLicenseSettings() {
  return {
    type: settings.GET_LICENSE,
    meta: {
      promise: Settings.getLicense()
    }
  };
}

export function saveLicense(license) {
  return {
    type: settings.LICENSE_SETTINGS_SAVE,
    meta: {
      promise: Settings.updateLicense(license),
      notifications: {
        pending: 'Saving license',
        success: 'License saved successfully',
        error: 'There was an error saving your license'
      }
    }
  };
}

export function autorefreshLicense(autorefresh) {
  return {
    type: settings.TOGGLE_LICENSE_AUTOREFRESH,
    meta: {
      promise: Settings.autorefreshLicense({ auto_refresh: autorefresh })
    }
  };
}


/**
 * This accepts a single YAML string and updates registry settings via
 * an API call
 */
export function saveYamlStorage(yaml) {
  return {
    type: settings.YAML_STORAGE_SAVE,
    meta: {
      promise: Settings.updateStorageYaml(yaml),
      notifications: {
        pending: 'Saving storage settings',
        success: 'Settings saved',
        error: (resp) => {
          return `An error occured: ${resp.data.message}`;
        }
      }
    }
  };
}

export function saveFormStorage(data) {
  return {
    type: settings.FORM_STORAGE_SAVE,
    meta: {
      promise: Settings.updateStorage(data),
      notifications: {
        pending: 'Saving storage settings',
        success: 'Settings saved',
        error: (resp) => resp.data.message
      }
    }
  };
}

export function fetchStorageOptions() {
  return {
    type: settings.FETCH_SETTINGS,
    meta: {
      promise: Settings.getStorageDrivers()
    }
  };
}

export function clearRegistryStatus() {
  return { type: settings.REGISTRY_CLEAR };
}

export function checkLdapSettings(data) {
  return {
    type: settings.LDAP_CHECK,
    meta: {
      promise: LDAP.verifyLdapSettings(data)
    }
  };
}

export function syncLdap() {
  return {
    type: settings.LDAP_SYNC,
    meta: {
      promise: Enzi.default.createJob({}, { action: 'ldap-sync' }),
      notifications: {
        pending: 'Starting LDAP sync',
        success: 'LDAP sync started',
        error: 'There was an error starting an LDAP sync'
      }
    }
  };
}

export function getGCSchedule() {
  return {
    type: settings.GET_GC_SCHEDULE,
      meta: {
          promise: Settings.getGCSchedule()
      }
  };
}

export function updateGCSchedule(data) {
  let successString = `Garbage collection set for '${data.schedule}'`;
  if (data.schedule === '') {
      successString = 'Garbage collection schedule removed';
  }

  return {
    type: settings.UPDATE_GC_SCHEDULE,
    meta: {
      promise: Settings.updateGCSchedule(data),
      notifications: {
        pending: 'Updating garbage collection schedule',
        success: successString,
        error: (resp) => {
          return resp.data.message;
        }
      }
    }
  };
}

export function deleteGCSchedule() {
  return {
    type: settings.DELETE_GC_SCHEDULE,
    meta: {
      promise: Settings.deleteGCSchedule(),
      notifications: {
        pending: 'Removing the garbage collection schedule',
        success: 'Garbage collection schedule removed',
        error: (resp) => {
          return resp.data.message;
        }
      }
    }
  };
}

export function getLastGCSavings() {
  return {
    type: settings.GET_LAST_GC_SAVINGS,
    meta: {
      promise: Settings.getLastGCSavings().then(resp => {
        if (resp.data.layers === undefined) {
          return { layers: [] };
        }
        return resp.data;
      })
    }
  };
}

export function runGC() {
  return {
    type: settings.RUN_GC,
    meta: {
      promise: Jobs.createJob({job: 'registryGC'}),
      notifications: {
        pending: 'Attempting to run garbage collection',
        success: 'Garbage collection running',
        error: (resp) => {
          return resp.data.message;
        }
      }
    }
  };
}

export function stopGC() {
  return {
    type: settings.STOP_GC,
    meta: {
      promise: Settings.stopGC()
    }
  };
}

export function getGCStatus() {
  return {
    type: settings.GET_GC_STATUS,
    meta: {
      promise: Jobs.getJobs()
    }
  };
}
