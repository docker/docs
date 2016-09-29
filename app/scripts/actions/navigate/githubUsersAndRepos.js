'use strict';

var debug = require('debug')('navigate::githubRepos');
import async from 'async';
import {
  Builds
  } from 'hub-js-sdk';
import has from 'lodash/object/has';
import linkedAccountAction from './linkedAccountsSettings';

export default function getGithubRepos({actionContext, payload, done, maybeData}) {

  var _getLinkedAccountStatus = function(cb) {
    linkedAccountAction({
      actionContext: actionContext,
      payload: payload,
      done: cb,
      maybeData: maybeData});
  };

  var _getSourceRepos = function(cb) {
    Builds.getSourceRepos('github', maybeData.token, function (err, res) {
      if (err) {
        const { detail } = err.response.body;
        if (detail) {
          actionContext.dispatch('LINKED_REPO_SOURCES_ERROR', detail);
        }
        cb(null);
      } else{
        cb(null, res.body);
      }
    });
  };


  if (has(maybeData, 'token')) {
    async.parallel([
      _getLinkedAccountStatus,
      _getSourceRepos
    ], function(err, results) {
      actionContext.dispatch('SET_LINKED_REPO_TYPE', 'github');
      actionContext.dispatch('RECEIVE_LINKED_REPO_SOURCES', results[1]);
      done();
    });
  } else {
    done();
  }
}
