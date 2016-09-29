/* @flow */
'use strict';

import type FluxibleActionContext from '../../../flow-libs/fluxible';

module.exports = function(actionContext: FluxibleActionContext) {
  actionContext.dispatch('LOGIN_CLEAR', {});
};
