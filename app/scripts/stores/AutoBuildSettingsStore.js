'use strict';
import createStore from 'fluxible/addons/createStore';
import sortBy from 'lodash/collection/sortBy';
import { STATUS } from './common/Constants';

var AutoBuildSettingsStore = createStore({
  storeName: 'AutoBuildSettingsStore',
  handlers: {
    AB_TRIGGER_BY_TAG_ERROR: '_triggerByTagError',
    AB_TRIGGER_BY_TAG_SUCCESS: '_triggerByTagSuccess',
    ATTEMPT_TRIGGER_BY_TAG: '_triggerByTagAttempt',
    RECEIVE_AUTOBUILD_SETTINGS: '_receiveAutoBuildSettings',
    UPDATE_AUTO_BUILD_SETTINGS: '_updateFields',
    RECEIVE_AUTOBUILD_LINKS: '_receiveAutoBuildLinks',
    LINK_AUTOBUILD_ERROR: '_linkAutoBuildError',
    LINK_AUTOBUILD_SUCCESS: '_linkAutoBuildSuccess',
    UPDATE_AUTOBUILD_PUSH_TRIGGER_ITEM: '_updateBuildTriggerItem',
    UPDATE_AUTOBUILD_NEW_TAG_ITEM: '_updateBuildNewTagItem',
    DELETE_AUTOBUILD_PUSH_TRIGGER_ITEM: '_deleteBuildTriggerItem',
    DELETE_AUTOBUILD_NEW_TAG_ITEM: '_deleteBuildsNewTagItem',
    ADD_AUTOBUILD_PUSH_TRIGGER_ITEM: '_addBuildTriggerItem',
    SAVE_BUILD_TAGS_SUCCESS: '_saveTagsSuccess',
    SAVE_BUILD_TAGS_ERROR: '_saveTagsError',
    RECEIVE_TRIGGER_STATUS: '_receiveTriggerStatus',
    RECEIVE_TRIGGER_LOGS: '_receiveTriggerLogs'
  },
  initialize: function() {
    this.autoBuildStore = {
      repository: '',
      build_name: '',
      provider: '',
      source_url: '',
      docker_url: '',
      repo_web_url: '',
      repo_type: '',
      active: false,
      deleted: false,
      repo_id: '',
      build_tags: [],
      deploykey: null,
      hook_id: ''
    };
    this.newTags = [];
    this.validations = {
      buildTags: {
        hasError: false,
        success: false,
        errors: []
      },
      links: {
        hasError: false,
        success: false,
        error: ''
      },
      trigger: {
        success: '',
        error: ''
      }
    };
    this.autoBuildBlankSlate = {};
    this.autoBuildLinks = [];
    this.triggerLinkForm = {
      repoName: ''
    };
    this.triggerStatus = {
      token: '',
      trigger_url: '',
      active: false
    };
    this.triggerLogs = [];
    this.STATUS = STATUS.DEFAULT;
  },
  _resetValidations: function(field) {
    if (field === 'buildTags') {
      this.validations.buildTags = {
        hasError: false,
        success: false,
        errors: []
      };
    } else {
      this.validations[field] = {
        hasError: false,
        success: false,
        error: ''
      };
    }
    this.emitChange();
  },
  _receiveAutoBuildSettings: function(payload) {
    this.autoBuildStore = payload;
    const sorted = sortBy(payload.build_tags, 'id'); // ensure build_tags received are sorted
    this.autoBuildStore.build_tags = sorted;
    this.autoBuildBlankSlate = this.autoBuildStore;
    this.emitChange();
  },
  _receiveAutoBuildLinks: function(payload) {
    this.autoBuildLinks = payload;
    this.triggerLinkForm.repoName = '';
    this.emitChange();
  },
  _linkAutoBuildError: function() {
    this.validations.links = {
      hasError: true,
      success: false,
      error: 'Failed to link this repository to your Automated Build.'
    };
    this.emitChange();
  },
  _linkAutoBuildSuccess: function() {
    this.validations.links = {
      hasError: false,
      success: true,
      error: ''
    };
    this.emitChange();
  },
  _addBuildTriggerItem: function() {
    this.newTags.push({
      name: '',
      dockerfile_location: '',
      source_name: '',
      source_type: 'Branch',
      isNew: true
    });
    this._resetValidations('buildTags');
    this.emitChange();
  },
  _deleteBuildTriggerItem: function(index) {
    this.autoBuildStore.build_tags[index].toDelete = true;
    this._resetValidations('buildTags');
    this.emitChange();
  },
  _deleteBuildsNewTagItem: function(index) {
    this.newTags[index].toDelete = true;
    this._resetValidations('buildTags');
    this.emitChange();
  },
  _updateBuildTriggerItem: function({ index, fieldkey, value}) {
    this.autoBuildStore.build_tags[index][fieldkey] = value;
    this._resetValidations('buildTags');
    this.emitChange();
  },
  _updateBuildNewTagItem: function({ index, fieldkey, value}) {
    this.newTags[index][fieldkey] = value;
    this._resetValidations('buildTags');
    this.emitChange();
  },
  _updateFields: function({ field, key, value }) {
    this[field][key] = value;
    if (field === 'triggerLinkForm') {
      this._resetValidations('links');
    }
    this.emitChange();
  },
  _saveTagsSuccess: function() {
    this.validations.buildTags = {
      success: true,
      hasError: false,
      errors: []
    };
    this.newTags = [];
    setTimeout(this._resetValidations.bind(this), 3000, 'buildTags');
    this.emitChange();
  },
  _saveTagsError: function(tag) {
    let currentErrors = this.validations.buildTags.errors;
    if (tag.error) {
      currentErrors.push(`${tag.name}: ${tag.error}`);
    }
    this.validations.buildTags = {
      success: false,
      hasError: true,
      errors: currentErrors
    };
    this.emitChange();
  },
  _receiveTriggerStatus: function(triggerStatus){
    this.triggerStatus = triggerStatus;
    this.emitChange();
  },
  _receiveTriggerLogs: function(triggerLogs) {
    this.triggerLogs = triggerLogs;
    this.emitChange();
  },
  _triggerByTagAttempt: function() {
    this.STATUS = STATUS.ATTEMPTING;
    this.emitChange();
  },
  _triggerByTagError: function(err) {
    this.validations.trigger.error = err;
    setTimeout(this._clearTriggerStatus.bind(this), 5000);
    this.emitChange();
  },
  _triggerByTagSuccess: function(success) {
    this.validations.trigger.success = success;
    setTimeout(this._clearTriggerStatus.bind(this), 5000);
    this.emitChange();
  },
  _clearTriggerStatus: function() {
    this.validations.trigger.success = '';
    this.validations.trigger.error = '';
    this.STATUS = STATUS.DEFAULT;
    this.emitChange();
  },
  getState: function() {
    return {
      autoBuildStore: this.autoBuildStore,
      autoBuildBlankSlate: this.autoBuildBlankSlate,
      autoBuildLinks: this.autoBuildLinks,
      autoTriggerForm: this.autoTriggerForm,
      triggerStatus: this.triggerStatus,
      triggerLinkForm: this.triggerLinkForm,
      triggerLogs: this.triggerLogs,
      validations: this.validations,
      newTags: this.newTags,
      STATUS: this.STATUS
    };
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this.autoBuildStore = state.autoBuildStore;
    this.autoBuildLinks = state.autoBuildLinks;
    this.autoTriggerForm = state.autoTriggerForm;
    this.triggerStatus = state.triggerStatus;
    this.triggerLinkForm = state.triggerLinkForm;
    this.triggerLogs = state.triggerLogs;
    this.validations = state.validations;
    this.newTags = state.newTags;
    this.STATUS = state.STATUS;
  }
});

module.exports = AutoBuildSettingsStore;
