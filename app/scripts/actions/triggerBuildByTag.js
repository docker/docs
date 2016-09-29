'use strict';

import { Autobuilds } from 'hub-js-sdk';
var debug = require('debug')('hub:actions:triggerBuildByTag');

module.exports = function(actionContext, {JWT, tag, triggerId}) {
  let successStatus, errorStatus;
  actionContext.dispatch('ATTEMPT_TRIGGER_BY_TAG', triggerId);
  Autobuilds.triggerBuildByTag(JWT, tag, (err, res) => {
    if (err) {
      debug(err);
      actionContext.dispatch('AB_TRIGGER_BY_TAG_ERROR',
        {
          error: 'Error triggering the build. Please check your source name.',
          id: triggerId
        });
    } else {
      if (res.ok) {
        switch(res.status) {
          case 202:
            successStatus = 'Successfully triggered a build.';
            break;
          case 200:
            successStatus = 'Attempted to trigger a build. Please check the build details page for more information.';
            //TODO: report this as error or success once it is clear what exactly this means
            break;
          default:
            break;
        }
        actionContext.dispatch('AB_TRIGGER_BY_TAG_SUCCESS', {success: successStatus, id: triggerId});
      }
    }
  });
};
