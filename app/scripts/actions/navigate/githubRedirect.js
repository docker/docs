'use strict';

import _ from 'lodash';
import async from 'async';
import { Builds } from 'hub-js-sdk';
import request from 'superagent';
var debug = require('debug')('navigate::githubRedirect');
import GithubLinkStore from '../../stores/GithubLinkStore';
import linkedAccountSettingsAction from './linkedAccountsSettings';

/**
 * This action is hit when the site gets redirected from the github oauth workflow
 * @param actionContext
 * @param payload
 * @param done
 * @param maybeData
 */
export default function githubRedirect({actionContext, payload, done, maybeData}){

  const SOURCES = {
    GITHUB: 'github',
    BITBUCKET: 'bitbucket'
  };

  debug('In Gihub Redirect Route Handler !');
  //If token is available the user is already logged in
  var token;
  if (_.has(maybeData, 'token')) {
    token = maybeData.token;
    var code = payload.location.query.code;
    var state = payload.location.query.state;

    var _associateGithubAccount = function(cb) {
      Builds.associateGithubAccount(token, code, function(err, res) {
        if (err) {
          debug(JSON.stringify(err));
          debug('ERROR ASSOCIATING GITHUB ACCOUNT: ' + err);
          actionContext.dispatch('GITHUB_ASSOCIATE_ERROR', res.body);
          cb(err);
        } else {
          if (res.body) {
            debug('GitHub account association Success: ' + res.body);
            actionContext.dispatch('GITHUB_ASSOCIATE_SUCCESS', res.body);
            cb(null, res.body);
          }
        }
      });
    };

    var _getLinkedAccounts = function(callback) {
      debug('Getting linked account settings state.');
      linkedAccountSettingsAction(
        {
          actionContext: actionContext,
          payload: payload,
          done: callback,
          maybeData: maybeData
        }
      );
    };

    debug('Payload Cookies: ' + payload.cookies);
    //TODO: validate state before linking
    if (payload.cookies && payload.cookies.ghOauthKey === state) {
      async.series([
        _associateGithubAccount,
        _getLinkedAccounts
      ], function(err, results) {
          request.post('/oauth/github-done/')
            .end(function(e, r) { debug('github-oauth done or exited.'); });
          done();
        }
      );
    } else {
      //The validation of the state failed
      actionContext.dispatch('GITHUB_SECURITY_ERROR', 'There was a security issue with your request. Please try again later.');
      done();
    }
  } else {
    token = '';
    actionContext.dispatch('GITHUB_ASSOCIATE_ERROR');
    done();
  }
}
