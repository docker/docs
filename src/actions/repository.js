const request = require('superagent-promise')(require('superagent'), Promise);
import {
  DEFAULT_TAGS_PAGE_SIZE,
  DEFAULT_COMMENTS_PAGE_SIZE,
} from 'lib/constants/defaults';
import { jwt } from 'lib/utils/authHeaders';

/* eslint-disable max-len */
export const REPOSITORY_FETCH_COMMENTS = 'REPOSITORY_FETCH_COMMENTS';
export const REPOSITORY_FETCH_OWNED_NAMESPACES = 'REPOSITORY_FETCH_OWNED_NAMESPACES';
export const REPOSITORY_FETCH_REPOSITORIES_FOR_NAMESPACE = 'REPOSITORY_FETCH_REPOSITORIES_FOR_NAMESPACE';
export const REPOSITORY_FETCH_IMAGE_DETAIL = 'REPOSITORY_FETCH_IMAGE_DETAIL';
export const REPOSITORY_FETCH_IMAGE_TAGS = 'REPOSITORY_FETCH_IMAGE_TAGS';
/* eslint-enable max-len */

const DOCKERHUB_HOST = '';
const DOCKERHUB_BASE_URL = `${DOCKERHUB_HOST}/v2/repositories`;

export function repositoryFetchImageDetail({
  namespace,
  reponame,
}) {
  const url = `${DOCKERHUB_BASE_URL}/${namespace}/${reponame}/`;
  return ({
    type: REPOSITORY_FETCH_IMAGE_DETAIL,
    meta: { namespace, reponame },
    payload: {
      promise:
        request
          .get(url)
          .accept('application/json')
          .end()
          .then((res) => res.body),
    },
  });
}

export function repositoryFetchImageTags({
  // Temporary until the marketplace tags service is up
  id,
  isCertified = false,
  namespace,
  page = 1,
  page_size = DEFAULT_TAGS_PAGE_SIZE,
  reponame,
}) {
  const url = `${DOCKERHUB_BASE_URL}/${namespace}/${reponame}/tags/`;
  const params = { page_size, page };
  return ({
    type: REPOSITORY_FETCH_IMAGE_TAGS,
    meta: { id, isCertified, namespace, reponame, page, page_size },
    payload: {
      promise:
        request
          .get(url)
          .accept('application/json')
          .query(params)
          .end()
          .then((res) => res.body),
    },
  });
}


export function repositoryFetchComments({
  id,
  isCertified = false,
  namespace,
  reponame,
  page = 1,
  page_size = DEFAULT_COMMENTS_PAGE_SIZE,
}) {
  const url = `${DOCKERHUB_BASE_URL}/${namespace}/${reponame}/comments/`;
  const params = { page_size, page };
  return ({
    type: REPOSITORY_FETCH_COMMENTS,
    meta: { id, isCertified, page, page_size, namespace, reponame },
    payload: {
      promise:
        request
          .get(url)
          .accept('application/json')
          .query(params)
          .end()
          .then((res) => res.body),
    },
  });
}

export function repositoryFetchOwnedNamespaces() {
  const url = `${DOCKERHUB_BASE_URL}/namespaces/`;
  return {
    type: REPOSITORY_FETCH_OWNED_NAMESPACES,
    payload: {
      promise:
        request
          .get(url)
          .set(jwt())
          .end()
          .then((res) => res.body),
    },
  };
}

export function repositoryFetchRepositoriesForNamespace({ namespace }) {
  const url = `${DOCKERHUB_BASE_URL}/${namespace}/`;
  return {
    type: REPOSITORY_FETCH_REPOSITORIES_FOR_NAMESPACE,
    meta: { namespace },
    payload: {
      promise:
        request
          .get(url)
          .set(jwt())
          .query({ page_size: 0 })
          .end()
          .then((res) => res.body),
    },
  };
}
