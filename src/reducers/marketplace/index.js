import {
  MARKETPLACE_FETCH_BUNDLE_DETAIL,
  MARKETPLACE_FETCH_BUNDLE_SUMMARY,
  MARKETPLACE_FETCH_CATEGORIES,
  MARKETPLACE_FETCH_PLATFORMS,
  MARKETPLACE_FETCH_REPOSITORY_DETAIL,
  MARKETPLACE_FETCH_REPOSITORY_SUMMARY,
  MARKETPLACE_SEARCH,
} from 'actions/marketplace';
import {
  REPOSITORY_FETCH_IMAGE_DETAIL,
  REPOSITORY_FETCH_COMMENTS,
  REPOSITORY_FETCH_IMAGE_TAGS,
} from 'actions/repository';
import {
  NAUTILUS_FETCH_SCAN_DETAIL,
  NAUTILUS_FETCH_TAGS_AND_SCANS,
} from 'actions/nautilus';
import cloneDeep from 'lodash/cloneDeep';
import forEach from 'lodash/forEach';
import makeRepositoryId from 'lib/utils/repo-image-name.js';
import merge from 'lodash/merge';
import set from 'lodash/set';
import get from 'lodash/get';
import setWith from 'lodash/setWith';

export const DEFAULT_STATE = {
  bundles: {},
  filters: {
    categories: {},
    platforms: {},
  },
  images: {
    certified: {},
    community: {},
  },
  search: {
    isFetching: false,
    pages: {},
  },
};

//
// https://visionmedia.github.io/superagent/#error-handling
const reduceError = (payload) => {
  // Status will be a code like 404, etc.
  const status = get(payload, ['response', 'error', 'status']);
  // Message may come from the API
  const message = get(payload, ['response', 'error', 'message']);
  // TODO Kristie 6/7/16 Handle more error cases from different APIs
  return message || status;
};

// Helper function for saving the results of a paginated call
// path should be an array
const reducePaginatedResults = ({ nextState, path, count, page, newPage }) => {
  set(nextState, [...path, 'isFetching'], false);
  set(nextState, [...path, 'count'], count);
  // If 'pages' doesn't exist, use an object instead of an array
  setWith(nextState, [...path, 'pages', page], newPage, Object);
  return nextState;
};

const reduceSearch = ({ nextState, meta, payload }) => {
  const page = meta && meta.page || payload && payload.page || 1;
  const { count, summaries = [] } = payload;
  const newPage = { results: summaries };
  const path = ['search'];
  return reducePaginatedResults({ nextState, path, count, page, newPage });
};

const reduceCommunityComments = ({ nextState, meta, payload, id }) => {
  const page = meta && meta.page || payload && payload.page || 1;
  const { count, results } = payload;
  const newPage = { results };
  const path = ['images', 'community', id, 'comments'];
  return reducePaginatedResults({ nextState, path, count, page, newPage });
};

const reduceCertifiedComments = ({ nextState, meta, payload, id }) => {
  const page = meta && meta.page || payload && payload.page || 1;
  const { count, results } = payload;
  const newPage = { results };
  const path = ['images', 'certified', id, 'comments'];
  return reducePaginatedResults({ nextState, path, count, page, newPage });
};

const reduceScanDetail = ({ nextState, meta, payload, id }) => {
  const path = ['images', 'certified', id, 'scanDetail', meta.tag];
  const scanResults = { ...payload.entities, isFetching: false };
  setWith(nextState, path, scanResults, Object);
  return nextState;
};

const reduceCertifiedTags = ({ nextState, meta, payload, id }) => {
  const page = meta && meta.page || payload && payload.page || 1;
  const { count, results } = payload;
  const newPage = { results };
  const path = ['images', 'certified', id, 'tags'];
  return reducePaginatedResults({ nextState, path, count, page, newPage });
};

const reduceCommunityTags = ({ nextState, meta, payload, id }) => {
  const page = meta && meta.page || payload && payload.page || 1;
  const { count, results } = payload;
  const newPage = { results };
  const path = ['images', 'community', id, 'tags'];
  return reducePaginatedResults({ nextState, path, count, page, newPage });
};

const reduceTagsAndScans = ({ nextState, meta, payload, id }) => {
  const page = meta && meta.page || payload && payload.page || 1;
  const [tags, scans] = payload;
  const { count, results } = tags;
  // keep track of the order of the results in orderedIds
  const newPage = { orderedIds: [] };
  const path = ['images', 'certified', id, 'tagsAndScans'];
  // Take both payloads from the getTags and getScans calls and merge any
  // available scans into the tags to display them
  forEach(results, (tag) => {
    newPage.orderedIds.push(tag.name);
    newPage[tag.name] = tag;
  });
  forEach(scans, (scan) => {
    // Nautilus uses "tag" and Hub uses "name" to represent the tagname
    const { tag } = scan;
    const existingTag = newPage[tag];
    if (existingTag) {
      // add nautilus scan information if we got this tag from the hub api
      newPage[tag] = merge(existingTag, scan);
    }
  });
  return reducePaginatedResults({ nextState, path, count, page, newPage });
};

const formatFilter = (filterArray) => {
  const filters = {};
  filterArray.forEach(({ name, label }) => {
    filters[name] = label;
  });
  return filters;
};

export default function marketplace(state = DEFAULT_STATE, action) {
  const nextState = cloneDeep(state);
  const { meta, payload, type } = action;
  // Community image (combo of namespace + reponame)
  let repositoryId;
  // Certified image
  let productId;
  let currentRepository;
  let currentBundle;
  // Grab the current bundle or repository's results, if they exist
  if (meta && meta.reponame && meta.namespace) {
    repositoryId = makeRepositoryId(meta);
    currentRepository = nextState.images.community[repositoryId];
  }
  if (meta && (meta.id || meta.isCertified)) {
    productId = meta.id;
    // It's fine for either of these to be undefined
    currentRepository = nextState.images.certified[productId];
    currentBundle = nextState.bundles[productId];
  }

  switch (type) {
    //--------------------------------------------------------------------------
    // MARKETPLACE_FETCH_CATEGORIES
    //--------------------------------------------------------------------------
    case `${MARKETPLACE_FETCH_CATEGORIES}_ACK`:
      nextState.filters.categories = formatFilter(payload);
      return nextState;

    //--------------------------------------------------------------------------
    // MARKETPLACE_FETCH_PLATFORMS
    //--------------------------------------------------------------------------
    case `${MARKETPLACE_FETCH_PLATFORMS}_ACK`:
      nextState.filters.platforms = formatFilter(payload);
      return nextState;

    //--------------------------------------------------------------------------
    // MARKETPLACE_FETCH_REPOSITORY_DETAIL
    //--------------------------------------------------------------------------
    case `${MARKETPLACE_FETCH_REPOSITORY_DETAIL}_REQ`:
      nextState.images.certified[productId] = merge({}, currentRepository, {
        isFetching: true,
      });
      return nextState;
    case `${MARKETPLACE_FETCH_REPOSITORY_DETAIL}_ACK`:
      nextState.images.certified[productId] = merge({}, currentRepository, {
        isFetching: false,
        ...payload,
      });
      return nextState;
    case `${MARKETPLACE_FETCH_REPOSITORY_DETAIL}_ERR`:
      nextState.images.certified[productId].isFetching = false;
      nextState.images.certified[productId].error = reduceError(payload);
      return nextState;

    //--------------------------------------------------------------------------
    // MARKETPLACE_FETCH_REPOSITORY_SUMMARY
    //--------------------------------------------------------------------------
    case `${MARKETPLACE_FETCH_REPOSITORY_SUMMARY}_REQ`:
      nextState.images.certified[productId] = merge({}, currentRepository, {
        isFetching: true,
      });
      return nextState;
    case `${MARKETPLACE_FETCH_REPOSITORY_SUMMARY}_ACK`:
      nextState.images.certified[productId] = merge({}, currentRepository, {
        isFetching: false,
        ...payload,
      });
      return nextState;
    case `${MARKETPLACE_FETCH_REPOSITORY_SUMMARY}_ERR`:
      nextState.images.certified[productId].isFetching = false;
      nextState.images.certified[productId].error = reduceError(payload);
      return nextState;

    //--------------------------------------------------------------------------
    // MARKETPLACE_FETCH_BUNDLE_DETAIL
    //--------------------------------------------------------------------------
    case `${MARKETPLACE_FETCH_BUNDLE_DETAIL}_REQ`:
      nextState.bundles[productId] = merge({}, currentBundle, {
        isFetching: true,
      });
      return nextState;
    case `${MARKETPLACE_FETCH_BUNDLE_DETAIL}_ACK`:
      nextState.bundles[productId] = merge({}, currentBundle, {
        isFetching: false,
        ...payload,
      });
      return nextState;
    case `${MARKETPLACE_FETCH_BUNDLE_DETAIL}_ERR`:
      nextState.bundles[productId].isFetching = false;
      nextState.bundles[productId].error = reduceError(payload);
      return nextState;

    //--------------------------------------------------------------------------
    // MARKETPLACE_FETCH_BUNDLE_SUMMARY
    //--------------------------------------------------------------------------
    case `${MARKETPLACE_FETCH_BUNDLE_SUMMARY}_REQ`:
      nextState.bundles[productId] = merge({}, currentBundle, {
        isFetching: true,
      });
      return nextState;
    case `${MARKETPLACE_FETCH_BUNDLE_SUMMARY}_ACK`:
      nextState.bundles[productId] = merge({}, currentBundle, {
        isFetching: false,
        ...payload,
      });
      return nextState;
    case `${MARKETPLACE_FETCH_BUNDLE_SUMMARY}_ERR`:
      nextState.bundles[productId].isFetching = false;
      nextState.bundles[productId].error = reduceError(payload);
      return nextState;


    //--------------------------------------------------------------------------
    // MARKETPLACE_SEARCH
    //--------------------------------------------------------------------------
    case `${MARKETPLACE_SEARCH}_REQ`:
      nextState.search.isFetching = true;
      return nextState;
    case `${MARKETPLACE_SEARCH}_ACK`:
      return reduceSearch({ nextState, meta, payload });
    case `${MARKETPLACE_SEARCH}_ERR`:
      nextState.search.isFetching = false;
      nextState.search.error = reduceError(payload);
      return nextState;

    //--------------------------------------------------------------------------
    // REPOSITORY_FETCH_IMAGE_DETAIL
    //--------------------------------------------------------------------------
    case `${REPOSITORY_FETCH_IMAGE_DETAIL}_REQ`:
      nextState.images.community[repositoryId] = merge({}, currentRepository, {
        isFetching: true,
      });
      return nextState;
    case `${REPOSITORY_FETCH_IMAGE_DETAIL}_ACK`:
      nextState.images.community[repositoryId] = merge({}, currentRepository, {
        isFetching: false,
        ...payload,
        // Hub uses `name` instead of `reponame`
        reponame: payload.name,
      });
      return nextState;
    case `${REPOSITORY_FETCH_IMAGE_DETAIL}_ERR`:
      nextState.images.community[repositoryId].isFetching = false;
      nextState.images.community[repositoryId].error = reduceError(payload);
      return nextState;

    //--------------------------------------------------------------------------
    // REPOSITORY_FETCH_IMAGE_TAGS
    //--------------------------------------------------------------------------
    case `${REPOSITORY_FETCH_IMAGE_TAGS}_REQ`:
      // TODO Kristie 6/7/16 Remove isCertified once marketplace tags Service
      // is up (this will be a different dispatch / API)
      if (meta && meta.isCertified) {
        nextState.images.certified[productId] =
          merge({}, currentRepository, { tags: { isFetching: true } });
      } else {
        nextState.images.community[repositoryId] =
          merge({}, currentRepository, { tags: { isFetching: true } });
      }
      return nextState;
    case `${REPOSITORY_FETCH_IMAGE_TAGS}_ACK`:
      // TODO Kristie 6/7/16 Remove isCertified once marketplace tags Service
      // is up (this will be a different dispatch / API)
      if (meta && meta.isCertified) {
        return reduceCertifiedTags({ nextState, meta, payload, id: productId });
      }
      return reduceCommunityTags({
        nextState, meta, payload, id: repositoryId,
      });
    case `${REPOSITORY_FETCH_IMAGE_TAGS}_ERR`:
      if (meta && meta.isCertified) {
        nextState.images.certified[productId].tags.isFetching = false;
        nextState.images.certified[productId].tags.error =
          reduceError(payload);
      } else {
        nextState.images.community[repositoryId].tags.isFetching = false;
        nextState.images.community[repositoryId].tags.error =
          reduceError(payload);
      }
      return nextState;

    //--------------------------------------------------------------------------
    // REPOSITORY_FETCH_COMMENTS
    //--------------------------------------------------------------------------
    case `${REPOSITORY_FETCH_COMMENTS}_REQ`:
      if (meta && meta.isCertified) {
        nextState.images.certified[productId] =
          merge({}, currentRepository, { comments: { isFetching: true } });
      } else {
        nextState.images.community[repositoryId] =
          merge({}, currentRepository, { comments: { isFetching: true } });
      }
      return nextState;
    case `${REPOSITORY_FETCH_COMMENTS}_ACK`:
      if (meta && meta.isCertified) {
        return reduceCertifiedComments({
          nextState, meta, payload, id: productId,
        });
      }
      return reduceCommunityComments({
        nextState, meta, payload, id: repositoryId,
      });
    case `${REPOSITORY_FETCH_COMMENTS}_ERR`:
      if (meta && meta.isCertified) {
        nextState.images.certified[productId].comments.isFetching = false;
        nextState.images.certified[productId].comments.error =
          reduceError(payload);
      } else {
        nextState.images.community[repositoryId].comments.isFetching = false;
        nextState.images.community[repositoryId].comments.error =
          reduceError(payload);
      }
      return nextState;

    //--------------------------------------------------------------------------
    // NAUTILUS_FETCH_SCAN_DETAIL
    //--------------------------------------------------------------------------
    case `${NAUTILUS_FETCH_SCAN_DETAIL}_REQ`:
      // Save scan results in the tree like
      // images.someRepoId.scanDetail.someTagName
      setWith(
        nextState,
        ['images', 'certified', productId, 'scanDetail', meta.tag],
        { isFetching: true },
        Object,
      );
      return nextState;
    case `${NAUTILUS_FETCH_SCAN_DETAIL}_ACK`:
      return reduceScanDetail({ nextState, meta, payload, id: productId });
    case `${NAUTILUS_FETCH_SCAN_DETAIL}_ERR`:
      // TODO Kristie 5/11/16 Handle Errors properly!!
      setWith(
        nextState,
        ['images', 'certified', productId, 'scanDetail', meta.tag],
        { isFetching: false, error: payload.error },
        Object,
      );
      return nextState;

    //--------------------------------------------------------------------------
    // NAUTILUS_FETCH_TAGS_AND_SCANS
    //--------------------------------------------------------------------------
    case `${NAUTILUS_FETCH_TAGS_AND_SCANS}_REQ`:
      setWith(
        nextState,
        ['images', 'certified', productId, 'tagsAndScans'],
        { isFetching: true },
        Object,
      );
      return nextState;
    case `${NAUTILUS_FETCH_TAGS_AND_SCANS}_ACK`:
      return reduceTagsAndScans({ nextState, meta, payload, id: productId });
    case `${NAUTILUS_FETCH_TAGS_AND_SCANS}_ERR`:
      // TODO Kristie 5/11/16 Handle Errors properly!!
      setWith(
        nextState,
        ['images', 'certified', productId, 'tagsAndScans'],
        { isFetching: false, error: payload.error },
        Object,
      );
      return nextState;

    default:
      return state;
  }
}
