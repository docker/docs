'use strict';
const debug = require('debug')('hub:actions:regenTriggerToken');
import {
  Autobuilds as AutoBuild
  } from 'hub-js-sdk';

function getTriggerStatus(JWT, actionContext, namespace, name) {
  AutoBuild.getTriggerStatus(JWT, namespace, name, function(err, res) {
    if (err) {
      debug('error', err);
    }
    if (res.ok) {
      actionContext.dispatch('RECEIVE_TRIGGER_STATUS', res.body);
    }
  });
}

export default function regenTriggerToken(actionContext, {JWT, namespace, name}) {
  AutoBuild.regenBuildTriggerToken(JWT, namespace, name, function(err, res) {
    if (err) {
      debug('error', err);
    }
    if (res.ok) {
      getTriggerStatus(JWT, actionContext, namespace, name);
    }
  });
}
