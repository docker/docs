'use strict';

import _ from 'lodash';
import { Accounts } from 'dtr-js-sdk';
import consts from 'consts';
import { TeamRecord, UserRecord, OrgMemberRecord } from 'records';
import { TeamSchema, UserSchema } from 'normalizers';
import { arrayOf, normalize } from 'normalizr';
import { normalizedToRecords } from 'utils/records';

const AccountsAPI = Accounts.default;

// The fucking enzi api split getting team information into two endpoints:
// one for the team ID, name and description and the other for sync info
// (ldap/managed etc)
//
// here we're just going to combine both APIs and ignore it actionside.
//
// we return objects and make a new team record based off of both responses
export function getTeam(orgName, teamName) {
  const combinator = Promise.all([
    AccountsAPI.getTeam({ name: orgName, team: teamName }).then((resp) => {
      return {
        orgName,
        team: resp.data
      };
    }),
    AccountsAPI.getTeamSyncConfig({ name: orgName, team: teamName }).then((resp) => {
      return {
        orgName,
        team: { name: teamName, ...resp.data }
      };
    })
  ]);

  return {
    type: consts.teams.GET_TEAM,
    meta: {
      promiseScope: [orgName, teamName],
      promise: combinator
    }
  };
}

export function listTeams({ orgName, limit = 100 }) {
  return {
    type: consts.teams.LIST_TEAMS,
    meta: {
      promiseScope: [orgName],
      promise: AccountsAPI.listTeams({ name: orgName }, { limit }).then(
        (resp) => {
          const norm = normalize(resp.data.teams, arrayOf(TeamSchema));
          return {
            orgName,
            data: normalizedToRecords(norm, { 'team': TeamRecord })
          };
        }
      )
    }
  };
}

export function listTeamMembers({ orgName, teamName, limit = 50 }) {
  return {
    type: consts.teams.LIST_MEMBERS,
    meta: {
      promiseScope: [orgName, teamName],
      promise: AccountsAPI
        .listTeamMembers({ name: orgName, team: teamName }, { limit })
        .then( response => {
          const data = normalize(response.data.members, arrayOf(UserSchema));

          return {
            orgName,
            teamName,
            data: normalizedToRecords(data, { 'user': OrgMemberRecord })
          };
        })
    }
  };
}

export function createTeam({ orgName, type, name, ldapDN, ldapGroupMemberAttribute }) {
  return {
    type: consts.teams.CREATE_TEAM,
    meta: {
      promiseScope: [orgName, name],
      promise: AccountsAPI.createTeam({ name: orgName }, {
        type,
        name,
        ldapDN,
        ldapGroupMemberAttribute
      }).then(resp => {
        // after creating a team we must save ldap details, if we need to.
        if (type === 'managed') {
          return {
            orgName,
            team: new TeamRecord(resp.data)
          };
        }

        // format the ldap data into what enzi expects...
        const data = {
          enableSync: true,
          selectGroupMembers: true,
          groupDN: ldapDN,
          groupMemberAttr: ldapGroupMemberAttribute
        };

        return AccountsAPI
          .updateTeamSyncConfig({ name: orgName, team: name }, data)
          .then(syncResp => {
            return {
              orgName,
              team: new TeamRecord({
                ...resp.data,
                ...syncResp.data
              })
            };
          });
      })
    }
  };
}

/**
 * Adds a team within an org
 *
 * TODO: Is this needed?
 */
export function createTeamWithUsers({ orgName, team }, users = []) {
  // TODO:  Improve error handling with N+1 API calls.  If the actual team
  // creation API call fails the team resource has not yet been created.  If
  // adding users to the team fails, the team exists and certain members (who
  // knows which ones...) may not be in the team.
  //
  // The error message should reflect these states.
  //
  // See pending API update regarding adding users to the team at once
  return {
    type: consts.teams.CREATE_TEAM,
    meta: {
      promiseScope: [orgName, team.name],
      promise: AccountsAPI.createTeam({ name: orgName }, team)
        .then( () => {
          // If this is an LDAP team we won't have any users - only add users
          // if this is a managed team and users exist
          if (users.length > 0) {
            let promises = users.map( user => {
              return AccountsAPI.addTeamMember({
                name: orgName,
                team: team.name,
                member: user.name
              });
            });
            // After all members have been added return the new team and all
            // members for the reducer to add to the teams list
            return Promise.all(promises).then( () => {
              return {
                orgName,
                team: new TeamRecord(team),
                members: _.mapValues(_.indexBy(users, 'name'), (member) => new UserRecord(member))
              };
            });
          } else {
            return {
              orgName,
              team: new TeamRecord(team),
              members: {}
            };
          }
        }, error => {
          throw {
            _error: error.data.errors
          };
        }
      ),
      notifications: {
        pending: 'Creating team',
        success: 'Your team was created',
        error: 'There was an error creating your team'
      }
    }
  };
}

export function getTeamsForUser({ name }) {
  return {
    type: consts.teams.GET_TEAMS_FOR_USER,
    meta: {
      promise: AccountsAPI.listUserOrganizations({ name }, { limit: 100 })
        .then((response) => {
          const {
            data: {
              memberOrgs
            }
          } = response;

          // extract the org names from the response
          const orgNames = memberOrgs.map((org) => {
            return org.org.name;
          });

          // get the list of teams for each org
          const promises = orgNames.map((org) => {
            return AccountsAPI.listOrganizationMemberTeams({ name: org, member: name }).then((r) => {
              return r.data.memberTeams.map((team) => {
                team.team.orgName = org;
                return team;
              });
            });
          });

          return Promise.all(promises).then((r) => {
            return {
              user: name,
              teamsByOrg: r
            };
          });

        }, error => {
          throw {
            _error: error.data.errors
          };
        })
    }
  };
}

export function deleteTeam({ orgName, teamName }) {
  return {
    type: consts.teams.DELETE_TEAM,
    meta: {
      promiseScope: [orgName, teamName],
      promise: AccountsAPI.deleteTeam({ name: orgName, team: teamName })
        .then(() => {
          return {
            orgName,
            teamName
          };
        }).then(() => {
          // TODO make this more not suck
          window.location.href = `/orgs/${ orgName }/users`;
        }),
      notifications: {
        pending: 'Deleting team',
        success: `Team ${teamName} was deleted`,
        error: `There was an error deleting ${teamName}`
      }
    }
  };
}

export function updateTeam(orgName, teamName, teamData) {
  // When switching from an LDAP team to managed we need to remove any LDAP
  // fields such as the DN; the API fails if these are sent alongside a managed
  // team.
  // We handle this here so every form need not worry about deleting these
  // fields during submission
  if (teamData.type === 'managed') {
    delete teamData.ldapDN;
    delete teamData.ldapGroupMemberAttribute;
  }

  return {
    type: consts.teams.UPDATE_TEAM,
    meta: {
      promiseScope: [orgName, teamName],
      promise: AccountsAPI.updateTeam({ name: orgName, team: teamName }, teamData)
        .then((resp) => {
          // format the ldap data into what enzi expects...
          const data = {
            enableSync: teamData.type === 'ldap',
            ...teamData,
            selectGroupMembers: true,
            searchBaseDN: '',
            searchScopeSubtree: false,
            searchFilter: ''
          };

          return AccountsAPI
            .updateTeamSyncConfig({ name: orgName, team: resp.data.name }, data)
            .then(syncResp => {
              if (teamName !== resp.data.name) {
                window.location.href = `/orgs/${orgName}/teams/${resp.data.name}/settings`;
              }

              return {
                orgName,
                team: new TeamRecord({
                  ...resp.data,
                  ...syncResp.data
                })
              };
            });
        }),
      notifications: {
        pending: `Updating team ${teamName}`,
        success: `Successfully updated ${teamName}`,
        error: `There was an error updating ${teamName}`
      }
    }
  };
}

export function getTeamSyncOptions({ orgName, teamName }) {
  return {
    type: consts.teams.GET_TEAM_SYNC,
    meta: {
      promiseScope: [orgName, teamName],
      promise: AccountsAPI
        .getTeamSyncConfig({ name: orgName, team: teamName })
    }
  };
}

export function checkTeamMembership(orgName, teamName, memberName) {
  return {
    type: consts.teams.GET_TEAM_MEMBER,
    meta: {
      promiseScope: [orgName, teamName, memberName],
      promise: AccountsAPI
        .checkTeamMembership({ name: orgName, team: teamName, member: memberName })
        .then(
          () => ({
            // This endpoint returns a 204 No Content if the user is
            // within this group.  We only need to return an object for
            // the teams reducer.
            orgName,
            teamName,
            memberName,
            isMember: true
          }),
          (error) => {
            if (error.status === 404 || error.status === 403) {
              // Not found is not an error
              return {
                orgName,
                teamName,
                memberName,
                isMember: false
              };
            }
          }
        )
    }
  };
}

export function addTeamMember({ name, team, member }) {
  return {
    type: consts.teams.ADD_TEAM_MEMBER,
    meta: {
      promiseScope: [name, team],
      promise: AccountsAPI
        .addTeamMember({ name, team, member })
        .then((resp) => ({
          name,
          team,
          member: new UserRecord(resp.data.member)
        }))
    },
    notifications: {
      pending: `Adding ${member} to team ${team}`,
      success: `${member} added to ${team}`,
      error: `There was an error adding ${member} to ${team}`
    }
  };
}

// Add multiple users to a team
export function addTeamMembers(orgName, teamName, memberNames) {
  return {
    type: consts.teams.ADD_TEAM_MEMBERS,
    meta: {
      // TODO:how to scope this? If you call addTeamMembers with a different set of members in succession before the
      // old call can finish, the original call won't fire it's resolved action. This can make our local state's list
      // of members get out of sync.
      // Maybe make sure to block user from adding more members until this is done?
      promiseScope: [orgName, teamName],
      // TODO: If for example one of the add member fails, then Promise.all will fail but the ones which succeeded
      // should still be added to our local state
      promise: Promise.all(
                 memberNames.map((memberName) => {
                   return AccountsAPI.addTeamMember({ name: orgName, team: teamName, member: memberName })
                     .then((resp) => {
                       return resp.data;
                     });
                 })
               ).then((addedMembers) => {
                  const membersArray = addedMembers.map((m) => {
                    return m.member;
                  });
                   return {
                     orgName,
                     teamName,
                     members: _.mapValues(_.keyBy(membersArray, 'name'), (member) => new UserRecord(member))
                   };
               }),
      notifications: {
        pending: `Adding to team ${teamName}`,
        success: `${memberNames.length} new members were added to ${teamName}`,
        error: `There was an error adding members to ${teamName}`
      }
    }
  };
}

// Delete a single member from a team
export function deleteTeamMember({ orgName, memberName, teamName }) {
  return {
    type: consts.teams.DELETE_TEAM_MEMBER,
    meta: {
      promiseScope: [orgName, teamName, memberName],
      promise: AccountsAPI.deleteTeamMember({ name: orgName, team: teamName, member: memberName })
        .then(() => ({
          orgName,
          teamName,
          memberName
        }))
    }
  };
}
