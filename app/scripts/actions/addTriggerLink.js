'use strict';
import _ from 'lodash';
import { Autobuilds as AutoBuild } from 'hub-js-sdk';
var debug = require('debug')('hub:actions:addTriggerLink');

export default function(actionContext, { JWT, namespace, name, to_repo }) {
  // to_repo needs to be a repo id
  AutoBuild.addAutomatedBuildLink(JWT, namespace, name, to_repo, function(err, res) {
    if (err) {
      debug(err);
      actionContext.dispatch('LINK_AUTOBUILD_ERROR');
    } else {
      actionContext.dispatch('LINK_AUTOBUILD_SUCCESS');
      AutoBuild.getAutomatedBuildLinks(JWT, namespace, name, function(getErr, getRes){
        if (getErr) {
          debug(err);
        } else {
          actionContext.dispatch('RECEIVE_AUTOBUILD_LINKS', getRes.body.results);
        }
      });
    }
  });
}
