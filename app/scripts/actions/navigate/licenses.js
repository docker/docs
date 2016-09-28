'use strict';
var debug = require('debug')('navigate::licenses');
import async from 'async';
import request from 'superagent';
import _ from 'lodash';
import {
  Repositories as Repos,
  Notifications,
  Orgs,
  Users
} from 'hub-js-sdk';

function getLicenses({ userID, token, actionContext}, done) {
  request(process.env.REGISTRY_API_BASE_URL + '/api/licensing/v3/license/' + userID + '/')
  .accept('application/json')
  .set('Authorization', 'JWT ' + token)
  .end((err, res) => {
    if (err) {
      done(null, null);
    } else {

      let licenses = _.map(res.body.licenses, function(obj) {
        return _.merge({},
                       obj,
                       {
                         orgname: userID
                       });
      });
      done(null, licenses);
    }
    });
}

export default function licensesFn({
  actionContext, payload, done, maybeData
}){

  if(maybeData.token) {
    Users.getNamespacesForUser(maybeData.token, function(err, res){
      async.map(res.body.namespaces,
                function(item, cb) {
                  getLicenses({
                    userID: item,
                    token: maybeData.token,
                    actionContext
                  }, cb);
                },
                function(error, results) {
                  actionContext.dispatch('RECEIVE_LICENSES', _.compact(results));
                  done();
                });
    });
  } else {
    // user must be logged in; they aren't
    done();
  }
}
