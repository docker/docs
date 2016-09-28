'use strict';
import {
  Repositories as R
  } from 'hub-js-sdk';
const debug = require('debug')('hub:actions:addCollaborator');

export default function addTeamCollaborator(actionContext, {JWT, namespace, name, id, permission}, done) {
  actionContext.dispatch('ADD_COLLAB_START');
  R.addTeamCollaborator(JWT, { namespace, name, group_id: id, permission }, (err, res) => {
    if(err) {
      debug('failed');
    } else {
      debug('succeeded');
      R.getTeamCollaboratorsForRepo(JWT, `${namespace}/${name}`, (getErr, getRes) => {
        if (getErr) {
          // 'User repository does not have any teams yet'
          actionContext.dispatch('COLLAB_RECEIVE_TEAMS', {});
        } else {
          actionContext.dispatch('COLLAB_RECEIVE_TEAMS', getRes.body);
        }
      });
    }
    actionContext.dispatch('ADD_COLLAB_FINISH');
    done();
  });
}
