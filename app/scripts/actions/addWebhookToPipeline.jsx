'use strict';

export default function addWebhookToPipeline({ dispatch },
                                             params,
                                             done) {
  dispatch('ADD_WEBHOOK_NEW_HOOK');
  done();
}
