'use strict';

export default function(actionContext, { fieldKey, fieldValue}) {
  actionContext.dispatch('CREATE_REPO_UPDATE_FIELD_WITH_VALUE', {
    fieldKey, fieldValue
  });
}
