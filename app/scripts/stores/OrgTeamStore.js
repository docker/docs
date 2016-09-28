'use strict';

import createStore from 'fluxible/addons/createStore';
import { STATUS } from './orgteamstore/Constants';

var OrgTeamStore = createStore({
  storeName: 'OrgTeamStore',
  handlers: {
    CREATE_ORG_TEAM: '_createOrgTeam',
    RECEIVE_ORG_TEAM: '_receiveOrgTeam',
    RECEIVE_TEAM_MEMBERS: '_receiveOrgMembers',
    TEAM_ERROR: '_orgTeamError',
    TEAM_BAD_REQUEST: '_teamBadRequest',
    TEAM_UNAUTHORIZED: '_teamUnauthorized',
    ORG_TEAM_CLEAR_ERROR_STATES: '_clearErrorStates'
  },
  initialize: function() {
    // initialize
    this.name = '';
    this.description = '';
    this.members = [];
    this.errorDetails = {};
    this.success = '';
    this.STATUS = STATUS.DEFAULT;
  },
  //TODO: this will be removed once we have API
  _createOrgTeam: function(payload) {
    this.name = payload.name;
    this.description = payload.description;
    this.STATUS = STATUS.CREATE_TEAM_SUCCESS;
    this.emitChange();
  },
  _receiveOrgTeam: function(orgTeam) {
    this.name = orgTeam.name;
    this.description = orgTeam.description;
    this.emitChange();
  },
  _receiveOrgMembers: function(members) {
    this.members = members;
    this.emitChange();
  },
  _orgTeamError: function(err) {
    this.STATUS = STATUS.TEAM_ERROR;
    this.STATUS = STATUS.GENERAL_SERVER_ERROR;
    this.errorDetails = {detail: 'Error updating team. Check if name is between 3 and 30 characters with no spaces.'};
    this.emitChange();
  },
  _teamBadRequest: function(err) {
    this.STATUS = STATUS.TEAM_BAD_REQUEST;
    this.errorDetails = {detail: 'Please check your input values. The team name may already exist or the characters may be invalid.'};
    this.emitChange();
  },
  _teamUnauthorized: function(err) {
    this.STATUS = STATUS.TEAM_UNAUTHORIZED;
    this.errorDetails = {detail: 'You have no permission to edit this team.'};
    this.emitChange();
  },
  _clearErrorStates: function() {
    this.errorDetails = {};
    this.success = '';
    this.STATUS = STATUS.DEFAULT;
    this.emitChange();
  },
  getState: function() {
    return {
      name: this.name,
      description: this.description,
      members: this.members,
      errorDetails: this.errorDetails,
      success: this.success,
      STATUS: this.STATUS
    };
  },
  getMembers: function() {
    return this.members;
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this.name = state.name;
    this.description = state.description;
    this.members = state.members;
    this.errorDetails = state.errorDetails;
    this.success = state.success;
    this.STATUS = state.STATUS;
  }
});

module.exports = OrgTeamStore;
