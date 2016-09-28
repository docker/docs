const request = require('superagent-promise')(require('superagent'), Promise);
import { bearer } from 'lib/utils/authHeaders';

/* eslint-disable max-len */
export const WHITELIST_FETCH_AUTHORIZATION = 'WHITELIST_FETCH_AUTHORIZATION';
export const WHITELIST_SUBSCRIBE_TO_BETA = 'WHITELIST_SUBSCRIBE_TO_BETA';
export const WHITELIST_AM_I_WAITING = 'WHITELIST_AM_I_WAITING';

const WHITELIST_HOST = '';
export const WHITELIST_BASE_URL = `${WHITELIST_HOST}/api/whitelist/v1`;

//------------------------------------------------------------------------------
// AUTHORIZATION
//------------------------------------------------------------------------------

export const whitelistFetchAuthorization = () => {
  const url = `${WHITELIST_BASE_URL}/authorize`;
  return {
    type: WHITELIST_FETCH_AUTHORIZATION,
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

// http://store-stage.docker.com/api/whitelist/v1/consumer_interest
export const whitelistAmIWaiting = () => {
  const url = `${WHITELIST_BASE_URL}/consumer_interest`;
  return {
    type: WHITELIST_AM_I_WAITING,
    payload: {
      promise:
        request
          .get(url)
          .set(bearer())
          .end()
          .then(res => res.body),
    },
  };
};

export const whitelistSubscribeToBeta = ({
  first_name,
  last_name,
  company,
  email,
}) => {
  const url = `${WHITELIST_BASE_URL}/consumer_interest`;
  return {
    type: WHITELIST_SUBSCRIBE_TO_BETA,
    payload: {
      promise:
        request
          .post(url)
          .set(bearer())
          .send({
            first_name,
            last_name,
            company,
            email,
          })
          .end()
          .then(res => res.body),
    },
  };
};
