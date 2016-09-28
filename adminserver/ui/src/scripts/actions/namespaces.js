'use strict';

import { Index } from 'dtr-js-sdk';
import consts from 'consts';
import { normalizeSearchResults } from 'normalizers';
import { normalizedToRecords } from 'utils/records';
import { NamespaceRecord } from 'records';

export function searchNamespaces (query = '') {
  return {
    type: consts.SEARCH_NAMESPACES,
    meta: {
      promiseScope: [],
      promise: Index.autocomplete({
        query: query,
        includeRepositories: false,
        includeAccounts: true
      })
      .then(resp => {
        const data = normalizeSearchResults(resp.data);
        return normalizedToRecords(data, {user: NamespaceRecord, org: NamespaceRecord});
      })
    }
  };
}
