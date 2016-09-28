'use strict';

import { Index } from 'dtr-js-sdk';
import consts from 'consts';

import { normalizeSearchResults } from 'normalizers';

export function searchAll(term) {
  return {
    type: consts.SEARCH,
    meta: {
      promiseScope: [], // Don't care about previous searchTerms
      promise: Index.autocomplete({
          query: term,
          includeRepositories: true,
          includeAccounts: true,
          limit: 100
        })
        .then(resp => {
          return normalizeSearchResults(resp.data);
        })
    }
  };
}

export function clearSearch() {
  return {
    type: consts.SEARCH,
    ready: true,
    payload: {}
  };
}
