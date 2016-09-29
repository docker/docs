'use strict';
import { Repositories } from 'hub-js-sdk';
const { addCollaborator, getCollaboratorsForRepo } = Repositories;
import has from 'lodash/object/has';
const debug = require('debug')('hub:actions:addCollaborator');

export default function addCollaboratorAction(actionContext, {JWT, namespace, name, user}, done) {
  actionContext.dispatch('ADD_COLLAB_START');
  addCollaborator(JWT, { namespace, name, user }, (err, res) => {
    if(err) {
      if(has(err.response.body, 'detail')) {
        debug('failed');
        actionContext.dispatch('ADD_COLLAB_ERROR', err.response.body.detail);
      }
    } else {
      debug('succeeded');
      actionContext.dispatch('ADD_COLLAB_SUCCESS');
      getCollaboratorsForRepo(JWT, `${namespace}/${name}`, (getErr, getRes) => {
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
