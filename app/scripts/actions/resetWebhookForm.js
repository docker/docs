'use strict';

export default function resetWebhookForm({ dispatch },
                                             params,
                                             done) {
  dispatch('ADD_WEBHOOK_RESET');
  done();
}
