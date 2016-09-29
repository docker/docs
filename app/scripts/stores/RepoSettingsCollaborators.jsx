'use strict';

import createStore from 'fluxible/addons/createStore';
import { STATUS } from './collaborators/Constants.js';
var debug = require('debug')('RepoSettingsCollaborators');

export default createStore({
  storeName: 'RepoSettingsCollaborators',
  handlers: {
    ADD_COLLAB_START: '_addCollabStart',
    ADD_COLLAB_ERROR: '_addCollabError',
    ADD_COLLAB_SUCCESS: '_addCollabSuccess',
    COLLAB_RECEIVE_COLLABORATORS: '_receiveCollaborators',
    COLLAB_RECEIVE_TEAMS: '_receiveTeams',
    COLLAB_RECEIVE_ALL_TEAMS: '_receiveAllTeams',
    DEL_COLLABORATORS_SET_LOADING: 'setLoadingFor',
    DEL_COLLABORATORS_SET_ERROR: 'setErrorFor',
    DEL_COLLABORATORS_SET_SUCCESS: 'setSuccessFor',
    LOGOUT: 'initialize',
    ON_ADD_COLLAB_CHANGE: 'onAddCollabChange'
  },
  initialize() {
    // these are full request objects. Only one will succeed and have a `count` key
    this.collaborators = {};
    this.teams = {};
    this.allTeams = {results: []};
    this.newCollaborator = '';
    this.error = '';
    this.requests = {};
    this.STATUS = STATUS.DEFAULT;
  },
  getState() {
    return {
      collaborators: this.collaborators,
      teams: this.teams,
      allTeams: this.allTeams,
      newCollaborator: this.newCollaborator,
      error: this.error,
      requests: this.requests,
      STATUS: this.STATUS
    };
  },

  setLoadingFor(username) {
    this.requests[username] = STATUS.ATTEMPTING;
    this.emitChange();
  },
  setErrorFor(username) {
    this.requests[username] = STATUS.ERROR;
    this.emitChange();
  },
  setSuccessFor(username) {
    this.requests[username] = STATUS.DEFAULT;
    this.newCollaborator = '';
    this.emitChange();
  },
  onAddCollabChange(collaborator) {
    this.newCollaborator = collaborator;
    this.error = '';
    this.emitChange();
  },
  _addCollabSuccess() {
    this.STATUS = STATUS.SUCCESS;
    this.newCollaborator = '';
    this.error = '';
    this.emitChange();
  },
  _addCollabError(message) {
    this.STATUS = STATUS.ERROR;
    this.error = message;
    this.emitChange();
  },
  _addCollabStart() {
    this.STATUS = STATUS.ATTEMPTING;
    this.emitChange();
  },
  _receiveCollaborators(collaborators) {
    debug(collaborators);
    this.newCollaborator = '';
    this.collaborators = collaborators;
    this.emitChange();
  },
  _receiveTeams(teams) {
    debug(teams);
    this.teams = teams;
    this.emitChange();
  },
  _receiveAllTeams(teams) {
    this.allTeams = teams;
    this.emitChange();
  },
  dehydrate() {
    return this.getState();
  },
  rehydrate(state) {
    this.collaborators = state.collaborators;
    this.teams = state.teams;
    this.allTeams = state.allTeams;
    this.newCollaborator = state.newCollaborator;
    this.error = state.error;
    this.requests = state.requests;
    this.STATUS = state.STATUS;
  }
});
