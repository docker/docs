'use strict';

const debug = require('debug')('hub:actions:githubOauth');
const _ = require('lodash');
const request = require('superagent');

module.exports = function(actionContext, {stateString}) {
  request.post('/oauth/github-attempt/')
    .send( {ghk: stateString} )
    .end(function(err, res) {
    });
};
