'use strict';
import async from 'async';
import request from 'superagent';
const debug = require('debug')('createWebhook');

// TODO: REMOVE. This is deprecated in favor of AddPipeline
module.exports = function(actionContext, { JWT, namespace, name, webhookName }) {

  request.post(`${process.env.REGISTRY_API_BASE_URL}/v2/repositories/${namespace}/${name}/webhooks/`)
         .accept('application/json')
         .set('Authorization', `JWT ${JWT}`)
         .send({
           name: webhookName
         })
         .end((err, results) => {
           if(err) {
             debug(err);
             if(err.response.badRequest) {
               debug('badrequest');
             } else {
               debug('facepalm');
             }
           } else {
             actionContext.history.push(`/r/${namespace}/${name}/~/settings/webhooks/`);
           }
         });
};
