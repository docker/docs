'use strict';

import async from 'async';
import has from 'lodash/object/has';
import { Builds } from 'hub-js-sdk';
import linkedAccountSettingsAction from './linkedAccountsSettings';
var debug = require('debug')('navigate::bitbucketRedirect');

/**
 * This action is hit when the site gets redirected from the github oauth workflow
 * @param actionContext
 * @param payload
 * @param done
 * @param maybeData
 */
export default function bitbucketRedirect({actionContext, payload, done, maybeData}){

  const SOURCES = {
    GITHUB: 'github',
    BITBUCKET: 'bitbucket'
  };

  debug('In Bitbucket Redirect Route Handler !');
  //If token is available the user is already logged in
  var token;
  if (has(maybeData, 'token')) {
    token = maybeData.token;
    var oauthVerifier = payload.location.query.oauth_verifier;
    var oauthToken = payload.location.query.oauth_token;

    var _associateBitbucketAccount = function(cb) {
      Builds.associateBitbucketAccount(token, {oauth_verifier: oauthVerifier, oauth_token: oauthToken}, function(err, res) {
        if (err) {
          debug(JSON.stringify(err));
          debug('ERROR ASSOCIATING BITBUCKET ACCOUNT: ' + err);
          actionContext.dispatch('BITBUCKET_ASSOCIATE_ERROR', res.body);
          cb(err);
        } else {
          if (res.body) {
            debug('Bitbucket account association Success: ' + res.body);
            actionContext.dispatch('BITBUCKET_ASSOCIATE_SUCCESS', res.body);
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

    async.series([
      _associateBitbucketAccount,
      _getLinkedAccounts
    ], function(err, results) {
        done();
      }
    );
  } else {
    token = '';
    //TODO: redirect to login page or test redirect happens automatically
    done();
  }
}
