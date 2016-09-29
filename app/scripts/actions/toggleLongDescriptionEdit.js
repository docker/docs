'use strict';

export default function (actionContext, { isEditing }) {
  actionContext.dispatch('TOGGLE_LONG_DESCRIPTION_EDIT', { isEditing });
}
