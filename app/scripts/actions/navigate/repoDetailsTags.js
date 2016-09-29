'use strict';
const debug = require('debug')('navigate::repo');
import { parallel } from 'async';
import { PENDING_DELETE } from 'common/enums/RepoStatus';
import has from 'lodash/object/has';
import merge from 'lodash/object/merge';
import request from 'superagent';
import {
  RECEIVE_NAUTILUS_TAGS_FOR_REPOSITORY,
  RECEIVE_TAGS_FOR_REPOSITORY,
  RECEIVE_REPO
} from 'reduxConsts.js';
import { normalize, arrayOf } from 'normalizr';
import { tag } from 'normalizers';
import {
  Repositories as Repos
} from 'hub-js-sdk';

const getRepo = ({maybeToken, actionContext, user, splat}) => (callback) => {
  Repos.getRepo(maybeToken, `${user}/${splat}`, function(err, res) {
    let status;
    if (res && res.body) {
      status = res.body.status;
    }

    if (err || status === PENDING_DELETE) {
      return callback(err);
    } else {
      return callback(null, res.body);
    }
   });
};

// getTagsForRepo uses the hub API for loading tag information.  This is used to
// show information about tags that are not scanned.
//
// We use this and the nautilus API because each API has incomplete information:
//   - The nautilus API has vulnerability information
//   - This API has the image size
//
// NOTE: This dispatches to Redux reducers and fluxible stores in the final callback
const getTagsForRepo = ({actionContext, maybeToken, user, splat}) => (callback) => {
  const namespace = user;
  const reponame = splat;
  Repos.getTagsForRepo(maybeToken, `${user}/${splat}`, function(err, res) {
    if (err) {
      return callback(err);
    }
    const { results } = res.body;
    // TODO: Assert normalize works as expected via jest/mocha
    const tags = normalize(results, arrayOf(tag));
    return callback(null, { namespace, reponame, tags });
  });
};

// getNautilusTagsForRepo uses the nautilus API to load tags and their
// vulnerability information via the nautilus API.
//
// NOTE: This dispatches to Redux reducers, not fluxible stores
const getNautilusTagsForRepo = ({actionContext, maybeToken, user, splat}) => (callback) => {
  const namespace = user;
  const reponame = splat;
  const req = request.get(`${process.env.NAUTILUS_API_BASE_URL}/repositories/summaries/${namespace}/${reponame}`);
  if (maybeToken) {
    req.set('Authorization', 'JWT ' + maybeToken);
  }
  req.timeout(7000);
  req.end((err, res) => {
    if (err) {
      // Nautilus call does NOT trigger an error page
      return callback(null, null);
    }
    // TODO: Assert normalize works as expected via jest/mocha
    const tags = normalize(
      res.body,
      arrayOf(tag),
      {
        // TODO: Add records
        // NOTE: The nautius API uses a field called 'tag' to represent the
        // tag name, whereas the HUB api uses 'name'.
        //
        // Our record, frontend code and normalizr key expects us to use
        // the 'name' field. This normalizes the record to tag.
        assignEntity: (obj, key, val) => {
          obj[key] = val;
          if (key === 'tag') {
            obj.name = val;
            delete obj.tag;
          }
        }
      }
    );
    return callback(null, { namespace, reponame, tags });
  });
};

export default function repoDetailsTags({actionContext, payload, done, maybeData}) {
  debug('maybeData:', maybeData);
  let token = '';
  if (has(maybeData, 'token')) {
    token = maybeData.token;
  }
  const { user, splat } = payload.params;
  const args = {
    actionContext,
    maybeToken: token,
    user,
    splat
  };

  parallel({
    repo: getRepo(args),
    tags: getTagsForRepo(args),
    scans: getNautilusTagsForRepo(args)
  }, (err, res) => {
    if (err) {
      // Tags or repo error
      actionContext.dispatch('REPO_NOT_FOUND', err);
    } else {
      const { repo, tags, scans } = res;
      /* REPO */
      actionContext.dispatch('RECEIVE_REPOSITORY', repo);
      // We also need to dispatch to Redux; this will store the current repo
      // within the repos reducer allowing us to find the namespace and repo
      // name for the current route.
      actionContext.reduxStore.dispatch({
        type: RECEIVE_REPO,
        payload: repo
      });

      /* TAGS */
      actionContext.reduxStore.dispatch({
        type: RECEIVE_TAGS_FOR_REPOSITORY,
        payload: tags
      });

      /* SCANS */
      if (scans) {
        actionContext.reduxStore.dispatch({
          type: RECEIVE_NAUTILUS_TAGS_FOR_REPOSITORY,
          payload: scans
        });
      }
    }
    done();
  });
}
