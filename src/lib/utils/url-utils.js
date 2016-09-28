import { flow } from 'lodash';
import routes from 'lib/constants/routes';

// TODO: This module was copy/pasted from apiUtils and needs to
// be more generic now.

const regexp = /\/?api\/v\d+\//;

const removePrexif = (u) => u.replace(regexp, '');
const removeSuffix = (u) => u.slice(0, -1);

const normalize = flow(removeSuffix, removePrexif);

export const tagName = (u) => normalize(u).split('/tag/').pop();

export const removeTrailingSlash = (u) => {
  return u && u.charAt(u.length - 1) === '/' ? removeSuffix(u) : u;
};

export const getQueryVariable = (qs, variable) => {
  const vars = qs.substring(1).split('&');
  for (let i = 0; i < vars.length; i++) {
    const [key, val] = vars[i].split('=');
    if (decodeURIComponent(key) === variable) {
      return decodeURIComponent(val);
    }
  }

  return null;
};

export const getLoginRedirectURL = () => {
  const base = encodeURIComponent(window && window.location.pathname);
  const query = window && window.location.search;
  const nextUri = `${base}${query || ''}`;
  return `${routes.login()}?next=${nextUri}`;
};

export const isValidUrl = url => {
  // Let localhost pass
  if (/localhost/.test(url)) {
    return true;
  }
  // Matches a simple protocol://pathname with optional query string
  // eslint-disable-next-line max-len
  return /(http|ftp|https):\/\/[\w-]+(\.[\w-]+)+([\w.,@?^=%&amp;:\/~+#-]*[\w@?^=%&amp;\/~+#-])?/
    .test(url);
};

export const isQsPathSecure = (path = '') => {
  return (
    path &&
    !/^\w+:/.test(path) &&
    path.charAt(0) === '/' &&
    path.charAt(1) !== '/'
  );
};
