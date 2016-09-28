'use strict';
import createStore from 'fluxible/addons/createStore';
import _ from 'lodash';
const STATUS = require('./repostore/Constants').STATUS;
var debug = require('debug')('RepositoryPageStore');

var RepoStore = createStore({
  storeName: 'RepositoryPageStore',
  handlers: {
    RECEIVE_REPOSITORY: '_receiveRepository',
    CREATE_REPO_ERROR: '_createRepoError',
    TOGGLE_STARRED_STATE: '_toggleStarred',
    TOGGLE_VISIBILITY_SUCCESS: '_toggleVisibility',
    REPO_NOT_FOUND: '_repoNotFound'
  },
  initialize: function() {
    this.canEdit = false;
    this.description = '';
    this.fullDescription = '';
    this.hasStarred = false;
    this.isPrivate = true;
    this.isAutomated = false;
    this.name = '';
    this.namespace = '';
    this.status = 0;
    this.lastUpdated = '';
    this.globalFormError = '';
    this.STATUS = STATUS.DEFAULT;
  },
  _createRepoError: function(err) {
    if (err) {
      var errResponse = err.response.body;
      this.globalFormError = '';
      if (!_.isEmpty(errResponse)) {
        if (err.response.badRequest) {
          this.STATUS = STATUS.BAD_REQUEST;
          if (_.has(errResponse, '__all__')) {
            this.globalFormError = errResponse.__all__.toString();
          } else if (_.has(errResponse, 'detail')) {
            this.globalFormError = errResponse.detail.toString();
          } else {
            _.forIn(errResponse, function(v, k) {
              this.globalFormError += k + ': ' + v.join(',') + '\n';
            }.bind(this));
          }
        }
      } else {
        this.STATUS = STATUS.GENERAL_SERVER_ERROR;
        this.globalFormError = 'An error occurred while creating your repository. Please try again later.';
      }
    }
    this.emitChange();
  },
  _receiveRepository: function(res) {
    debug('receive repo', res);
    this.STATUS = STATUS.DEFAULT;
    this.canEdit = res.can_edit;
    this.description = res.description;
    // full_description can come in as null; Default to string
    this.fullDescription = res.full_description || '';
    this.hasStarred = res.has_starred;
    this.isPrivate = res.is_private;
    this.isAutomated = res.is_automated;
    this.lastUpdated = res.last_updated;
    this.name = res.name;
    this.namespace = res.namespace;
    this.status = res.status;

    this.emitChange();
  },
  _toggleStarred: function(status) {
    this.hasStarred = status;
    this.emitChange();
  },
  _toggleVisibility: function(vis) {
    this.isPrivate = vis;
    this.emitChange();
  },
  _repoNotFound: function(err) {
    this.STATUS = STATUS.REPO_NOT_FOUND;
    this.emitChange();
  },
  getState: function() {
    return {
      canEdit: this.canEdit,
      description: this.description,
      fullDescription: this.fullDescription,
      hasStarred: this.hasStarred,
      isPrivate: this.isPrivate,
      isAutomated: this.isAutomated,
      lastUpdated: this.lastUpdated,
      name: this.name,
      namespace: this.namespace,
      status: this.status,
      globalFormError: this.globalFormError,
      STATUS: this.STATUS
    };
  },
  rehydrate: function(state) {
    this.canEdit = state.canEdit;
    this.description = state.description;
    this.fullDescription = state.fullDescription;
    this.hasStarred = state.hasStarred;
    this.isPrivate = state.isPrivate;
    this.isAutomated = state.isAutomated;
    this.name = state.name;
    this.lastUpdated = state.lastUpdated;
    this.namespace = state.namespace;
    this.status = state.status;
    this.globalFormError = state.globalFormError;
    this.STATUS = state.STATUS;
  },
  dehydrate: function() {
    return this.getState();
  }
});

module.exports = RepoStore;
