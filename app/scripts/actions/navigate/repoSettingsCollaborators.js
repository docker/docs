'use strict';
const debug = require('debug')('navigate::repo');
import { parallel } from 'async';
import { PENDING_DELETE } from 'common/enums/RepoStatus';
import _ from 'lodash';
import {
  Repositories as Repos,
  Orgs
} from 'hub-js-sdk';

function getRepo({maybeToken, actionContext, user, splat}) {
  return function(callback){
    Repos.getRepo(maybeToken, `${user}/${splat}`, function(err, res) {
      let status;
      if (res && res.body) {
        status = res.body.status;
      }

      if (err || status === PENDING_DELETE) {
        actionContext.dispatch('REPO_NOT_FOUND', null);
        return callback(err);
      }
      actionContext.dispatch('RECEIVE_REPOSITORY', res.body);
      return callback();
     });
};
}

// GET's the collaborators for a user repo (which should be individuals)
function getCollaborators({maybeToken, actionContext, user, splat}) {
  return function(cb) {
    Repos.getCollaboratorsForRepo(maybeToken, `${user}/${splat}`, (err, res) => {
                                  if(err) {
                                    // 'Org repositories do not have collaborators.'
                                    actionContext.dispatch('COLLAB_RECEIVE_COLLABORATORS', {});
                                  } else {
                                    actionContext.dispatch('COLLAB_RECEIVE_COLLABORATORS', res.body);
                                  }
                                  cb();
                                 });
  };
}

 //GET's the collaborators for an organization repo (which should be teams)
function getTeamCollaborators({maybeToken, actionContext, user, splat}) {
  return function(cb) {
    Repos.getTeamCollaboratorsForRepo(maybeToken, `${user}/${splat}`, (err, res) => {
                                      if (err) {
                                        // 'User repository does not have any teams yet'
                                        actionContext.dispatch('COLLAB_RECEIVE_TEAMS', {});
                                      } else {
                                        actionContext.dispatch('COLLAB_RECEIVE_TEAMS', res.body);
                                      }
                                      cb();
                                     });
  };
}

// GET's all teams for the organization
function getOrgTeams({maybeToken, actionContext, user, splat}) {
  return function(cb) {
    Orgs.getTeams(maybeToken, user, (err, res) => {
      if (err) {
        // 'No such organization'
        actionContext.dispatch('COLLAB_RECEIVE_TEAMS', {});
        actionContext.dispatch('COLLAB_RECEIVE_ALL_TEAMS', {results: []});
      } else {
        actionContext.dispatch('COLLAB_RECEIVE_ALL_TEAMS', res.body);
      }
      cb();
    });
  };
}

export default function repoSettingsMain({actionContext, payload, done, maybeData}){
  debug('maybeData:', maybeData);
  if (_.has(maybeData, 'token')) {
    let args = {
      actionContext,
      maybeToken: maybeData.token,
      user: payload.params.user,
      splat: payload.params.splat
    };

    parallel([
      getRepo(args),
      getCollaborators(args),
      getTeamCollaborators(args),
      getOrgTeams(args)
    ], function(err, res){
      done();
    });
  } else {
    actionContext.dispatch('REPO_NOT_FOUND', null);
    done();
  }
}
