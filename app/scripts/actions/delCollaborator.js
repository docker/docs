'use strict';
import {
  Repositories as R
  } from 'hub-js-sdk';

export default function delCollaborator(actionContext, { JWT, namespace, name, username }, done) {
  actionContext.dispatch('DEL_COLLABORATORS_SET_LOADING', username);
  R.delCollaborator(JWT, { namespace, name, username }, (err, res) => {
    if(err) {
      actionContext.dispatch('DEL_COLLABORATORS_SET_ERROR', username);
    } else {
      actionContext.dispatch('DEL_COLLABORATORS_SET_SUCCESS', username);
      R.getCollaboratorsForRepo(JWT, `${namespace}/${name}`, (getErr, getRes) => {
        if(getErr) {
          // 'Org repositories do not have collaborators.'
          actionContext.dispatch('COLLAB_RECEIVE_COLLABORATORS', {});
        } else {
          actionContext.dispatch('COLLAB_RECEIVE_COLLABORATORS', getRes.body);
        }
      });
    }
    done();
  });
}
