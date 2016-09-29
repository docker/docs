'use strict';

import createStore from 'fluxible/addons/createStore';
import _ from 'lodash';
import { STATUS } from './orgteamstore/Constants';
var debug = require('debug')('DashboardTeamsStore');

var DashboardTeamsStore = createStore({
  storeName: 'DashboardTeamsStore',
  handlers: {
    RECEIVE_DASHBOARD_ORG_TEAMS: '_receiveDashboardOrgTeams',
    TEAM_ERROR: '_orgTeamError',
    TEAM_BAD_REQUEST: '_teamBadRequest',
    TEAM_UNAUTHORIZED: '_teamUnauthorized',
    UPDATE_TEAM_ERROR: '_updateTeamError',
    UPDATE_TEAM_SUCCESS: '_updateTeamSuccess',
    TEAM_READ_ONLY: '_isTeamReadOnly'
  },
  initialize() {
    this.teams = [];
    this.count = 0;
    this.teamReadOnly = false;
    this.errorDetails = {detail: ''};
    this.success = '';
    this.STATUS = STATUS.DEFAULT;
  },
  _receiveDashboardOrgTeams(orgTeams) {
    debug(orgTeams);
    this.teams = _.sortBy(orgTeams.results, 'name');
    this.count = orgTeams.count;
    this.emitChange();
  },
  _orgTeamError: function(err) {
    this.STATUS = STATUS.TEAM_ERROR;
    this.STATUS = STATUS.GENERAL_SERVER_ERROR;
    this.errorDetails = {detail: 'Username does not exist or it is invalid.'};
    this.emitChange();
  },
  _teamBadRequest: function(err) {
    this.STATUS = STATUS.TEAM_BAD_REQUEST;
    this.errorDetails = err;
    this.emitChange();
  },
  _teamUnauthorized: function(err) {
    this.STATUS = STATUS.TEAM_UNAUTHORIZED;
    this.errorDetails = err;
    this.emitChange();
  },
  _isTeamReadOnly: function(flag) {
    this.teamReadOnly = flag;
    this.emitChange();
  },
  _updateTeamError: function(err) {
    this.STATUS = STATUS.UPDATE_TEAM_ERROR;
    if (err.response) {
      var errResp = err.response;
      if (errResp.badRequest) {
        this.errorDetails = err;
      } else if (errResp.unauthorized || errResp.forbidden) {
        this.errorDetails = {detail: 'You are not permitted to edit this team.'};
      } else {
        this.errorDetails = {detail: 'Error updating team. Check if name is between 3 and 30 characters with no spaces.'};
      }
    }
    this.emitChange();
  },
  _updateTeamSuccess: function(err) {
    this.STATUS = STATUS.UPDATE_TEAM_SUCCESS;
    this.success = 'Team successfully updated.';
    setTimeout(this._clearErrorStates.bind(this), 5000);
    this.emitChange();
  },
  _clearErrorStates: function() {
    this.errorDetails = {detail: ''};
    this.success = '';
    this.STATUS = STATUS.DEFAULT;
    this.emitChange();
  },
  getState() {
    return {
      teams: this.teams,
      count: this.count,
      teamReadOnly: this.teamReadOnly,
      success: this.success,
      errorDetails: this.errorDetails
    };
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this.teams = state.teams;
    this.teamReadOnly = state.teamReadOnly;
    this.count = state.count;
    this.success = state.success;
    this.errorDetails = state.errorDetails;
  }
});

module.exports = DashboardTeamsStore;
