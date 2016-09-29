'use strict';
import { Builds } from 'hub-js-sdk';
import has from 'lodash/object/has';
const debug = require('debug')('hub:actions:cancelBuild');

export default function cancelBuildAction(actionContext, {JWT, id, namespace, name, build_code}, done) {
  actionContext.dispatch('CANCEL_BUILD_START', id);
  Builds.cancelBuild(JWT, { namespace, name, build_code }, (err, res) => {
    if(err) {
      debug('failed');
      const detail = has(err.response.body, 'detail') ? err.response.body : '';
      actionContext.dispatch('CANCEL_BUILD_ERROR', { id, detail: err.response.body.detail });

    } else {
      debug('succeeded');
      actionContext.dispatch('CANCEL_BUILD_SUCCESS', id);
    }
    done();
  });
}
