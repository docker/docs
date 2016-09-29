'use strict';

var debug = require('debug')('navigate::linkedAccounts');
import async from 'async';
import {
  Builds
  } from 'hub-js-sdk';
import _ from 'lodash';

const SOURCES = {
  GITHUB: 'github',
  BITBUCKET: 'bitbucket'
};

export default function linked({actionContext, payload, done, maybeData}){
  if (_.has(maybeData, 'token')) {

    var _checkGithub = function(cb) {
      Builds.checkGithub(maybeData.token, function(err, res) {
        if(err) {
          debug(err);
          cb(null, false);
        } else {
          debug('Github account exists');
          cb(null, true);
        }
      });
    };

    var _checkBitbucket = function(cb) {
      Builds.checkBitbucket(maybeData.token, function(err, res) {
        if(err) {
          debug(err);
          cb(null, false);
        } else {
          cb(null, true);
        }
      });
    };

    var _getGithubAccount = function(accountExists, cb) {
      if(accountExists) {
        Builds.getSourceAccount(SOURCES.GITHUB, maybeData.token, function (err, res) {
          if (err) {
            cb(null, err);
          } else {
            cb(null, res.body);
          }
        });
      } else {
        cb(null, {detail: 'No associated Github user'});
      }
    };

    var _getBitbucketAccount = function(accountExists, cb) {
      if (accountExists) {
        Builds.getSourceAccount(SOURCES.BITBUCKET, maybeData.token, function (err, res) {
          if(err) {
            cb(null, err);
          } else {
            cb(null, res.body);
          }
        });
      } else {
        Builds.getBitbucketAuthUrl(maybeData.token, function(err, res) {
          if (err) {
            debug(err);
            actionContext.dispatch('BITBUCKET_AUTH_URL_ERROR', err);
          } else if (res.ok) {
            actionContext.dispatch('RECEIVE_BITBUCKET_AUTH_URL', res.body);
          }
        });
        cb(null, {detail: 'No associated Bitbucket user'});
      }
    };

    var _wfGetGithubAccount = function(cb) {
      async.waterfall([
        _checkGithub,
        _getGithubAccount
      ], function(err, result) {
        debug('Github WF - result', result);
        if (!err) {
          cb(null, result);
        }
      });
    };

    var _wfGetBitbucketAccount = function(cb) {
      async.waterfall([
        _checkBitbucket,
        _getBitbucketAccount
      ], function(err, result) {
        debug('Bitbucket WF - result', result);
        if (!err) {
          cb(null, result);
        }
      });
    };

    //Check github and bitbucket and then get the account
    async.parallel({
      _wfGetGithubAccount,
      _wfGetBitbucketAccount
    }, function(err, results) {
      debug('results: ', results);
      let accounts = {
        github: results._wfGetGithubAccount,
        bitbucket: results._wfGetBitbucketAccount,
        gitlab: {detail: 'No associated GitLab user'}
      };
      actionContext.dispatch('RECEIVE_SOURCE_ACCOUNTS', accounts);
      return done();
    });
  } else {
    done();
  }
}
