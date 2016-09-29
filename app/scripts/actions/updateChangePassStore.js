'use strict';

var updateStore = function(actionContext, payload) {
  actionContext.dispatch('CHANGE_PASS_UPDATE', payload);
};

module.exports = updateStore;
