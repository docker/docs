const request = require('superagent-promise')(require('superagent'), Promise);
import { bearer } from 'lib/utils/authHeaders';

/* eslint-disable max-len */
export const PUBLISH_ACCEPT_VENDOR_AGREEMENT = 'PUBLISH_ACCEPT_VENDOR_AGREEMENT';
export const PUBLISH_ADD_PRODUCT_TIERS = 'PUBLISH_ADD_PRODUCT_TIERS';
export const PUBLISH_CREATE_PRODUCT = 'PUBLISH_CREATE_PRODUCT';
export const PUBLISH_DELETE_PRODUCT_TIERS = 'PUBLISH_DELETE_PRODUCT_TIERS';
export const PUBLISH_FETCH_PRODUCT_DETAILS = 'PUBLISH_FETCH_PRODUCT_DETAILS';
export const PUBLISH_FETCH_PRODUCT_LIST = 'PUBLISH_FETCH_PRODUCT_LIST';
export const PUBLISH_FETCH_PRODUCT_TIERS = 'PUBLISH_FETCH_PRODUCT_TIERS';
export const PUBLISH_GET_PUBLISHERS = 'PUBLISH_GET_PUBLISHERS';
export const PUBLISH_GET_SIGNUP = 'PUBLISH_GET_SIGNUP';
export const PUBLISH_GET_VENDOR_AGREEMENT = 'PUBLISH_GET_VENDOR_AGREEMENT';
export const PUBLISH_SUBSCRIBE = 'PUBLISH_SUBSCRIBE';
export const PUBLISH_UPDATE_PRODUCT_DETAILS = 'PUBLISH_UPDATE_PRODUCT_DETAILS';
export const PUBLISH_UPDATE_PRODUCT_REPOS = 'PUBLISH_UPDATE_PRODUCT_REPOS';
export const PUBLISH_UPDATE_PRODUCT_TIERS = 'PUBLISH_UPDATE_PRODUCT_TIERS';
export const PUBLISH_UPDATE_PUBLISHER_INFO = 'PUBLISH_UPDATE_PUBLISHER_INFO';

const PUBLISH_HOST = '';
export const PUBLISH_BASE_URL = `${PUBLISH_HOST}/api/publish/v1`;

export const publishSubscribe = ({
  first_name,
  last_name,
  company,
  phone_number,
  email,
}) => {
  const url = `${PUBLISH_BASE_URL}/signups`;
  return {
    type: PUBLISH_SUBSCRIBE,
    payload: {
      promise:
        request
          .post(url)
          .set(bearer())
          .send({
            first_name,
            last_name,
            company,
            phone_number,
            email,
          })
          .end()
          .then(res => res.body),
    },
  };
};

export const publishGetSignup = () => {
  const url = `${PUBLISH_BASE_URL}/signups`;
  return {
    type: PUBLISH_GET_SIGNUP,
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
};

export const publishGetPublishers = () => {
  const url = `${PUBLISH_BASE_URL}/publishers`;
  return {
    type: PUBLISH_GET_PUBLISHERS,
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
};

// NOTE: no reducers for this action yet
export const publishUpdatePublisherInfo = (data) => {
  const url = `${PUBLISH_BASE_URL}/publishers`;
  /*
    data:
    {
      "email": "nandhini@docker.com",
      "first_name": "Nandhini",
      "last_name": "Santhanam",
      "company": "Docker,Inc",
      "phone_number": "2222222222",
      "links": [{
          "name": "google.com",
          "label": "website"
      }]
    }
  */
  return {
    type: PUBLISH_UPDATE_PUBLISHER_INFO,
    payload: {
      promise:
        request
          .patch(url)
          .accept('application/json')
          .set(bearer())
          .send(data)
          .end()
          .then((res) => res.body),
    },
  };
};

export const publishFetchProductList = () => {
  const url = `${PUBLISH_BASE_URL}/products`;
  return {
    type: PUBLISH_FETCH_PRODUCT_LIST,
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
};

export const publishFetchProductDetails = ({ product_id }) => {
  const url = `${PUBLISH_BASE_URL}/products/${product_id}`;
  return {
    type: PUBLISH_FETCH_PRODUCT_DETAILS,
    payload: {
      promise:
        request
          .get(url)
          .set(bearer())
          .accept('application/json')
          .end()
          .then((res) => res.body),
    },
  };
};

export const publishGetVendorAgreement = () => {
  const url = `${PUBLISH_BASE_URL}/vendor-agreement`;
  return {
    type: PUBLISH_GET_VENDOR_AGREEMENT,
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
};

// NOTE: no reducers for this action yet
export const publishCreateProduct = ({ name, status, repositories }) => {
  const body = {
    name,
    status,
    repositories,
  };
  const url = `${PUBLISH_BASE_URL}/products`;
  return {
    type: PUBLISH_CREATE_PRODUCT,
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
};

// NOTE: no reducers for this action yet
export const publishUpdateProductRepos = ({ product_id, repoSources }) => {
  const url = `${PUBLISH_BASE_URL}/products/${product_id}/repositories`;
  return {
    type: PUBLISH_UPDATE_PRODUCT_REPOS,
    payload: {
      promise:
        request
          .put(url)
          .set(bearer())
          .send(repoSources)
          .accept('application/json')
          .type('application/json')
          .end()
          .then((res) => res.body),
    },
  };
};

// NOTE: no reducers for this action yet
export const publishUpdateProductDetails = ({ product_id, details }) => {
  const url = `${PUBLISH_BASE_URL}/products/${product_id}/`;
  /*
    details:
    {
      name, * required
      status, * required
      product_type,
      full_description,
      short_description,
      categories,
      platforms,
      links,
      screenshots,
    }
  */
  return {
    type: PUBLISH_UPDATE_PRODUCT_DETAILS,
    payload: {
      promise:
        request
          .patch(url)
          .set(bearer())
          .send(details)
          .accept('application/json')
          .type('application/json')
          .end()
          .then((res) => res.body),
    },
  };
};

export const publishAcceptVendorAgreement = () => {
  const url = `${PUBLISH_BASE_URL}/accept-vendor-agreement`;
  return {
    type: PUBLISH_ACCEPT_VENDOR_AGREEMENT,
    payload: {
      promise:
        request
          .post(url)
          .set(bearer())
          .end()
          .then((res) => res.body),
    },
  };
};

export const publishFetchProductTiers = ({ product_id }) => {
  const url = `${PUBLISH_BASE_URL}/products/${product_id}/rate-plans`;
  return {
    type: PUBLISH_FETCH_PRODUCT_TIERS,
    payload: {
      promise:
        request
          .get(url)
          .set(bearer())
          .end()
          .then((res) => res.body),
    },
  };
};

export const publishCreateProductTiers = ({ product_id, tiersList }) => {
  const url = `${PUBLISH_BASE_URL}/products/${product_id}/rate-plans`;
  return {
    type: PUBLISH_ADD_PRODUCT_TIERS,
    payload: {
      promise:
        request
          .post(url)
          .set(bearer())
          .send(tiersList)
          .end()
          .then((res) => res.body),
    },
  };
};

export const publishUpdateProductTiers = ({ product_id, tiersObject }) => {
  const url = `${PUBLISH_BASE_URL}/products/${product_id}/rate-plans`;
  return {
    type: PUBLISH_UPDATE_PRODUCT_TIERS,
    payload: {
      promise:
        request
          .put(url)
          .set(bearer())
          .send(tiersObject)
          .end()
          .then((res) => res.body),
    },
  };
};

export const publishDeleteProductTiers = ({ product_id, tier_id }) => {
  const url = `${PUBLISH_BASE_URL}/products/${product_id}/rate-plans/${tier_id}`;
  return {
    type: PUBLISH_DELETE_PRODUCT_TIERS,
    payload: {
      promise:
        request
          .del(url)
          .set(bearer())
          .end()
          .then((res) => res.body),
    },
  };
};
