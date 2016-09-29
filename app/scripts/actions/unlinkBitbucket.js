'use strict';

import { Builds } from 'hub-js-sdk';
const debug = require('debug')('hub:actions:unlinkBitbucket');
import linkedAccountSettingsAction from './navigate/linkedAccountsSettings';

module.exports = function(actionContext, jwt) {
  Builds.unlinkBitbucket(jwt, function(err, res) {
    if (err) {
      debug(err);
      const { detail } = err.response.body;
      if(detail) {
        actionContext.dispatch('BITBUCKET_UNLINK_ERROR', detail);
      }
    } else if (res.ok) {
      linkedAccountSettingsAction(
        {
          actionContext: actionContext,
          payload: {},
          done: function() { debug('done unlinking bitbucket.'); },
          maybeData: {token: jwt}
        }
      );
    }
  });
};
