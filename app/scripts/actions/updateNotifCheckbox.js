'use strict';

export default function(actionContext, cboxType) {
  actionContext.dispatch('NOTIF_CHECKBOX_CLICK', cboxType);
}
