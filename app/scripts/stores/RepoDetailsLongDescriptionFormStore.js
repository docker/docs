'use strict';

import _ from 'lodash';
import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('RepoDetailsLongDescriptionFormStore');

export default createStore({
  storeName: 'RepoDetailsLongDescriptionFormStore',
  handlers: {
    RECEIVE_REPOSITORY: '_receiveRepository',
    LONG_DESCRIPTION_ATTEMPT_START: '_longDescriptionAttemptStart',
    LONG_DESCRIPTION_SUCCESS: '_longDescriptionSuccess',
    DETAILS_UNAUTHORIZED: '_detailsUnauthorized',
    DETAILS_UNAUTHORIZED_DETAIL: '_detailsUnauthorizedDetail',
    LONG_BAD_REQUEST: '_badRequest',
    DETAILS_ERROR: '_detailsError',
    LONG_DESCRIPTION_UPDATE_FIELD_WITH_VALUE: '_updateFieldWithValue',
    DETAILS_RESET_FORMS: '_detailsResetForms',
    TOGGLE_LONG_DESCRIPTION_EDIT: '_toggleEditMode'
  },
  initialize: function() {
    this.isEditing = false;
    this.successfulSave = false;
    this.fields = {
      longDescription: {}
    };

    this._defaultValues = {
      longDescription: ''
    };

    this.values = {
      longDescription: ''
    };
  },
  _longDescriptionAttemptStart() {
    debug('starting long description update');
  },
  _longDescriptionSuccess() {
    this.fields.longDescription.success = 'Successfully updated full description.';
    //switch back to viewing mode with green outline
    this.successfulSave = true;
    this.isEditing = false;
    this._defaultValues.longDescription = this.values.longDescription;
    //clear the green outline and text on successful save
    setTimeout(this._clearFeedbackStates.bind(this), 5000);
    setTimeout(this._clearSuccessfulSave.bind(this), 5000);
    this.emitChange();
  },
  _receiveRepository(repo) {
    this.isEditing = false;
    this.values.longDescription = repo.full_description || '';
    this._defaultValues.longDescription = repo.full_description || '';
    this.emitChange();
  },
  _detailsResetForms() {
    // reset form value to repo longdescription
    this.values.longDescription = this._defaultValues.longDescription;
    // reset errors
    this.fields.longDescription = {};
    this.emitChange();
  },
  _badRequest(obj) {
    this.fields.longDescription.hasError = !!obj.full_description;
    this.fields.longDescription.error = obj.full_description[0];
    setTimeout(this._clearFeedbackStates.bind(this), 5000);
    this.emitChange();
  },
  _detailsError() {
    this.fields.longDescription.hasError = true;
    this.fields.longDescription.error = 'Sorry, your long description could not be saved.';
    setTimeout(this._clearFeedbackStates.bind(this), 5000);
    this.emitChange();
  },
  _clearFeedbackStates() {
    this.fields.longDescription.error = '';
    this.fields.longDescription.hasError = false;
    this.fields.longDescription.success = '';
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
    this.values.longDescription = this._defaultValues.longDescription;
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
