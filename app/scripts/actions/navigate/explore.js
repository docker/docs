'use strict';
var debug = require('debug')('navigate::home');
import async from 'async';
import _ from 'lodash';
import {
  Repositories as Repos,
  Notifications,
  Orgs,
  Users
} from 'hub-js-sdk';

export default function explore({actionContext, payload, done, maybeData}){

  //Get repos for library
  // We have, hijacked the repos store. This might not be good
  // in the long run
  Repos.getReposForUser(maybeData.token, 'library', function(err, res) {
    if (err) {
      actionContext.dispatch('ERROR_RECEIVING_REPOS');
      done();
    } else {
      actionContext.dispatch('RECEIVE_REPOS', res.body);
      done();
    }
  }, payload.location.query.page);
}
