/* @flow */
/*global VisibilityFormFieldPayload */
'use strict';

var debug = require('debug')('hub:actions:updateAutoBuildPushTriggerItem');
export default function(actionContext, { isNew, index, fieldkey, value}) {
  if (isNew) {
    actionContext.dispatch('UPDATE_AUTOBUILD_NEW_TAG_ITEM', {index, fieldkey, value});
  } else {
    actionContext.dispatch('UPDATE_AUTOBUILD_PUSH_TRIGGER_ITEM', {index, fieldkey, value});
  }
}
