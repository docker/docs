'use strict';

export default function(actionContext, payload) {
  actionContext.dispatch('UPDATE_TO_ORG_OWNER', payload);
}
