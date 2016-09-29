'use strict';

import createStore from 'fluxible/addons/createStore';
import _ from 'lodash';
var debug = require('debug')('createFormStore');
/**
 * @param {Array} fields - array of objects with field names
 *   as keys and inital values as values
 * @param {Function} init - old initialize function
 */

export default function createFormStore(fields, oldSpec) {

  var spec = {};

  spec.initialize = function() {
    this.globalFormError = '';
    _.forOwn(fields, function(key, val){
      this.fields[key] = {};
      this.values[key] = val;
    });
    spec.initialize();
  }.bind(spec);

  spec._badRequest = function(obj) {
    /**
     * obj is an Object with keys that are field names
     *     and values that are arrays of errors
     *
     * This function should be used as the handler for
     * an HTTP 400 BadRequest
     *
     * obj = {
     *   username: ["cannot be empty"]
     * }
     */
    // did we update state?
    var dirty = false;

    _.forOwn(obj, function(key, val) {
      if(_.includes(fields, key)) {
        this.fields[key].hasError = !!val;
        this.fields[key].error = val[0];
        dirty = true;
      }
    }, this);

    if(dirty) {
      this.emitChange();
    }
  };

  spec._getState = function() {
    return {
      fields: this.fields,
      values: this.values,
      globalFormError: this.globalFormError
    };
  };


  spec._updateFieldWithValue = function({fieldKey, fieldValue}){
    this.values[fieldKey] = fieldValue;
    this.emitChange();
  };


  spec.dehydrate = function() {
    return {};
  },
  spec.rehydrate = function(state) {
    this.state = state;
  };

  _.merge(spec, oldSpec, function(objectValue, sourceValue, key) {
    if(key === 'initializer') {
      return sourceValue;
    } else if(key === 'getState') {
      return function() {
        debug('state', this.state);
        return _.merge({},
                       objectValue.getState(),
                       sourceValue._getState());
      };
    }
  });
  return createStore(spec);
}
