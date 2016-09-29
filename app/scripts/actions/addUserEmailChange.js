'use strict';

export default function(actionContext, { email }) {
  actionContext.dispatch('UPDATE_ADD_EMAIL', email);
}
