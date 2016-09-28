'use strict';

import { createSelector } from 'reselect';

/**
 *
 * @return array  Array of NamespaceRecords in the order of the search results
 */
export const getNamespacesFromSearch = (state) => {
  const search = state.namespaces.get('search');
  const results = search.getIn(['result', 'orgs'], []);
  return results.map(id => {
    return search.getIn(['entities', 'org', id]);
  });
};

export const getNamespaceNamesFromSearch = createSelector(
  [getNamespacesFromSearch],
  (namespaces) => namespaces.map(n => n.name)
);
