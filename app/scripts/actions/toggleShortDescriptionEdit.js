'use strict';

export default function (actionContext, { isEditing }) {
  actionContext.dispatch('TOGGLE_SHORT_DESCRIPTION_EDIT', { isEditing });
}

