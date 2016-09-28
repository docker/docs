'use strict';
const debug = require('debug')('navigate::repoBuildSettings');
import { parallel, waterfall } from 'async';
import { PENDING_DELETE } from 'common/enums/RepoStatus';
import _ from 'lodash';
import {
  Repositories as Repos,
  Autobuilds as AutoBuild
  } from 'hub-js-sdk';

function getRepo({maybeToken, actionContext, user, splat}) {
  return function(callback){
    Repos.getRepo(maybeToken, `${user}/${splat}`, function(err, res) {
      let status;
      if (res && res.body) {
        status = res.body.status;
      }

      if (err || status === PENDING_DELETE) {
        debug('GET REPO ERR::', err);
        actionContext.dispatch('REPO_NOT_FOUND', err);
        return callback(err);
      } else {
        debug('GETTING REPOSITORY::', res.body);
        actionContext.dispatch('RECEIVE_REPOSITORY', res.body);
        return callback(null, res.body);
      }
    });
  };
}

function getAutobuildSettings({maybeToken, actionContext, user, splat}) {
  return function(callback){
    AutoBuild.getAutomatedBuildSettings(maybeToken, user, splat, function(err, res) {
      if (err) {
        debug('GET AUTOBUILD SETTINGS ERR::', err);
        return callback(null, null);
      }
      debug('AUTOBUILDSETTINGS', res.body);
      actionContext.dispatch('RECEIVE_AUTOBUILD_SETTINGS', res.body);
      //also dispatch to the triggerByTag store to initialize the status
      const {build_tags} = res.body;
      if (build_tags) {
        actionContext.dispatch('INITIALIZE_AB_TRIGGERS', build_tags);
      }
      return callback(null, res.body);
    });
  };
}

function getAutobuildLinks({maybeToken, actionContext, user, splat}) {
  return function(callback) {
    AutoBuild.getAutomatedBuildLinks(maybeToken, user, splat, function(err, res) {
      if (err) {
        debug('GET AUTOBUILD LINKS ERR::', err);
        return callback(null, null);
      }
      debug('AUTOBUILDLINKS', res.body);
      actionContext.dispatch('RECEIVE_AUTOBUILD_LINKS', res.body.results);
      return callback(null, res.body);
    });
  };
}

function getTriggerStatus({maybeToken, actionContext, user, splat}) {
  return function(callback) {
    AutoBuild.getTriggerStatus(maybeToken, user, splat, function(err, res) {
      if (err) {
        debug('GET TRIGGER STATUS ERR::', err);
        const defaultStatus = {
          token: '',
          trigger_url: '',
          active: false
        };
        actionContext.dispatch('RECEIVE_TRIGGER_STATUS', defaultStatus);
        return callback(null, null);
      }
      debug('AUTOBUILDTRIGGERS', res.body);
      actionContext.dispatch('RECEIVE_TRIGGER_STATUS', res.body);
      return callback(null, res.body);
    });
  };
}

function getTriggerLogs({maybeToken, actionContext, user, splat}) {
  return function(callback) {
    AutoBuild.getTriggerLogs(maybeToken, user, splat, function(err, res) {
      if (err) {
        debug('GET TRIGGER LOGS ERROR', err.res);
        return callback(null, null);
      }
      debug('AUTOBUILDTRIGGERLOGS', res.body);
      actionContext.dispatch('RECEIVE_TRIGGER_LOGS', res.body.results);
      return callback(null, res.body.results);
    });
  };
}

function getRest(args) {
  return function(repoDetails, cbk) {
    if (repoDetails) {
      parallel({
        getAutoBuild: getAutobuildSettings(args),
        getAutoLinks: getAutobuildLinks(args),
        getTriggerStatus: getTriggerStatus(args),
        getTriggerLogs: getTriggerLogs(args)
      }, function (err, res) {
        cbk(null, res);
      });
    } else {
      cbk(null);
    }
  };
}

export default function repoSettingsBuilds({actionContext, payload, done, maybeData}) {
  debug('maybeData:', maybeData);
  if (_.has(maybeData, 'token')) {
    let args = {
      actionContext,
      maybeToken: maybeData.token,
      user: payload.params.user,
      splat: payload.params.splat
    };

    waterfall([
      getRepo(args),
      getRest(args)
    ], function (e, r) {
      done();
    });
  } else {
    actionContext.dispatch('REPO_NOT_FOUND', null);
    done();
  }
}
