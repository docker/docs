'use strict';
import _ from 'lodash';

export default function(actionContext, slateArray) {
  if (_.includes(slateArray, 'notifications')) {
    actionContext.dispatch('RESET_EMAIL_NOTIFICATIONS_STORE');
  }
  if (_.includes(slateArray, 'outbound')) {
    actionContext.dispatch('RESET_OUTBOUND_EMAILS_STORE');
  }
}
