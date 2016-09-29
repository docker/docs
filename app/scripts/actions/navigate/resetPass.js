'use strict';

export default function resetPass({actionContext, payload, done, maybeData}){
  actionContext.dispatch('CHANGE_PASS_CLEAR');
  done();
}
