/* @flow */
'use strict';

//TODO: organization object needs to be documented in hub-js-sdk when available
module.exports = function(actionContext: {dispatch: Function},
                          orgName: {orgName: string}) {
  actionContext.dispatch('SELECT_ORGANIZATION', orgName);
};
