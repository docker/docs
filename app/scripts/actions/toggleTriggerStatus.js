'use strict';
const debug = require('debug')('hub:actions:toggleTriggerStatus');
import {
  Autobuilds as AutoBuild
  } from 'hub-js-sdk';

function getTriggerStatus(JWT, actionContext, namespace, name) {
  AutoBuild.getTriggerStatus(JWT, namespace, name, function(err, res) {
    if (err) {
      debug(err);
    }
    if (res.ok) {
      actionContext.dispatch('RECEIVE_TRIGGER_STATUS', res.body);
    }
  });
}

export default function toggleTriggerStatus(actionContext, {JWT, namespace, name, active}) {
  AutoBuild.toggleTriggerStatus(JWT, namespace, name, active, function(err, res) {
    if (err) {
      debug(err);
    }
    if (res.ok) {
      getTriggerStatus(JWT, actionContext, namespace, name);
    }
  });
}
