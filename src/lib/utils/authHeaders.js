import cookie from 'react-cookie';
import createBearer from 'lib/utils/create-bearer';
import createJWT from 'lib/utils/create-jwt';

const defaultHeaders = () => {
  const headers = {
    Accept: 'application/json',
  };

  const csrfToken = cookie.load('csrftoken');
  if (csrfToken) {
    headers['X-CSRFToken'] = csrfToken;
  } else if (process.env && process.env.DOCKERSTORE_CSRF) {
    headers['X-CSRFToken'] = process.env.DOCKERSTORE_CSRF;
  }

  return headers;
};

export default {
  jwt: () => {
    const headers = defaultHeaders();
    if (process.env.DOCKERSTORE_TOKEN) {
      headers.Authorization = createJWT(process.env.DOCKERSTORE_TOKEN);
    }

    return headers;
  },
  bearer: () => {
    const headers = defaultHeaders();
    if (process.env.DOCKERSTORE_TOKEN) {
      headers.Authorization = createBearer(process.env.DOCKERSTORE_TOKEN);
    }

    return headers;
  },
};
