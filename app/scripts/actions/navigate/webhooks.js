'use strict';

import _ from 'lodash';
import { parallel } from 'async';
import { PENDING_DELETE } from 'common/enums/RepoStatus';
import {
  Repositories as Repos
} from 'hub-js-sdk';
import getPipelines from '@dux/hub-sdk/webhooks/getPipelines';

var debug = require('debug')('action::webhooks');

const fetchRepo = ({ token, namespace, name }) => callback => {
  Repos.getRepo(token, `${namespace}/${name}`, (err, { body }) => {
    const { status } = body;
    if (err || status === PENDING_DELETE) {
      // We should handle other possible errors here (500, etc)
      return callback('REPO_NOT_FOUND');
    } else {
      return callback(null, body);
    }
  });
};

const fetchPipelines = ({ token, namespace, name }) => callback => {
  getPipelines(token, {
    namespace,
    name
  }, (err, res) => {
    if (err) {
      return callback();
    } else {
      return callback(null, res.body);
    }
  });
};

export default function fetchWebhooksPageData({
  actionContext: { dispatch },
  payload: { params },
  done,
  maybeData: { token, user }
}) {

  if (!token) {
    dispatch('REPO_NOT_FOUND', null);
    return done();
  }

  debug(params);
  const {
    user: namespace,
    splat: name
  } = params;

  parallel({
    repository: fetchRepo({ token, namespace, name }),
    pipelines: fetchPipelines({ token, namespace, name })
  },
           (err, { repository, pipelines }) => {
             if(err) {
               debug('err', err);
               /**
                * If there's an error, 404 by default, but handle other
                * errors in the future
                */
               dispatch('REPO_NOT_FOUND', null);
               done();
             } else {
               debug('receive', repository, pipelines);
               dispatch('RECEIVE_REPOSITORY', repository);
               dispatch('RECEIVE_WEBHOOKS', pipelines);
               done();
             }
           });
}
