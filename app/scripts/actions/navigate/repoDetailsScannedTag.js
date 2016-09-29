'use strict';

import { parallel } from 'async';
import { PENDING_DELETE } from 'common/enums/RepoStatus';
import has from 'lodash/object/has';
import merge from 'lodash/object/merge';
import { Repositories as Repos } from 'hub-js-sdk';
import {
  RECEIVE_REPO,
  RECEIVE_SCANNED_TAG_DATA,
  ERROR
} from 'reduxConsts';
const debug = require('debug')('navigate::repoDetailsScannedTagData');
import request from 'superagent';
import { normalize, arrayOf } from 'normalizr';
import { scan } from 'normalizers';

const getRepo = ({maybeToken, actionContext, user, splat}) => (callback) => {
  Repos.getRepo(maybeToken, `${user}/${splat}`, function(err, res) {
    let status;
    if (res && res.body) {
      status = res.body.status;
    }

    if (err || status === PENDING_DELETE) {
      actionContext.dispatch('REPO_NOT_FOUND', err);
      return callback(err);
    } else {
      return callback(null, res.body);
    }
   });
};

/*
 * Get nautilus scan info for the tag.
 * TODO: Once this API is in it's final form for the jan release, move it to the SDK
 */
const getScanForTag = ({actionContext, maybeToken, payload, user, splat, tagname, done}) => (callback) => {
  const namespace = user;
  const name = splat;
  const req = request.get(`${process.env.NAUTILUS_API_BASE_URL}/repositories/result?namespace=${namespace}&reponame=${name}&tag=${tagname}&detailed=1`);
  if (maybeToken) {
    req.set('Authorization', 'JWT ' + maybeToken);
  }
  req.timeout(7000);
  req.end((err, res) => {
    if (err) {
      // NOTE: Suffixing the dispatch type with _STATUS means the status
      //       reducer will record this in the same way as our SDK calls
      //       via the middleware.
      //       When we move this to the SDK, the error will be dispatched automatically
      actionContext.reduxStore.dispatch({
        type: `${RECEIVE_SCANNED_TAG_DATA}_STATUS`,
        payload: {
          status: ERROR,
          statusKey: ['getScanForTag', namespace, name, tagname],
          error: err
        },
        error: true
      });
      return callback(err);
    } else {
      // The API response contains a 'scan' resource within an object
      // inside 'scan_details'
      const { scan_details, image: { reponame, tag } } = res.body;
      const { latest_scan_status } = res.body;
      const result = { latest_scan_status, reponame, tag, ...scan_details };
      const normalized = normalize(result, scan);
      return callback(null, normalized);
    }
  });
};

export default function repoDetailsScannedTag({actionContext, payload, done, maybeData}){
  let token = '';
  if (has(maybeData, 'token')) {
    token = maybeData.token;
  }
  const args = {
    actionContext,
    maybeToken: token,
    user: payload.params.user,
    splat: payload.params.splat,
    tagname: payload.params.tagname
  };

  parallel({
    repo: getRepo(args),
    tagScan: getScanForTag(args)
  }, function(err, res){
    if (err) {
      actionContext.dispatch('REPO_NOT_FOUND', err);
    } else {
      const { repo, tagScan } = res;
      /* REPOS */
      //required for repository page header to display
      actionContext.dispatch('RECEIVE_REPOSITORY', repo);
      // We also need to dispatch to Redux; this will store the current repo
      // within the repos reducer allowing us to find the namespace and repo
      // name for the current route.
      actionContext.reduxStore.dispatch({
        type: RECEIVE_REPO,
        payload: repo
      });

      /* SCANS */
      actionContext.reduxStore.dispatch({
        type: RECEIVE_SCANNED_TAG_DATA,
        payload: tagScan
      });
    }
    done();
  });
}
