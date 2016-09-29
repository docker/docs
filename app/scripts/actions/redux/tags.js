'use strict';

import { Repositories as Repos } from 'hub-js-sdk';
import { DELETE_REPO_TAG } from 'reduxConsts';
const debug = require('debug')('hub:actions:redux:deleteRepoTag');

export const deleteRepoTag = ({ JWT, namespace, name, tagName }) => ({
  type: DELETE_REPO_TAG,
  payload: {
    namespace,
    reponame: name,
    tagName
  },
  meta: {
    sdk: {
      call: Repos.deleteRepoTag,
      args: [JWT, { namespace, name, tagName }],
      callback: (err, res) => ({}),
      statusKey: ['deleteRepoTag', tagName]
    }
  }
});
