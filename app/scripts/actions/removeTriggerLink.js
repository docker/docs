'use strict';
import _ from 'lodash';
import { Autobuilds as AutoBuild } from 'hub-js-sdk';
var debug = require('debug')('hub:actions:removeTriggerLink');

export default function(actionContext, { JWT, namespace, name, repo_id }) {
  AutoBuild.deleteAutomatedBuildLink(JWT, namespace, name, repo_id, function(err, res) {
    if (err) {
      debug('error', err);
    } else {
      AutoBuild.getAutomatedBuildLinks(JWT, namespace, name, function(getErr, getRes){
        if (getErr) {
          debug('getAutomatedBuildLinks error', err);
        } else {
          actionContext.dispatch('RECEIVE_AUTOBUILD_LINKS', getRes.body.results);
        }
      });
    }
  });
}
