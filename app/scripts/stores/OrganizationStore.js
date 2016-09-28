'use strict';

import createStore from 'fluxible/addons/createStore';
import _ from 'lodash';

var OrganizationStore = createStore({
  storeName: 'OrganizationStore',
  handlers: {
    RECEIVE_ORGANIZATION: '_updateOrg',
    CREATED_ORGANIZATION: '_onCreateOrg',
    SELECT_ORGANIZATION: '_onCurrentOrgChange',
    RECEIVE_ORG_TEAMS: '_receiveOrgTeams',
    CURRENT_USER_ORGS: '_onGetCurrentOrgs',
    UPDATE_ORG_SUCCESS: '_updateOrgSuccess',
    UPDATE_ORG_ERROR: '_updateOrgError'
  },
  initialize: function() {
    // initialize
    this.name = '';
    this.gravatarURL = 'https://secure.gravatar.com/avatar/00000000000000000000000000000000?d=retro&f=y';
    this.currentOrg = {};
    this.orgs = [];
    this.currentOrgTeams = [];
    this.success = '';
    this.error = '';
  },
  _onCreateOrg: function(payload) {
    this.receiveState({
      currentOrg: payload.newOrg,
      orgs: payload.userOrgs
    });
    this.emitChange();
  },
  _updateOrg: function(payload) {
    this.receiveState({
      currentOrg: payload
    });
    this.emitChange();
  },
  _onGetCurrentOrgs: function(payload) {
    this.receiveState({
      orgs: payload.results
    });
    this.emitChange();
  },
  _onCurrentOrgChange: function(payload) {
    this.receiveState({
      currentOrg: this.getOrg(payload)
    });
    this.emitChange();
  },
  _receiveOrgTeams: function(orgTeams) {
    this.receiveState({
      currentOrgTeams: orgTeams.results
    });
    this.emitChange();
  },
  _updateOrgSuccess: function() {
    this.success = 'Updated Organization Details Successfully!';
    setTimeout(this._clearOrgErrorStates.bind(this), 5000);
    this.emitChange();
  },
  _updateOrgError: function(err) {
    var errResponse = err.response;
    if (errResponse.badRequest) {
      _.forIn(errResponse.body, function(v, k) {
        this.error += k + ': ' + v.join(',') + '\n';
      }.bind(this));
    } else if(errResponse.unauthorized || errResponse.forbidden) {
      this.error = 'You have no permission to edit this organization.';
    } else {
      this.error = 'An error occurred during the organization update. Please try again later';
    }
    setTimeout(this._clearOrgErrorStates.bind(this), 5000);
    this.emitChange();
  },
  _clearOrgErrorStates: function() {
    this.success = '';
    this.error = '';
    this.emitChange();
  },
  receiveState: function(payload) {
    this.name = payload.orgname || this.name;
    this.gravatarURL = payload.gravatar_url || this.gravatarURL;
    this.currentOrg = payload.currentOrg || this.currentOrg;
    this.currentOrgTeams = payload.currentOrgTeams || this.currentOrgTeams;
    this.orgs = payload.orgs || this.orgs;
  },
  getState: function() {
    return {
      name: this.name,
      gravatarURL: this.gravatarURL,
      currentOrg: this.currentOrg,
      currentOrgTeams: this.currentOrgTeams,
      orgs: this.orgs,
      error: this.error,
      success: this.success
    };
  },
  getOrgs: function() {
    return this.orgs;
  },
  getCurrentOrg: function() {
    //returns currently selected org
    return this.currentOrg;
  },
  getOrg: function(name) {
    //Assuming org names are unique and expecting filter to return an array of exactly 1 item
    return _.filter(this.orgs, function(org) {
      return org.orgname === name;
    })[0];
  },
  getOrgTeams: function() {
    return this.currentOrgTeams;
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this.name = state.name;
    this.gravatarURL = state.gravatarURL;
    this.currentOrg = state.currentOrg;
    this.currentOrgTeams = state.currentOrgTeams;
    this.orgs = state.orgs;
    this.error = state.error;
    this.success = state.success;
  }
});

module.exports = OrganizationStore;
