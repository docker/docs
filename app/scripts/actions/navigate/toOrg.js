'use strict';
var debug = require('debug')('navigate::toOrg');

export default function updateOrgOwner({actionContext, payload, done, maybeData}){
  actionContext.dispatch('UPDATE_TO_ORG_OWNER', {owner: ''});
  done();
}
