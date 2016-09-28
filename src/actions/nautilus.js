const request = require('superagent-promise')(require('superagent'), Promise);
import { normalize } from 'normalizr';
import { jwt } from 'lib/utils/authHeaders';
import { normalizers } from 'lib/constants/nautilus';
import { DEFAULT_TAGS_PAGE_SIZE } from 'lib/constants/defaults';

export const NAUTILUS_FETCH_TAGS_AND_SCANS = 'NAUTILUS_FETCH_TAGS_AND_SCANS';
export const NAUTILUS_FETCH_SCAN_DETAIL = 'NAUTILUS_FETCH_SCAN_DETAIL';

// Don't need to proxy the nautilus api
const NAUTILUS_HOST = '';
const NAUTILUS_BASE_URL = `${NAUTILUS_HOST}/api/nautilus/v1/repositories`;

const DOCKERHUB_HOST = '';
const DOCKERHUB_BASE_URL = `${DOCKERHUB_HOST}/v2/repositories`;


export const nautilusFetchTagsAndScans = ({
  id,
  namespace,
  reponame,
  page = 1,
  page_size = DEFAULT_TAGS_PAGE_SIZE,
}) => {
  const scansUrl = `${NAUTILUS_BASE_URL}/summaries/${namespace}/${reponame}/`;
  const tagsUrl = `${DOCKERHUB_BASE_URL}/${namespace}/${reponame}/tags/`;
  return {
    type: NAUTILUS_FETCH_TAGS_AND_SCANS,
    meta: { id, namespace, reponame, page, page_size },
    payload: {
      promise: new Promise((resolve, reject) => {
        // Fetch all of the tags for this repository (paginated)
        request.get(tagsUrl)
          .set(jwt())
          .query({ page, page_size })
          .end()
          .then(
            // onSuccess of fetch tags, attempt to fetch the scans
            ({ body: tags }) => {
              // Fetch all of the scans for this repository
              request.get(scansUrl)
                .set(jwt())
                .end()
                .then(
                  // onSuccess of fetch scans, return both tags and scans
                  ({ body: scans }) => { resolve([tags, scans]); },
                  // onError of fetch scans, only return tags
                  () => { resolve([tags]); },
                );
            },
            // onError of fetch tags, reject the promise (error)
            reject
          );
      }),
    },
  };
};

export const nautilusFetchScanDetail = ({ id, namespace, reponame, tag }) => {
  const url = `${NAUTILUS_BASE_URL}/result`;
  const params = { detailed: 1, namespace, reponame, tag };
  return {
    type: NAUTILUS_FETCH_SCAN_DETAIL,
    meta: { ...params, id },
    payload: {
      promise:
        request
          .get(url)
          .set(jwt())
          .query(params)
          .end()
          .then(({ body }) => {
            // The API response contains a 'scan' resource within an object
            // inside 'scan_details'
            const { image, latest_scan_status, scan_details } = body;
            // Normalize the API result using the 'scan' normalizer schema
            // TODO Kristie 5/11/16 Do we need the reponame and tag?
            const normalized = normalize({
              latest_scan_status,
              reponame: image.reponame,
              tag: image.tag,
              ...scan_details,
            }, normalizers.scan);
            return normalized;
          }),
    },
  };
};
