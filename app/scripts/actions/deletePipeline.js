'use strict';

import request from 'superagent';

export default function deletePipeline({
  dispatch, history
}, {
  jwt, namespace, name, slug
}, done) {
  dispatch('DELETE_PIPELINE_ATTEMPTING');
  request.del(`${process.env.HUB_API_BASE_URL}/v2/repositories/${namespace}/${name}/webhook_pipeline/${slug}/`)
         .set('Authorization', `JWT ${jwt}`)
         .type('json')
         .accept('json')
         .end((err, res) => {
           if(err) {
             dispatch('DELETE_PIPELINE_FACEPALM');
             done();
           } else {
             dispatch('DELETE_PIPELINE_SUCCESS');
             history.push(`/r/${namespace}/${name}/~/settings/webhooks/`);
             done();
           }
         });
}
