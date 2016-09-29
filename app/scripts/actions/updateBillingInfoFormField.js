'use strict';

export default function(actionContext, { field, fieldKey, fieldValue}) {
  if (field === 'card' && fieldKey === 'number') {
    var type = window.recurly.validate.cardType(fieldValue.toString());
    actionContext.dispatch('BILLING_INFO_UPDATE_FIELD_WITH_VALUE', {
      field: 'card',
      fieldKey: 'type',
      fieldValue: type
    });
  }
  actionContext.dispatch('BILLING_INFO_UPDATE_FIELD_WITH_VALUE', {
    field, fieldKey, fieldValue
  });
}
