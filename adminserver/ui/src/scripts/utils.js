'use strict';

import _ from 'lodash';
import { bindActionCreators } from 'redux';

/**
 * Given an object of actions, this returns a thunk which returns all actions
 * bound to dispatch using the same structure.
 *
 * This allows us to use `this.props.actions.$actionName` within components
 * after being connected to Redux.
 *
 * Example:
 *
 *   @connect(mapState, mapActions(Actions))
 *   class Basic extends Component {
 *     // this.props.actions now containers all keys in Actions bound to dispatch
 *   }
 *
 * You can also nest them to keep them organized, for example:
 *
 *   import * as TeamActions from 'actions/teams';
 *   import * as OrgActions from 'actions/organizations';
 *   @connect(mapState, mapActions({
 *     team: TeamActions,
 *     org: OrgActions,
 *   })
 *   class Basic extends Component {
 *     // this.props.actions.team.someTeamAction();
 *   }
 *
 */
export let mapActions = (actions) => {
  return (dispatch) => {
    const mapObject = (obj) => {
      return _.mapValues(obj, (val) => {
        if (typeof val === 'function') {
          return bindActionCreators(val, dispatch);
        } else if (typeof val === 'object') {
          return mapObject(val);
        }
        throw new Error('Invalid action object passed to mapActions');
      });
    };
    return {
      actions: mapObject(actions)
    };
  };
};

export let humanPermissions = (level) => {
  switch(level) {
    case 'admin':
      return 'Admin';
    case 'read-write':
      return 'Read & Write';
    case 'read-only':
      return 'Read only';
  }
};
let accessLevelRanks = {
  '': 1,
  'read-only': 2,
  'read-write': 3,
  'admin': 4
};
export let accessLevelCompare = (accessLevel1, accessLevel2) => {
  // TODO maybe check that keys are one of read/read-write/admin/empty
  return (accessLevelRanks[accessLevel1] || 0) - (accessLevelRanks[accessLevel2] || 0);
};
export let myGlobalAccessLevel = () => {
  if (window.user.isAdmin) {
    return 'admin';
  } else if (window.user.isReadWrite) {
    return 'read-write';
  } else if (window.user.isReadOnly) {
    return 'read-only';
  }
  return '';
};


export let fieldsWithErrorOnlyIfTouched = (fields) => {
  return Object.keys(fields).reduce((newFields, fieldName) => {
    if (!newFields[fieldName].touched) {
      newFields[fieldName].error = undefined;
    }
    return newFields;
  }, {...fields});
};

let kmExceptions = { 'props': '', '__esModule': '' };
export let keyMirror = function(obj) {
  let ret = {};
  if (!(obj instanceof Object) && !Array.isArray(obj)) {
    throw new Error('keyMirror(...): Argument must be an object or array.');
  }

  let keys = Array.isArray(obj) ? obj : Object.keys(obj);
  for (let key of keys) {
    ret[key] = key;
  }

  if (window.Proxy) {
    return new Proxy(ret, {
      get: function(target, name) {
        if ((name in target) || (name in kmExceptions)) {
          return target[name];
        }
        console.error(`Invalid key: ${name}`);
        throw new Error(`Invalid key: ${name}`);
      }
    });
  }
  return ret;
};

export let checkDuplicateConsts = (obj, used = {}) => {
  let subkeys = [];
  Object.keys(obj).forEach((key) => {
    const val = obj[key];
    if (val instanceof Object) {
      subkeys.push(key);
    } else {
      if (key in used) {
        throw new Error(`Duplicate constant found: ${key}`);
      }
      used[key] = true;
    }
  });
  subkeys.forEach((key) =>
    checkDuplicateConsts(obj[key], used)
  );
};

export let formatError = (error) => {
  if (!error) {
    return '';
  }

  if (typeof error === 'string' || error instanceof String) {
    return error;
  }

  if (error instanceof Error) {
    // A javascript error
    return error.toString();
  }

  if (error instanceof Object) {
    if (Array.isArray(error)) {
      // Probably an array of errors returned from our API
      return error.map((e) => formatError(e)).join();
    }

    if (error.toJS) {
      // Probably an array of errors converted to immutable js when merged into our state tree
      return formatError(error.toJS());
    }

    if (error.data && error.data.errors) {
      // Probably an unhandled rejected axios request that contain API errors
      return formatError(error.data.errors);
    }

    if (error.code && error.message) {
      // Probably is in the format of our API errors TODO: Don't ignore error.details
      return `${error.code}: ${error.message}`;
    }
  }

  return JSON.stringify(error) || '';
};

/*
    will return 0 or the page as an int pulled from location.search
*/
export const getPage = (search) => {
    return search.substring(1).split('&').map((param) => {
        const params = param.split('=');
        if (params[0] === 'page') {
            return parseInt(params[1]);
        }
    // user facing page numbers start at one so subtract 1
    })[0] - 1 || 0;
};
