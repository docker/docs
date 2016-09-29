'use strict';

import { createSelector } from 'reselect';
import { Map, List } from 'immutable';
import filter from 'lodash/collection/filter';
import size from 'lodash/collection/size';
import map from 'lodash/collection/map';
import values from 'lodash/object/values';


// Returns all repository tags (unordered)
export const getRepoTags = (state) => {
  const reponame = state.repos.get('name', '');
  const namespace = state.repos.get('namespace', '');
  // We need to use toJS() to deeply convert tags from immutable to objects.
  // We also return an array because getScannedTags and getUnscannedTags return
  // arrays - keeping things consistent.
  return values(state.tags.getIn([namespace, reponame, 'tags'], new Map()).toJS());
};

// Returns all repository tags in the order of the hub API
// Note: This does _not_ return any tags that are returned by nautilus but not hub
//       so we use the getRepoTags for the scannedTag selector
export const getRepoTagsInOrder = (state) => {
  const reponame = state.repos.get('name', '');
  const namespace = state.repos.get('namespace', '');
  // We need to use toJS() to deeply convert tags from immutable to objects.
  // We also return an array because getScannedTags and getUnscannedTags return
  // arrays - keeping things consistent.
  let orderedTags = state.tags.getIn([namespace, reponame, 'result'], []);
  if (orderedTags.toArray) {
    orderedTags = orderedTags.toArray();
  }
  const tags = state.tags.getIn([namespace, reponame, 'tags'], new Map()).toJS();
  return map(orderedTags, (tagId) => tags[tagId]);
};


// Returns only tags which have been scanned by nautilus
export const getScannedTags = createSelector(
  [getRepoTags],
  (tags) => {
    // If the tag has a 'healthy' key then this has been scanned by nautilus
    return filter(tags, (tag) => tag.healthy !== undefined);
  }
);
// Number of tags scanned by nautilus
export const getScannedTagCount = createSelector(
  [getScannedTags],
  (tags) => size(tags)
);

// getUnscannedTags returns only tags **not** scanned by nautilus
export const getUnscannedTags = createSelector(
  [getRepoTagsInOrder],
  (tags) => {
    // If healthy is undefined this tag only has a hub response
    return filter(tags, (tag) => tag.healthy === undefined);
  }
);

export const getUnscannedTagCount = createSelector(
  [getUnscannedTags],
  (tags) => size(tags)
);
