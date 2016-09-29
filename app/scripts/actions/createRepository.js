/* @flow */
'use strict';

var createRepo = require('hub-js-sdk').Repositories.createRepository;
var debug: Function = require('debug')('hub:actions:createRepository');

var createRepository = function(actionContext: {dispatch: Function},
                          {jwt, repository}: {jwt: string; repository: any}) {
  createRepo(jwt, repository, (err, res) => {
    if (err) {

      if(res && res.badRequest) {
        actionContext.dispatch('CREATE_REPO_BAD_REQUEST', res.body);
      } else {
        actionContext.dispatch('CREATE_REPO_ERROR', err);
      }

    } else {

      if (res && res.ok) {
        actionContext.dispatch('RECEIVE_REPOSITORY', res.body);
        actionContext.dispatch('CREATE_REPO_CLEAR_FORM');
        actionContext.history.push(`/r/${repository.namespace}/${repository.name}/`);
      }

    }
  });
};

module.exports = createRepository;
