'use strict';

import _ from 'lodash';
import { STATUS } from './common/Constants';
import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('CreateRepositoryFormStore');

var noErrorObj = {
  hasError: false,
  error: ''
};

export default createStore({
  storeName: 'CreateRepositoryFormStore',
  handlers: {
    CREATE_REPO_CLEAR_FORM: 'initialize',
    CREATE_REPO_UPDATE_FIELD_WITH_VALUE: '_updateFieldWithValue',
    CREATE_REPO_ATTEMPT_START: '_attemptStart',
    CREATE_REPO_BAD_REQUEST: '_badRequest',
    CREATE_REPO_SUCCESS: '_success',
    CREATE_REPO_FACEPALM: '_facepalm',
    CREATE_REPO_RECEIVE_NAMESPACES: '_receiveNamespaces'
  },
  initialize() {
    this.STATUS = STATUS.DEFAULT;

    this.namespaces = [];
    this.globalFormError = '';

    this.fields = {
      user: {},
      namespace: {},
      name: {},
      description: {},
      full_description: {},
      is_private: {}
    };

    this.values = {
      user: '',
      namespace: '',
      name: '',
      description: '',
      full_description: '',
      is_private: true
    };
  },
  _receiveNamespaces({
    namespaces, selectedNamespace
  }) {
    debug('receiving namespaces', namespaces, selectedNamespace);
    /**
     * namespaces is equivalent to the response in the namespaces API call
     */
    this.namespaces = namespaces.namespaces;
    if(_.includes(namespaces.namespaces, selectedNamespace)) {
      this.values.namespace = selectedNamespace;
    } else {
      this.values.namespace = namespaces.namespaces[0];
    }
    this.emitChange();
  },
  _facepalm() {
    // this happens if things are screwed and we can't recover gracefully
    this.STATUS = STATUS.FACEPALM;
    this.emitChange();
  },
  _clearForm() {
    this.initialize();
    this.emitChange();
  },
  _updateFieldWithValue({fieldKey, fieldValue}) {
    debug(fieldKey, fieldValue);
    this.fields[fieldKey].hasError = false;
    this.fields[fieldKey].error = '';
    this.globalFormError = '';
    this.values[fieldKey] = fieldValue;
    this.emitChange();
  },
  _attemptStart() {
    this.STATUS = STATUS.ATTEMPTING;
    this.emitChange();
  },
  _success() {
    this.STATUS = STATUS.SUCCESSFUL_SIGNUP;
    this.emitChange();
  },
  _badRequest(obj) {
    /**
     * This function expects keys which match the `this.fields` keys
     * with an array of errors:
     *
     * {
     *   orgname: ['this field is required']
     * }
     */
    let shouldEmitChange = false;
    // So far obj.detail is only returned when there are no more private repo's
    // We really need to update the response from the api
    if (obj.detail) {
      obj.is_private = [obj.detail];
    }

    // cycle through the possible form fields
    this.fields = _.mapValues(this.fields, function (errorObject, key) {
      if(_.has(obj, key)) {
        shouldEmitChange = true;
        return {
          hasError: !!obj[key],
          error: obj[key][0]
        };
      } else {
        return errorObject;
      }
    });
    /**
     * __all__ occurs when "Repository with this Name and Namespace already exists."
     */
    if(obj.__all__) {
      this.globalFormError = obj.__all__[0];
      shouldEmitChange = true;
    }

    if(shouldEmitChange) {
      this.emitChange();
    }
  },
  getState() {
    return {
      fields: this.fields,
      values: this.values,
      STATUS: this.STATUS,
      namespaces: this.namespaces,
      globalFormError: this.globalFormError
    };
  },
  dehydrate() {
    return this.getState();
  },
  rehydrate(state) {
    this.fields = state.fields;
    this.values = state.values;
    this.namespaces = state.namespaces;
    this.STATUS = state.STATUS;
    this.globalFormError = state.globalFormError;
  }
});
