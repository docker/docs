'use strict';

export default function(actionContext, { fieldKey, fieldValue}) {
  actionContext.dispatch('ADD_ORG_UPDATE_FIELD_WITH_VALUE', { fieldKey, fieldValue });
}
