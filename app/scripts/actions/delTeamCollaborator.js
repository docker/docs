'use strict';
import {
  Repositories as R
  } from 'hub-js-sdk';

export default function delCollaborator(actionContext, { JWT, namespace, name, group_id }, done) {
  actionContext.dispatch('DEL_COLLABORATORS_SET_LOADING', group_id);
  R.delTeamCollaborator(JWT, { namespace, name, group_id }, (err, res) => {
    if(err) {
      actionContext.dispatch('DEL_COLLABORATORS_SET_ERROR', group_id);
    } else {
      actionContext.dispatch('DEL_COLLABORATORS_SET_SUCCESS', group_id);
      R.getTeamCollaboratorsForRepo(JWT, `${namespace}/${name}`, (getErr, getRes) => {
        if (getErr) {
          // 'User repository does not have any teams yet'
          actionContext.dispatch('COLLAB_RECEIVE_TEAMS', {});
        } else {
          actionContext.dispatch('COLLAB_RECEIVE_TEAMS', getRes.body);
        }
      });
    }
    done();
  });
}
