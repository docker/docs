'use strict';

export default function(actionContext, { fieldKey, fieldValue}) {
  actionContext.dispatch('AUTOBUILD_FORM_UPDATE_FIELD_WITH_VALUE', {
    fieldKey, fieldValue
  });
}
