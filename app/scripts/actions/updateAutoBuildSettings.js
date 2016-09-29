'use strict';

import { Autobuilds } from 'hub-js-sdk';
var debug = require('debug')('hub:actions:updateAutoBuildSettings');

module.exports = function(actionContext, {JWT, namespace, name, data}) {
  Autobuilds.updateAutomatedBuildSettings(JWT, namespace, name, data, function(err, res){
    if (err) {
      debug(res.body);
    } else if (res.ok) {
      Autobuilds.getAutomatedBuildSettings(JWT, namespace, name, function(getErr, getRes) {
        if (getErr) {
          debug(getErr);
        } else if (getRes.ok) {
          actionContext.dispatch('RECEIVE_AUTOBUILD_SETTINGS', getRes.body);
        }
      });
    }
  });
};
