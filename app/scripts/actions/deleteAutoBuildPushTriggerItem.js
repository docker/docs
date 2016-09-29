'use strict';

var debug = require('debug')('hub:actions:deleteAutoBuildPushTriggerItem');

export default function(actionContext, {isNew, index}) {
  if (isNew) {
    actionContext.dispatch('DELETE_AUTOBUILD_NEW_TAG_ITEM', index);
  } else {
    actionContext.dispatch('DELETE_AUTOBUILD_PUSH_TRIGGER_ITEM', index);
  }
}
