'use strict';

import officialLogos from 'common/data/officialLogos.js';
const officialImage = '/public/images/logos/mini-logo-white-inset.png';

const isOfficialNamespace = (namespace) => namespace === '_' || namespace === 'library';

// reponame is an optional param used to descriminate between official repos
// and it will be ignored unless it is an official repo
export function mkAvatarForNamespace(namespace, reponame) {
  if (isOfficialNamespace(namespace) && reponame && officialLogos[reponame]) {
    return `/public/images/official/${officialLogos[reponame]}`;
  }
  if(isOfficialNamespace(namespace) || !namespace) {
    return officialImage;
  }
  return `${process.env.HUB_API_BASE_URL}/v2/users/${namespace}/avatar/`;
}

/**
 * This function should be treated as deprecated. We can not reliably
 * detect the namespace type (User || Org) from the urls.
 */
export function mkAvatarForOrg(namespace) {
    return `${process.env.HUB_API_BASE_URL}/v2/orgs/${namespace}/avatar/`;
}

export function isOfficialAvatarURL(url) {
  return url === officialImage;
}
