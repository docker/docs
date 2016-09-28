'use strict';

import { Builds } from 'hub-js-sdk';
const debug = require('debug')('hub:actions:unlinkGithub');
import has from 'lodash/object/has';
import linkedAccountSettingsAction from './navigate/linkedAccountsSettings';

module.exports = function(actionContext, jwt) {
  Builds.unlinkGithub(jwt, function(err, res) {
    if (err) {
      debug(err);
      const { detail } = err.response.body;
      if(detail) {
        actionContext.dispatch('GITHUB_UNLINK_ERROR', detail);
      }
    } else if (res.ok) {
      linkedAccountSettingsAction(
        {
          actionContext: actionContext,
          payload: {},
          done: function() { debug('done unlinking github account.'); },
          maybeData: {token: jwt}
        }
      );
    }
  });
};
