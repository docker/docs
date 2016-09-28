'use strict';
import createStore from 'fluxible/addons/createStore';
import findIndex from 'lodash/array/findIndex';
import map from 'lodash/collection/map';
import { STATUS } from './common/Constants';

var AutobuildTriggerByTagStore = createStore({
  storeName: 'AutobuildTriggerByTagStore',
  handlers: {
    INITIALIZE_AB_TRIGGERS: '_initTriggers',
    AB_TRIGGER_BY_TAG_ERROR: '_triggerByTagError',
    AB_TRIGGER_BY_TAG_SUCCESS: '_triggerByTagSuccess',
    ATTEMPT_TRIGGER_BY_TAG: '_triggerByTagAttempt'
  },
  initialize: function() {
    this.triggers = [];
    this.tagStatuses = [];
  },
  _initTriggers: function(tags) {
    //on load of the build settings page
    this.initialize();

    this.triggers = map(tags, (tag) => {
      return {
        id: tag.id,
        success: '',
        error: ''
      };
    });

    this.tagStatuses = map(tags, (tag) => {
      return {
        id: tag.id,
        status: STATUS.DEFAULT
      };
    });
    this.emitChange();
  },
  _findIndices: function(id) {
    const statusIndex = findIndex(this.tagStatuses, (s) => {
      return s.id === id;
    });
    const triggerIndex = findIndex(this.triggers, (t) => {
      return t.id === id;
    });
    return {statusIndex, triggerIndex};
  },
  _triggerByTagAttempt: function(id) {
    const {statusIndex, triggerIndex} = this._findIndices(id);
    this.tagStatuses[statusIndex].status = STATUS.ATTEMPTING;
    this.triggers[triggerIndex].error = '';
    this.triggers[triggerIndex].success = '';
    this.emitChange();
  },
  _triggerByTagError: function(errObj) {
    const {id, error} = errObj;
    const { triggerIndex } = this._findIndices(id);
    this.triggers[triggerIndex].error = error;
    setTimeout(this._clearTriggerStatus.bind(this, id), 3000);
    this.emitChange();
  },
  _triggerByTagSuccess: function(successObj) {
    const {id, success} = successObj;
    const { triggerIndex } = this._findIndices(id);
    this.triggers[triggerIndex].success = success;
    setTimeout(this._clearTriggerStatus.bind(this, id), 3000);
    this.emitChange();
  },
  _clearTriggerStatus: function(id) {
    const {statusIndex, triggerIndex} = this._findIndices(id);
    this.tagStatuses[statusIndex].status = STATUS.DEFAULT;
    this.triggers[triggerIndex].error = '';
    this.triggers[triggerIndex].success = '';
    this.emitChange();
  },
  getState: function() {
    return {
      triggers: this.triggers,
      tagStatuses: this.tagStatuses
    };
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this.triggers = state.triggers;
    this.tagStatuses = state.tagStatuses;
  }
});

module.exports = AutobuildTriggerByTagStore;
