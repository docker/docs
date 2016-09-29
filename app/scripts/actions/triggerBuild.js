'use strict';

import { Autobuilds } from 'hub-js-sdk';
var debug = require('debug')('hub:actions:triggerBuild');

module.exports = function(actionContext, {JWT, name, namespace}) {
  Autobuilds.triggerBuild(JWT, namespace, name, (err, res) => {
    if (err) {
      debug(err);
      actionContext.dispatch('AB_TRIGGER_ERROR');
    } else {
      actionContext.dispatch('AB_TRIGGER_SUCCESS');
    }
  });
};
