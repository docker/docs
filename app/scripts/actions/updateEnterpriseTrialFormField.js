'use strict';

export default function(actionContext, {fieldKey, fieldValue}) {
  actionContext.dispatch('ENTERPRISE_TRIAL_UPDATE_FIELD_WITH_VALUE', { fieldKey, fieldValue });
}
