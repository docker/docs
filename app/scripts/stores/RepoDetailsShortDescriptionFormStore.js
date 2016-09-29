'use strict';

import _ from 'lodash';
import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('RepoDetailsShortDescriptionFormStore');

export default createStore({
  storeName: 'RepoDetailsShortDescriptionFormStore',
  handlers: {
    RECEIVE_REPOSITORY: '_receiveRepository',
    SHORT_DESCRIPTION_ATTEMPT_START: '_shortDescriptionAttemptStart',
    SHORT_DESCRIPTION_SUCCESS: '_shortDescriptionSuccess',
    DETAILS_UNAUTHORIZED: '_detailsUnauthorized',
    DETAILS_UNAUTHORIZED_DETAIL: '_detailsUnauthorizedDetail',
    SHORT_BAD_REQUEST: '_badRequest',
    DETAILS_ERROR: '_detailsError',
    SHORT_DESCRIPTION_UPDATE_FIELD_WITH_VALUE: '_updateFieldWithValue',
    DETAILS_RESET_FORMS: '_detailsResetForms',
    TOGGLE_SHORT_DESCRIPTION_EDIT: '_toggleEditMode'
  },
  initialize: function() {
    this.isEditing = false;
    this.successfulSave = false;
    this.fields = {
      shortDescription: {}
    };

    this._defaultValues = {
      shortDescription: ''
    };

    this.values = {
      shortDescription: ''
    };
  },
  _shortDescriptionAttemptStart() {
    debug('starting short description update');
  },
  _shortDescriptionSuccess() {
    this.fields.shortDescription.success = 'Successfully updated short description.';
    //switch back to viewing mode with green outline
    this.successfulSave = true;
    this.isEditing = false;
    this._defaultValues.shortDescription = this.values.shortDescription;
    //clear the green outline and text on successful save
    setTimeout(this._clearFeedbackStates.bind(this), 5000);
    setTimeout(this._clearSuccessfulSave.bind(this), 5000);
    this.emitChange();
  },
  _receiveRepository(repo) {
    this.isEditing = false;
    this.values.shortDescription = repo.description;
    this._defaultValues.shortDescription = repo.description;
    this.emitChange();
  },
  _detailsResetForms() {
    // reset form value to repo shortdescription
    this.values.shortDescription = this._defaultValues.shortDescription;
    // reset errors
    this.fields.shortDescription = {};
    this.emitChange();
  },
  _detailsError() {
    this.fields.longDescription.hasError = true;
    this.fields.longDescription.error = 'Sorry, your long description could not be saved.';
    setTimeout(this._clearFeedbackStates.bind(this), 5000);
    this.emitChange();
  },
  _badRequest(obj) {
    this.fields.shortDescription.hasError = true;
    this.fields.shortDescription.error = obj.description[0];
    setTimeout(this._clearFeedbackStates.bind(this), 5000);
    this.emitChange();
  },
  _clearFeedbackStates() {
    this.fields.shortDescription.error = '';
    this.fields.shortDescription.hasError = false;
    this.fields.shortDescription.success = '';
    this.emitChange();
  },
  _clearSuccessfulSave() {
    this.successfulSave = false;
    this.emitChange();
  },
  _updateFieldWithValue: function({fieldKey, fieldValue}){
    this.values[fieldKey] = fieldValue;
    this.emitChange();
  },
  _toggleEditMode( { isEditing }) {
    this.isEditing = isEditing;
    //if you cancel, clear the old input
    this.values.shortDescription = this._defaultValues.shortDescription;
    //in case you recently saved and you now cancel
    this._clearSuccessfulSave();
    this.emitChange();
  },
  getState() {
    return {
      _defaultValues: this._defaultValues,
      fields: this.fields,
      values: this.values,
      isEditing: this.isEditing,
      successfulSave: this.successfulSave
    };
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this._defaultValues = state._defaultValues;
    this.fields = state.fields;
    this.values = state.values;
    this.isEditing = state.isEditing;
    this.successfulSave = state.successfulSave;
  }
});
