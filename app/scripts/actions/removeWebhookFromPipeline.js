'use strict';

export default function removeWebhookFromPipeline({ dispatch },
                                                  params,
                                                  done) {
  dispatch('ADD_WEBHOOK_REMOVE_HOOK');
  done();
}
