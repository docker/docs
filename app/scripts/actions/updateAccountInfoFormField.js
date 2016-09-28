'use strict';

export default function(actionContext, { fieldKey, fieldValue}) {
  actionContext.dispatch('ACCOUNT_INFO_UPDATE_FIELD_WITH_VALUE', { fieldKey, fieldValue });
}
