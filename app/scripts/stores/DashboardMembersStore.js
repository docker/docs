'use strict';

import createStore from 'fluxible/addons/createStore';
import _ from 'lodash';
var debug = require('debug')('DashboardMembersStore');
import { STATUS } from './orgteamstore/Constants';

var DashboardMembersStore = createStore({
  storeName: 'DashboardMembersStore',
  handlers: {
    RECEIVE_DASHBOARD_TEAM_MEMBERS: '_receiveDashboardTeamMembers',
    ORG_DASHBOARD_MEMBERS_ERROR: '_errorReceivingMembers',
    TEAM_MEMBER_ERROR: '_teamMemberError',
    TEAM_MEMBER_BAD_REQUEST: '_teamMemberBadRequest',
    TEAM_MEMBER_UNAUTHORIZED: '_teamMemberUnauthorized',
    CLEAR_MEMBER_ERROR: '_clearErrorStates'
  },
  initialize() {
    this.members = [];
    this.count = 0;
    this.error = {};
    this.STATUS = STATUS.DEFAULT;
  },
  _errorReceivingMembers(err) {
    debug(err);
  },
  _receiveDashboardTeamMembers(members) {
    debug(members);
    this.members = members;
    this.count = members.length;
    this.emitChange();
  },
  _teamMemberError: function(err) {
    this.STATUS = STATUS.MEMBER_ERROR;
    this.STATUS = STATUS.GENERAL_SERVER_ERROR;
    this.error = err;
    this.emitChange();
  },
  _teamMemberBadRequest: function(err) {
    this.STATUS = STATUS.MEMBER_BAD_REQUEST;
    this.error = err;
    this.emitChange();
  },
  _teamMemberUnauthorized: function(err) {
    this.STATUS = STATUS.MEMBER_UNAUTHORIZED;
    this.error = err;
    this.emitChange();
  },
  _clearErrorStates: function() {
    this.error = {};
    this.STATUS = STATUS.DEFAULT;
    this.emitChange();
  },
  getState() {
    return {
      members: this.members,
      count: this.count,
      error: this.error,
      STATUS: this.STATUS
    };
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this.members = state.members;
    this.count = state.count;
    this.STATUS = state.STATUS;
    this.error = state.error;
  }
});

module.exports = DashboardMembersStore;
