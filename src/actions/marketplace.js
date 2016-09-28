const request = require('superagent-promise')(require('superagent'), Promise);
import { bearer } from 'lib/utils/authHeaders';
import { DEFAULT_SEARCH_PAGE_SIZE } from 'lib/constants/defaults';
import { OFFICIAL } from 'lib/constants/searchFilters/sources';
// TODO Kristie 5/17/16 Use actual API when it is ready
import searchPlatformFilters from 'lib/constants/searchFilters/platforms';
import isStaging from 'lib/utils/isStaging';
import sanitize from 'lib/utils/remove-undefined';
import trim from 'lodash/trim';

/* eslint-disable max-len */
export const MARKETPLACE_FETCH_AUTOCOMPLETE_SUGGESTIONS = 'MARKETPLACE_FETCH_AUTOCOMPLETE_SUGGESTIONS';
export const MARKETPLACE_FETCH_BUNDLE_DETAIL = 'MARKETPLACE_FETCH_BUNDLE_DETAIL';
export const MARKETPLACE_FETCH_BUNDLE_SUMMARY = 'MARKETPLACE_FETCH_BUNDLE_SUMMARY';
export const MARKETPLACE_FETCH_CATEGORIES = 'MARKETPLACE_FETCH_CATEGORIES';
export const MARKETPLACE_FETCH_MOST_POPULAR = 'MARKETPLACE_FETCH_MOST_POPULAR';
export const MARKETPLACE_FETCH_FEATURED = 'MARKETPLACE_FETCH_FEATURED';
export const MARKETPLACE_FETCH_PLATFORMS = 'MARKETPLACE_FETCH_PLATFORMS';
export const MARKETPLACE_FETCH_REPOSITORY_DETAIL = 'MARKETPLACE_FETCH_REPOSITORY_DETAIL';
export const MARKETPLACE_FETCH_REPOSITORY_SUMMARY = 'MARKETPLACE_FETCH_REPOSITORY_SUMMARY';
export const MARKETPLACE_SEARCH = 'MARKETPLACE_SEARCH';
export const MARKETPLACE_CREATE_REPOSITORY = 'MARKETPLACE_CREATE_REPOSITORY';
export const MARKETPLACE_EDIT_REPOSITORY = 'MARKETPLACE_EDIT_REPOSITORY';
export const MARKETPLACE_DELETE_REPOSITORY = 'MARKETPLACE_DELETE_REPOSITORY';

const MARKETPLACE_HOST = '';
const MARKETPLACE_BASE_URL = `${MARKETPLACE_HOST}/api/content/v1`;

const MARKETPLACE_PRIVATE_HOST = process.env.NODE_ENV !== 'production' || isStaging() ?
  'https://mercury-content.s.stage-us-east-1.aws.dckr.io' :
  'https://mercury-content.s.us-east-1.aws.dckr.io';
const MARKETPLACE_PRIVATE_BASE_URL = `${MARKETPLACE_PRIVATE_HOST}/api/private/content`;
/* eslint-enable max-len */

export function marketplaceSearch({
    q,
    category,
    order,
    page,
    page_size = DEFAULT_SEARCH_PAGE_SIZE,
    platform,
    sort,
    // Unless explicitly specified, do not include community results
    source = OFFICIAL,
  }) {
  // Trim whitespace from query
  // eslint-disable-next-line no-param-reassign
  q = trim(q);
  const params = {
    q,
    category,
    order,
    page,
    page_size,
    platform,
    sort,
    source,
  };

  // Track browse and search separately
  if (q === '' && page_size === '99') {
    analytics.track('browse', sanitize({ category, platform, source }));
  } else {
    analytics.track('search', sanitize({ q, category, platform, source }));
  }
  const url = `${MARKETPLACE_BASE_URL}/repositories/search`;
  return {
    type: MARKETPLACE_SEARCH,
    meta: { q, page, page_size },
    payload: {
      promise:
        request
          .get(url)
          .accept('application/json')
          .set(bearer())
          .query(params)
          .end()
          .then((res) => res.body),
    },
  };
}

// Fetch the suggestions for the drop down on global search
export function marketplaceFetchAutocompleteSuggestions({ q }) {
  const page = 1;
  const page_size = 6;
  const params = { q, page, page_size };

  const url = `${MARKETPLACE_BASE_URL}/repositories/search`;
  return {
    type: MARKETPLACE_FETCH_AUTOCOMPLETE_SUGGESTIONS,
    meta: params,
    payload: {
      promise:
        request
          .get(url)
          .accept('application/json')
          .set(bearer())
          .query(params)
          .end()
          .then((res) => res.body),
    },
  };
}

export function marketplaceFetchRepositorySummary({ id }) {
  const url =
    `${MARKETPLACE_BASE_URL}/repositories/${id}/summary`;
  return {
    type: MARKETPLACE_FETCH_REPOSITORY_SUMMARY,
    meta: { id },
    payload: {
      promise:
        request
          .get(url)
          .accept('application/json')
          .set(bearer())
          .end()
          .then((res) => res.body),
    },
  };
}

export function marketplaceFetchRepositoryDetail({ id }) {
  const url = `${MARKETPLACE_BASE_URL}/repositories/${id}`;
  return {
    type: MARKETPLACE_FETCH_REPOSITORY_DETAIL,
    meta: { id },
    payload: {
      promise:
        request
          .get(url)
          .accept('application/json')
          .set(bearer())
          .end()
          .then((res) => res.body),
    },
  };
}

export function marketplaceFetchBundleSummary({ id }) {
  const url = `${MARKETPLACE_BASE_URL}/bundles/${id}/summary`;
  return {
    type: MARKETPLACE_FETCH_BUNDLE_SUMMARY,
    meta: { id },
    payload: {
      promise:
        request
          .get(url)
          .accept('application/json')
          .set(bearer())
          .end()
          .then((res) => res.body),
    },
  };
}

export function marketplaceFetchBundleDetail({ id, shouldRedirectToLogin }) {
  const url = `${MARKETPLACE_BASE_URL}/bundles/${id}`;
  return {
    type: MARKETPLACE_FETCH_BUNDLE_DETAIL,
    meta: {
      id,
      shouldRedirectToLogin,
    },
    payload: {
      promise:
        request
          .get(url)
          .accept('application/json')
          .set(bearer())
          .end()
          .then((res) => res.body),
    },
  };
}

export function marketplaceFetchCategories() {
  const url = `${MARKETPLACE_BASE_URL}/categories`;
  return {
    type: MARKETPLACE_FETCH_CATEGORIES,
    payload: {
      promise:
        request
          .get(url)
          .accept('application/json')
          .set(bearer())
          .end()
          .then((res) => res.body),
    },
  };
}

// TODO Kristie 5/17/16 Use actual API when it is ready
export function marketplaceFetchPlatforms() {
  return {
    type: MARKETPLACE_FETCH_PLATFORMS,
    payload: {
      promise: new Promise((resolve) => {
        resolve(searchPlatformFilters);
      }),
    },
  };
}

// -----------------------------------------------------------------------------
// Home page special searches
// -----------------------------------------------------------------------------
// Most Popular Images (fetch 9)
export function marketplaceFetchMostPopular() {
  const params = { page_size: 9, source: OFFICIAL };
  const url = `${MARKETPLACE_BASE_URL}/repositories/search`;
  return {
    type: MARKETPLACE_FETCH_MOST_POPULAR,
    meta: params,
    payload: {
      promise:
        request
          .get(url)
          .accept('application/json')
          .set(bearer())
          .query(params)
          .end()
          .then((res) => res.body),
    },
  };
}

// Newest Images (fetch 9)
export function marketplaceFetchFeatured() {
  const params = {
    page_size: 9,
    category: 'featured',
    source: OFFICIAL,
  };
  const url = `${MARKETPLACE_BASE_URL}/repositories/search`;
  return {
    type: MARKETPLACE_FETCH_FEATURED,
    meta: params,
    payload: {
      promise:
        request
          .get(url)
          .accept('application/json')
          .set(bearer())
          .query(params)
          .end()
          .then((res) => res.body),
    },
  };
}

export function marketplaceCreateRepository({
  name,
  namespace,
  reponame,
  publisher,
  short_description,
  full_description,
  categories,
  platforms,
  source,
  logo_url,
  screenshots,
  links,
  eusa,
  download_attribute,
  instructions,
}) {
  const url = `${MARKETPLACE_PRIVATE_BASE_URL}/repositories/`;
  const body = {
    name,
    namespace,
    reponame,
    publisher,
    short_description,
    full_description,
    categories,
    platforms,
    source,
    logo_url,
    screenshots,
    links,
    eusa,
    download_attribute,
    instructions,
  };
  return {
    type: MARKETPLACE_CREATE_REPOSITORY,
    meta: body,
    payload: {
      promise:
        request
          .post(url)
          .set(bearer())
          .send(body)
          .accept('application/json')
          .type('application/json')
          .end()
          .then((res) => res.body),
    },
  };
}

export function marketplaceEditRepository({
  id,
  name,
  namespace,
  reponame,
  publisher,
  short_description,
  full_description,
  categories,
  platforms,
  source,
  logo_url,
  screenshots,
  links,
  eusa,
  download_attribute,
  instructions,
}) {
  const body = {
    name,
    namespace,
    reponame,
    publisher,
    short_description,
    full_description,
    categories,
    platforms,
    logo_url,
    source,
    screenshots,
    links,
    eusa,
    download_attribute,
    instructions,
  };
  const url = `${MARKETPLACE_PRIVATE_BASE_URL}/repositories/${id}`;
  return {
    type: MARKETPLACE_EDIT_REPOSITORY,
    meta: body,
    payload: {
      promise:
        request
          .patch(url)
          .set(bearer())
          .send(body)
          .accept('application/json')
          .type('application/json')
          .end()
          .then((res) => res.body),
    },
  };
}

export function marketplaceDeleteRepository({ id }) {
  const url = `${MARKETPLACE_PRIVATE_BASE_URL}/repositories/${id}`;
  return {
    type: MARKETPLACE_DELETE_REPOSITORY,
    payload: {
      promise:
        request
          .del(url)
          .set(bearer())
          .end()
          .then((res) => res.body),
    },
  };
}
