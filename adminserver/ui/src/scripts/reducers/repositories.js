'use strict';

import Immutable, { Map, List } from 'immutable';

import consts from 'consts';
import { RepositoryRecord } from 'records';

const defaultState = new Map({
  // Data

  // repositories stores a normalized result of repositories from listing all
  // repositories that a user can see.
  repositories: new Map(),

  // repositoriesByNamespace stores a map of namespace names to a normalized
  // list of repositories.
  //
  // We need this to list repositories for each namespace on the user/org page
  // list (repositories is replaced on each listRepositories call).
  //
  // Map of orgName => normalized repo data
  repositoriesByNamespace: new Map(),

  // Stores a list of repositories for each team within a namespace.
  // Map of orgName => teamName => normalized data
  repositoriesByTeam: new Map(),

  searchResults: new Map(),

  tags: new Map(),
  /* Repository tags
   * {
   *   org: {
   *     repo: {
   *       tagName: {
   *         'name': string,
   *         'inNotary': bool,
   *         'hasMismatch': bool, // if hashes from notary + registry mismatch
   *         'isTrusted': bool,
   *       },
   *       ...
   *     },
   *     ...
   *   },
   *   ...
   * }
   */

  userRepositoryAccess: new Map(),
  /* window.user's combined access levels to each repo
   * Takes into account his global admin status, org ownership, and repo admin from team or from user
   {
     userName: {
       namespace: {
         repoName: accessLevel,
         ...
       }
       ...
     },
     ...
   }
  */

  teamRepositoryAccess: new Map()
  /* Team access levels to each repository
  {
    orgName: {
      teamName: {
        repoName: accessLevel,
        ...

      },
      ...

    },
    ...

  }
  */
});

const actions = {
  [consts.repositories.UPDATE_SEARCH]: (state, action) => {
    if (!action.ready) {
      return state;
    }

    if (action.error) {
      return state.set('searchResults', new Map());
    }

    let repos = new Map().withMutations((mutableMap) => {
      action.payload.forEach((repoData) => {
        return mutableMap.setIn(
          [`${repoData.namespace}/${repoData.name}`, 'repo'],
          new RepositoryRecord(repoData)
        );
      });
    });

    return state.set('searchResults', repos);
  },

  [consts.repositories.LIST_ALL_REPOSITORIES]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    return state.set('repositories', new Map(action.payload.entities.repo));
  },

  [consts.repositories.LIST_SHARED_REPOSITORIES]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }

    // TODO: Remove duplication from above
    const { user, repositories } = action.payload;
    // Create a new map in the form of 'namespace/name' => RepositoryRacord
    let repos = new Map().withMutations(mutableMap => {
      repositories.forEach(repo => {
        mutableMap.set(`${repo.namespace}/${repo.name}`, new RepositoryRecord(repo));
      });
    });

    return state.setIn(['userRepositories', user, consts.repositories.SHOW_SHARED], repos);
  },

  [consts.repositories.LIST_REPOSITORIES]: (state, action) => {
    if (!action.ready) {
      return state;
    }

    if (action.error) {
      return state.set('repositories', new Map());
    }

    const { payload } = action;
    if (!payload.namespace) {
      // This is listing a user's repositories; this should be stored in the
      // plain 'repositories' map for the dashboard list.
      return state.setIn(['repositories'], payload.data);
    }

    return state.setIn(['repositoriesByNamespace', payload.namespace], payload.data);
  },

  [consts.repositories.LIST_REPO_USER_ACCESS]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }

    const { repository, userAccessList } = action.payload;
    return state.setIn(['repositories', repository.namespace, repository.name, 'access', 'users'], userAccessList);
  },

  [consts.repositories.GRANT_TEAM_ACCESS_TO_REPO]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }

    const { payload } = action;
    const orgName = payload.getIn(['result', 'org']);
    const teamName = payload.getIn(['result', 'team']);
    const repoName = payload.getIn(['result', 'repository']);

    return state.withMutations(s => {
      // Add the repository record to the entities.repo normalized data of
      // repositoriesByTeam
      s.setIn(['repositoriesByTeam', orgName, teamName, 'entities', 'repo', repoName], payload.getIn(['entities', 'repo', repoName]));
      // Add the repository name to the result list of arrays
      let result = s.getIn(['repositoriesByTeam', orgName, teamName, 'result'], new List());
      result = result.push(repoName);
      s.setIn(['repositoriesByTeam', orgName, teamName, 'result'], result);
    });
  },


  /**
   * Lists all repos that a team has access to.
   */
  [consts.repositories.LIST_TEAM_ACCESS_TO_REPO]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }

    const {
      orgName,
      teamName,
      data
    } = action.payload;

    return state.setIn(['repositoriesByTeam', orgName, teamName], data);
  },

  [consts.repositories.CREATE_REPO_AND_GRANT_TEAM_ACCESS]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }

    // Response data, normalised in action creator:
    //
    // {
    //   orgName: "name",
    //   teamName: "team",
    //   repositoryAccess: {
    //     accessLevel: "...",
    //     repository: {...}
    //   }
    // }

    // Append the repository to the list of team repositories
    const {
      orgName,
      repositoryAccess: {
        accessLevel,
        repository,
        repository: { name: repoName } // extract repo name
      },
      team: { name: teamName }
    } = action.payload;
    return state.merge({
      repositories: state.get('repositories').setIn([orgName, repoName, 'repo'], new RepositoryRecord(repository)),
      teamRepositoryAccess: state.get('teamRepositoryAccess').setIn([orgName, teamName, repoName], accessLevel)
    });
  },

  [consts.repositories.GET_TEAM_ACCESS_TO_REPO_NAMESPACE]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    const {
      accessLevel,
      namespace: {
        name: orgName
      },
      team: {
        name: teamName
      }
    } = action.payload;
    return state.setIn(['teamNamespaceAccess', orgName, teamName], accessLevel);
  },

  [consts.repositories.REVOKE_TEAM_ACCESS_TO_REPO_NAMESPACE]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    const {
      accessLevel,
      namespace: {
        name: orgName
      },
      team: {
        name: teamName
      }
    } = action.payload;
    return state.setIn(['teamNamespaceAccess', orgName, teamName], accessLevel);
  },

  // TODO: repeat of GRANT_TEAM_ACCESS
  [consts.repositories.CHANGE_TEAM_ACCESS_TO_REPO]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    // Append the repository to the list of team repositories
    const {
      orgName,
      teamName,
      teamId,
      repoName,
      accessLevel
    } = action.payload;

    return state
      .setIn(['repositoriesByTeam', orgName, teamName, 'entities', 'repo', `${orgName}/${repoName}`, 'accessLevel'], accessLevel)
      .setIn(['repositories', orgName, repoName, 'access', 'teams', teamId, 'accessLevel'], accessLevel);
  },

  [consts.repositories.REVOKE_TEAM_ACCESS_TO_REPO]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    const {
      orgName,
      teamName,
      teamId,
      repoName
    } = action.payload;

    // TODO: Also remove the ID from the result array of the normalized response
    return state
      .deleteIn(['repositoriesByTeam', orgName, teamName, 'entities', 'repo', `${orgName}/${repoName}`])
      .deleteIn(['repositories', orgName, repoName, 'access', 'teams', teamId]);
  },

  [consts.repositories.LIST_REPO_TEAM_ACCESS]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }

    const {
      repository,
      teamAccessList
    } = action.payload;

    return state.setIn(['repositories', repository.namespace, repository.name, 'access', 'teams'], teamAccessList);
  },

  [consts.repositories.CREATE_REPOSITORY]: (state, action) => {
    if (!action.ready || action.error) { return state; }

    return state
      .mergeIn(['repositories', 'entities', 'repo'], action.payload.getIn(['entities', 'repo']))
      .updateIn(['repositories', 'result'], r => r ? r.push(action.payload.get('result', new List())) : undefined);
  },

/** Repo specific actions **/

  [consts.repositories.GET_REPOSITORY]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }

    return state.setIn(['repositories', action.payload.namespace, action.payload.name, 'repo'], new RepositoryRecord(action.payload));
  },

  [consts.repositories.GET_REPOSITORY_WITH_USER_PERMISSIONS]: (state, action) => {
    // TODO: Copy into user reducer and upsert user information in state
    if (!action.ready || action.error) {
      return state;
    }

    // action.payload = {accessLevel: 'string', user: {...}, repository: {...}};
    const {
      user: { name: userName },
      repository: { namespace, name: repoName },
      repository,
      accessLevel
    } = action.payload;

    return state.withMutations( map => {
      return map
        .setIn(['repositories', namespace, repoName, 'repo'], new RepositoryRecord(repository))
        .setIn(['userRepositoryAccess', userName, namespace, repoName], accessLevel);
    });
  },

  [consts.repositories.UPDATE_REPOSITORY]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    return state
        .setIn(['repositories', action.payload.namespace, action.payload.name, 'repo'], new RepositoryRecord(action.payload));
  },

  [consts.repositories.DELETE_REPOSITORY]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    return state.deleteIn(['repositories', action.params.namespace, action.params.repo]);
  },

  /**
   * listing tags for a repository
   */
  [consts.repositories.GET_REPOSITORY_TAGS]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }

    const { namespace, name: repoName } = action.payload;

    // Add all tags to the immutable map keyed by their name
    const tags = new Map().withMutations( map => {
      action.payload.tags.forEach(tag => map.set(tag.name, Immutable.fromJS(tag)));
    });

    return state.setIn(['tags', namespace, repoName], tags);
  },

  [consts.repositories.DELETE_REPO_MANIFEST]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    const { namespace, repo: repoName, tag } = action.payload;

    // Delete tag from state
    return state.deleteIn(['tags', namespace, repoName, tag]);
  }
};

export default function repositoriesStore(state = defaultState, payload) {
  if (typeof actions[payload.type] === 'function') {
    return actions[payload.type](state, payload);
  }
  return state;
}
