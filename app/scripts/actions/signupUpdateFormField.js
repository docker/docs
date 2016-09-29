'use strict';

export default function(actionContext, { fieldKey, fieldValue }, done) {
  actionContext.dispatch('SIGNUP_UPDATE_FIELD_WITH_VALUE', {
    fieldKey, fieldValue
  });
}
