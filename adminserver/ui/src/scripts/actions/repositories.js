'use strict';
import { Repositories } from 'dtr-js-sdk';
import consts from 'consts';
import { accessLevelCompare } from 'utils';
import { RepoSchema, OrgSchema, TeamSchema } from 'normalizers';
import { normalize, arrayOf } from 'normalizr';
import { normalizedToRecords } from 'utils/records';
import { RepositoryRecord, TeamRecord, OrganizationRecord } from 'records';
import { fromJS } from 'immutable';

const { repositories } = consts;

const DEFAULT_LIMIT = 1000;

const Repos = Repositories.default;

// Even though this can only be accessed by the current user you must always
// pass in the username for our reducer.
export function listAllRepositories(user) {
  return {
    type: consts.repositories.LIST_ALL_REPOSITORIES,
    meta: {
      promiseScope: [user],
      promise: Repos
        .listAllRepositories({}, {limit: DEFAULT_LIMIT})
        .then(resp => normalizedToRecords(
          normalize(resp.data.repositories, arrayOf(RepoSchema)),
          {'repo': RepositoryRecord}
        ))
    }
  };
}

/**
 * listRepositories lists repositories for a given namespace.
 *
 * NOTE: This returns an object containing the namespace and `data` which
 * represents the normalized response data. This is so that we can save
 * normalized results to a map of namespaces within the reducer.
 */
export function listRepositories({ namespace, limit = DEFAULT_LIMIT }) {
  return {
    type: consts.repositories.LIST_REPOSITORIES,
    meta: {
      promiseScope: [namespace],
      promise: Repos
        .listRepositories({namespace}, {limit})
        .then(resp => ({
            namespace,
            data: normalizedToRecords(
              normalize(resp.data.repositories, arrayOf(RepoSchema)),
              {'repo': RepositoryRecord}
            )
          })
        )
    }
  };
}

/**
 * Lists which repositories a team has access to.
 */
export function listTeamRepositories({ orgName, teamName }) {
  return {
    type: consts.repositories.LIST_TEAM_ACCESS_TO_REPO,
    meta: {
      promiseScope: [orgName, teamName],
      promise: Repos.listTeamRepoAccess({user: orgName, team: teamName}).then(resp => {
        return {
          // Add the organization name from our request, as the organization ID
          // sucks, is unusable, and is unfortunately the only thing included in
          // the API response
          orgName,
          teamName,
          data: normalizedToRecords(
            normalize(
              resp.data.repositoryAccessList,
              arrayOf(RepoSchema),
              {
                assignEntity: (obj, key, val) => {
                  if (typeof val === 'object') {
                    Object.keys(val).forEach(i => obj[i] = val[i]);
                  } else {
                    obj[key] = val;
                  }
                }
              }),
            {'repo': RepositoryRecord}
          )
        };
      })
    }
  };
}

/**
 * Returns an org or team's repositories depending on the teamName - if the
 * teamName is undefined this will return the org repos
 */
export const listOrgOrTeamRepos = ({ orgName, teamName, limit = DEFAULT_LIMIT }) => {
  if (teamName === undefined || teamName === '') {
    return listRepositories({namespace: orgName, limit});
  } else {
    return listTeamRepositories({orgName, teamName, limit});
  }
};

export const grantTeamAccessToRepo = ({ orgName, teamName, repo, accessLevel}) => ({
  type: repositories.GRANT_TEAM_ACCESS_TO_REPO,
  meta: {
    promiseScope: [orgName, teamName, repo],
    promise: Repos
      .grantRepoTeamAccess({namespace: orgName, repo, team: teamName}, {accessLevel})
      .then(response => {
        // Add the org to the data for normalizr
        response.data.org = {name: orgName};
        response.data.team = teamName;
        // Quick hack to add accessLevel to the repository
        response.data.repository.accessLevel = response.data.accessLevel;
        // Return a normalized response of the team and repository
        return normalizedToRecords(
          normalize(response.data, {
            team: TeamSchema,
            repository: RepoSchema,
            org: OrgSchema
          }),
          {
            team: TeamRecord,
            repo: RepositoryRecord,
            org: OrganizationRecord
          }
        );
      })
  }
});

export function createRepository({ namespace, name, visibility, shortDescription }) {
  return (dispatch) => {
    dispatch({
      type: repositories.CREATE_REPOSITORY,
      meta: {
        promise: Repos.createRepository({namespace}, {name, visibility, shortDescription})
          .then(resp => {
            dispatch(listOrgOrTeamRepos({
              orgName: namespace,
              teamName: ''
            }));
            return normalizedToRecords(normalize(resp.data, RepoSchema), {repo: RepositoryRecord});
          })
      }
    });
  };
}


/**
 * Creates a repository then adds team access to the repository.
 *
 */
export function createRepoAndGrantTeamAccess({ namespace, teamName, repo, visibility, shortDescription, accessLevel }) {
  return (dispatch) => {

    dispatch({
      type: repositories.CREATE_REPOSITORY,
      meta: {
        promise: Repos.createRepository({ namespace }, { name: repo, visibility, shortDescription })
          .then(resp => {
            dispatch(listOrgOrTeamRepos({
              orgName: namespace,
              teamName: ''
            }));

            if (teamName) {
              dispatch(grantTeamAccessToRepo({
                orgName: namespace,
                teamName,
                repo,
                accessLevel
              }));
            }

            return normalizedToRecords(normalize(resp.data, RepoSchema), {repo: RepositoryRecord});
          })
      }
    });

  };
}


// TODO: Are the below used?

export function listRepoUserAccess(params) {
  return {
    type: consts.repositories.LIST_REPO_USER_ACCESS,
    meta: {
      promiseScope: [params.namespace, params.repo],
      promise: Repos.listRepoUserAccess(params).then((resp) => resp.data)
    }
  };
}

export function listRepoTeamAccess(params) {
  return {
    type: consts.repositories.LIST_REPO_TEAM_ACCESS,
    meta: {
      promiseScope: [params.namespace, params.repo],
      promise: Promise.all([
        Repos.listRepoTeamAccess(params),
        Repos.listRepoNamespaceTeamAccess(params)
      ]).then((responses) => {
        // Need to merge:
        // - teams who have access to the namespace this repo is in
        // - teams who have access to this repo specifically

        let accumulation = {};
        responses.map((resp) => {
          resp.data.teamAccessList.reduce((accum, teamData) => {
            const {
              accessLevel,
              team
              } = teamData;
            const existingTeam = accum[team.id];
            if (!existingTeam ||
              accessLevelCompare(accessLevel, existingTeam.accessLevel) > 0) {
              accum[team.id] = teamData;
            }
            return accum;
          }, accumulation);
        });
        return {
          repository: responses[0].data.repository,
          teamAccessList: fromJS(accumulation)
        };
      })
    }
  };
}

/** Repo specific actions **/

// Unless you have a good reason to use this endpoint you should *always* prefer
// getUserRepoAccess below.
export function getRepository(params) {
  return {
    type: consts.repositories.GET_REPOSITORY,
    meta: {
      promiseScope: [params.namespace, params.repo],
      promise: Repos.getRepository(params).then((resp) => resp.data)
    }
  };
}

/**
 * getUserRepoAccess fetches the repository resource, the user's access level to
 * the repository and user information in one call.
 *
 * @param string user  username to fetch access level for
 * @param string namespace  namespace of repo
 * @param string repo  repo name
 */
export function getUserRepoAccess({ user, namespace, repo }) {
  return {
    type: consts.repositories.GET_REPOSITORY_WITH_USER_PERMISSIONS,
    meta: {
      promiseScope: [namespace, repo, user],
      promise: Repos.getUserRepoAccess({user, namespace, repo})
        .then((resp) => resp.data)
    }
  };
}

export function updateRepository({ namespace, repo }, data) {
  return {
    type: consts.repositories.UPDATE_REPOSITORY,
    meta: {
      promiseScope: [namespace, repo],
      promise: Repos.updateRepository({namespace, repo}, data)
        .then(
          (resp) => resp.data,
          (error) => {
            throw {
              _error: error.data.errors
            };
          }
        ),
      notifications: {
        pending: 'Updating repository',
        success: 'Repository saved',
        error: 'There was an error saving your repository'
      }
    }
  };
}

export function deleteRepository(params) {
  return {
    type: consts.repositories.DELETE_REPOSITORY,
    params,
    meta: {
      promiseScope: [params.namespace, params.repo],
      promise: Repos.deleteRepository(params).then((resp) => resp.data),
      notifications: {
        pending: 'Deleting repository',
        success: 'Repository deleted',
        error: 'There was an error deleting your repository'
      }
    }
  };
}

export function deleteRepoManifest({ namespace, repo, reference }) {
  return {
    type: consts.repositories.DELETE_REPO_MANIFEST,
    meta: {
      promiseScope: [namespace, repo, reference],
      promise: Repos.deleteRepoManifest({namespace, repo, reference})
        .then(() => {
          return {namespace, repo, tag: reference};
        }),
      notifications: {
        pending: `Deleting tag ${reference}`,
        success: `Deleted tag ${reference}`,
        error: `There was an error deleting tag ${reference}`
      }
    }
  };
}

/**
 * Deletes multiple manifests
 *
 */
export function deleteRepoManifests({ namespace, repo, references }) {
  return (dispatch) => {

    let ref = references.shift().toJS();

    let deleteManifest = deleteRepoManifest({ namespace, repo, reference: ref.name });

    if (references.length > 0) {
      // recursively call this action with what's left of the references array
      deleteManifest.meta.promise.then(() => {
        dispatch(deleteRepoManifests({ namespace, repo, references }));
      });
    }

    dispatch(deleteManifest);
  };
}

export function getRepositoryTags({ namespace, repo }) {
  return {
    type: consts.repositories.GET_REPOSITORY_TAGS,
    meta: {
      promiseScope: [namespace, repo],
      promise: Repos.getRepoTags({namespace, repo})
        .then((response) => {
          return {
            namespace,
            name: repo,
            tags: response.data
          };
        })
      // API response:
      // {
      //   name: 'namespace/repo-name',
      //   tags: [
      //     {
      //       'name': string,
      //       "inRegistry": bool   // "true if the tag exists in Registry"`
      //       "hashMismatch": bool // "true if the hashes from notary and registry don't match"`
      //       "inNotary": bool     // "true if the tax exists in Notary"`
      //     }
      //     ...
      //   ]
      // }
    }
  };
}

/** Team namespace permissions **/
export function getTeamAccessToRepoNamespace({ orgName, teamName }) {
  return {
    type: repositories.GET_TEAM_ACCESS_TO_REPO_NAMESPACE,
    meta: {
      promiseScope: [orgName, teamName],
      promise: Repos.getRepoNamespaceTeamAccess({namespace: orgName, team: teamName}).then(
        response => {
          return response.data;
        },
        error => {
          throw {_error: error.data.errors};
        }
      )
    }
  };
}

export function grantTeamAccessToRepoNamespace({ orgName, teamName, accessLevel}) {
  return {
    type: repositories.GRANT_TEAM_ACCESS_TO_REPO_NAMESPACE,
    meta: {
      promiseScope: [orgName, teamName],
      promise: Repos.grantRepoNamespaceTeamAccess({namespace: orgName, team: teamName}, {accessLevel}).then(
        response => {
          return response.data;
        },
        error => {
          throw {_error: error.data.errors};
        }
      ),
      notifications: {
        pending: `Allowing ${accessLevel} for all repositories of ${orgName} for the team ${teamName}`,
        success: `Allowed ${accessLevel} for all repositories of ${orgName} for the team ${teamName}`,
        error: 'There was an error adding the repositories to the team'
      }
    }
  };
}

export function revokeTeamAccessToRepoNamespace({ orgName, teamName }) {
  return {
    type: repositories.REVOKE_TEAM_ACCESS_TO_REPO_NAMESPACE,
    meta: {
      promiseScope: [orgName, teamName],
      promise: Repos.revokeRepoNamespaceTeamAccess({namespace: orgName, team: teamName}).then(
        () => {
          return {
            namespace: {
              name: orgName
            },
            team: {
              name: teamName
            },
            accessLevel: ''
          };
        },
        error => {
          throw {_error: error.data.errors};
        }
      ),
      notifications: {
        pending: `Revoking access from all repositories of ${orgName} for the team ${teamName}`,
        success: `Revoked access from all repositories of ${orgName} for the team ${teamName}`,
        error: 'There was an error revoking permissions'
      }
    }
  };
}

/** Team repo permissions **/
export function changeTeamAccessToRepo({ orgName, teamName, repo, accessLevel}) {
  return {
    type: repositories.CHANGE_TEAM_ACCESS_TO_REPO,
    meta: {
      promiseScope: [orgName, teamName, repo],
      promise: Repos
        .grantRepoTeamAccess({namespace: orgName, repo, team: teamName}, {accessLevel})
        .then((response) => {
          return {
            orgName,
            teamName,
            teamId: response.data.team.id,
            repoName: repo,
            accessLevel
          };
        })
    }
  };
}

export function revokeTeamAccessToRepo({ orgName, teamName, repo }, teamId = '') {
  return {
    type: repositories.REVOKE_TEAM_ACCESS_TO_REPO,
    meta: {
      promiseScope: [orgName, teamName, repo],
      promise: Repos
        .revokeRepoTeamAccess({namespace: orgName, repo, team: teamName})
        .then(() => {
          return {
            orgName,
            teamName,
            teamId,
            repoName: repo
          };
        })
    }
  };
}
