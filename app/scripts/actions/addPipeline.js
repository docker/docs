'use strict';

import includes from 'lodash/collection/includes';

import createPipeline from '@dux/hub-sdk/webhooks/createPipeline';
import MissingArgError from '@dux/hub-sdk/utils/MissingArgError';
import ValidationError from '@dux/hub-sdk/utils/ValidationError';

function handleKnownErrors(err, dispatch) {
  if(err instanceof MissingArgError) {
    if(includes(err.missingArgs, 'namespace') || includes(err.missingArgs, 'name')) {
      /**
       * TODO: send something to bugsnag; The user doesn't control these values
       * This can happen if an engineer forgets to pass in a required value that
       * the user can't control
       */
      dispatch('ADD_WEBHOOK_ERROR');
    } else {
      dispatch('ADD_WEBHOOK_MISSING_ARGS', err.missingArgs);
    }
  } else if (err instanceof ValidationError) {
    dispatch('ADD_WEBHOOK_VALIDATION_ERRORS', err.validationErrors);
  } else {
    // unknown error
    dispatch('ADD_WEBHOOK_ERROR', err);
  }
}

export default function addPipeline({ dispatch, history },
                                    {
                                      jwt,
                                      namespace,
                                      name,
                                      pipelineName,
                                      expectFinalCallback,
                                      webhooks
                                    },
                                   done) {
  dispatch('ADD_WEBHOOK_START');
  createPipeline(jwt,
              {
                namespace,
                name,
                pipelineName,
                expectFinalCallback,
                webhooks
              },
              (err, res) => {
                if(err || !res.ok) {
                  handleKnownErrors(err, dispatch);
                  done();
                } else {
                  dispatch('ADD_WEBHOOK_SUCCESS');
                  history.push(`/r/${namespace}/${name}/~/settings/webhooks/`);
                  done();
                }
              });
}
