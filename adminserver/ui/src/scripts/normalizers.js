'use strict';

import {
  normalize,
  Schema,
  arrayOf
} from 'normalizr';

/**
 * Basic schema types
 **/

export const RepoSchema = new Schema('repo', { idAttribute: data => {
  // If this is in a search result the ID will be nested in a repository object
  if (typeof data.repository === 'object') {
    return `${data.repository.namespace}/${data.repository.name}`;
  }
  return `${data.namespace}/${data.name}`;
}});

// These should be keyed by name as each repo contains only the name (and not ID)
export const OrgSchema = new Schema('org', { idAttribute: data => {
  // If this is in a search result the ID will be nested in an account object
  if (typeof data.account === 'object') {
    return data.account.name;
  }
  return data.name;
}});

export const UserSchema = new Schema('user', { idAttribute: data => {
  // If this is in a search result the ID will be nested in an account object
  if (typeof data.account === 'object') {
    return data.account.name;
  }
  if (typeof data.member === 'object') {
    return data.member.name;
  }
  return data.name;
}});

export const TeamSchema = new Schema('team', { idAttribute: 'name' });

/**
 * Nested schema declarations
 **/

// TODO: Namespaces inside repositories

// normalizeSearchResults normalizes search results into repos, users and orgs.
//
// The result is:
// {
//    results: [{id: 'id1', schema: 'repo'}, {id: 'id2', schema: 'user'}, ...],
//    entities: {
//      repo: {
//        id1: {...}
//      }
//      user: {
//        id2: {...}
//      }
//    }
// }
//
// This makes it easy to list repos, users, and orgs separately in a search
// result.
//
// To show search results in their existing order iterate through results and
// pull out the correct resource from the schema and ID listed.
//
// This is a curried function: you only need to call this function with
// resp.data.results to make things work.
export const normalizeSearchResults = (data) => {
  const fromSchema = {
    repositoryResults: arrayOf(RepoSchema),
    accountResults: arrayOf(
      {
        users: UserSchema,
        orgs: OrgSchema
      },
      {
        schemaAttribute: (item) => {
          if (item.isOrg === true) {
            return 'orgs';
          }
          return 'users';
        }
      }
    )
  };
  return normalize(data, fromSchema, {
    // by default normalize will fuck up our results object, producing:
    // results: {
    //   repositoryResults: [...],
    //   accountResults: [...]
    // }
    //
    // We don't want this, we want repos, users, and orgs in results separately
    // - the same as the entities produced.
    //
    // This function does that for us.
    assignEntity: (normalized, key, entity) => {
      if (key === 'repositoryResults') {
        normalized.repos = entity;
        return;
      }
      if (key === 'accountResults') {
        entity.forEach(item => {
          const { id, schema } = item;
          if (Array.isArray(normalized[schema])) {
            normalized[schema].push(id);
          } else {
            normalized[schema] = [id];
          }
        });
        return;
      }
      normalized[key] = entity;
    }
  });
};

