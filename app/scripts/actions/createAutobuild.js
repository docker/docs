'use strict';
import { Autobuilds } from 'hub-js-sdk';
const getRepository = require('hub-js-sdk').Repositories.getRepo;
import async from 'async';
import has from 'lodash/object/has';
import omit from 'lodash/object/omit';
const debug = require('debug')('hub:actions:createAutobuild');

/**
 *
 * @param actionContext
 * @param jwt
 * @param autobuildConfig
 *        {name, namespace, description, active, is_automated,
 *         provider, is_private, dockerfileLocation, sourceName, sourceType}
 */
export default function(actionContext, {JWT, autobuildConfig}) {
  /**
   * CreateAutoBuildSerializer {
   *   vcs_repo_name (string),
   *   provider (choice) = ['github' or 'bitbucket'],
   *   dockerhub_repo_name (string),
   *   is_private (boolean),
   *   build_tags (array[string])
   *   description?
   * }
   * @param cb
   * @private
   */
  var _createAutobuild = function(cb) {
    var bTags = autobuildConfig.tags;
    for (var i = 0; i < bTags.length; ++i) {
      bTags[i] = omit(bTags[i], 'id');
    }
    var ab = {
      name: autobuildConfig.name,
      namespace: autobuildConfig.namespace,
      description: autobuildConfig.description,
      vcs_repo_name: autobuildConfig.build_name,
      provider: autobuildConfig.provider,
      dockerhub_repo_name: autobuildConfig.namespace + '/' + autobuildConfig.name,
      is_private: autobuildConfig.is_private,
      active: autobuildConfig.active,
      build_tags: bTags
    };
    Autobuilds.createAutomatedBuild(JWT, ab, function(err, res) {
      if (err) {
        debug('createAutomatedBuild error', err);
        if (err.response.badRequest) {
          //Check fields and set a better response for fields
          const { detail } = err.response.body;
          if(detail) {
            actionContext.dispatch('AUTOBUILD_BAD_REQUEST', detail);
          }
        } else if (err.response.unauthorized) {
          actionContext.dispatch('AUTOBUILD_UNAUTHORIZED', err);
        } else if (err.response.serverError) {
          actionContext.dispatch('AUTOBUILD_ERROR', err);
        }
        cb(err);
      } else if (res.body) {
        var repoUrl = res.body.docker_url;
        actionContext.dispatch('AUTOBUILD_SUCCESS');
        cb(null, repoUrl);
      }
    });
  };

  /**
   *
   * @param repoUrl something like arunan/d3
   * @param cb
   * @private
   */
  var _getRepository = function(repoUrl, cb) {
    getRepository(JWT, repoUrl, function(err, res) {
      if (err) {
        debug('getRepository error', err);
        actionContext.dispatch('GET_REPOSITORY_ERROR');
        cb(err);
      } else if (res.body) {
        cb(null, res.body);
      }
    });
  };

  actionContext.dispatch('ATTEMPTING_AUTOBUILD_CREATION');
  async.waterfall([
    _createAutobuild,
    _getRepository
  ], function(err, result) {
    if (!err) {
      actionContext.dispatch('RECEIVE_REPOSITORY', result);
    }
  });

}
