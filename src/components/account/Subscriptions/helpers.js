import map from 'lodash/map';

export const setIsOwnerParam = (fetchedNamespaces, ownedNamespaces) => {
  const results = {};
  map(fetchedNamespaces, (userOrg) => {
    const username = userOrg.username || userOrg.orgname;
    const isOwner = ownedNamespaces.indexOf(username) >= 0;
    results[username] = { ...userOrg, isOwner };
  });
  return results;
};
