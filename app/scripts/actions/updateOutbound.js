'use strict';

export default function(actionContext, newList) {
  actionContext.dispatch('UPDATE_OUTBOUND', newList);
}
