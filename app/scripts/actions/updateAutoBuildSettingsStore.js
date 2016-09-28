'use strict';

export default function(actionContext, { field, key, value}) {
  actionContext.dispatch('UPDATE_AUTO_BUILD_SETTINGS', {field, key, value});
}
