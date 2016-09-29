'use strict';

import request from 'superagent';

export default function getPipelineHistory({
  dispatch, history
}, {
  jwt, namespace, name, slug
}, done) {
  request.get(`${process.env.HUB_API_BASE_URL}/v2/repositories/${namespace}/${name}/webhook_pipeline/${slug}/history/`)
         .set('Authorization', `JWT ${jwt}`)
         .type('json')
         .accept('json')
         .end((err, res) => {
           if (err) {
             return done();
           } else {
             dispatch('RECEIVE_PIPELINE_HISTORY', {slug, payload: res.body});
             return done();
           }
         });
}
