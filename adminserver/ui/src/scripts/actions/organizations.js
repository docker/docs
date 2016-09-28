'use strict';

import { Accounts } from 'dtr-js-sdk';

import consts from 'consts';
import { OrganizationRecord, OrgMemberRecord } from 'records';
import { normalize, arrayOf } from 'normalizr';
import { OrgSchema, UserSchema } from 'normalizers';
import { normalizedToRecords } from 'utils/records';
import { addTeamMember } from 'actions/teams';
import { Map } from 'immutable';

const AccountsAPI = Accounts.default;

export function listOrganizations() {
  return {
    type: consts.organizations.LIST_ORGANIZATIONS,
    meta: {
      promise: AccountsAPI.listAccounts({}, { filter: 'orgs', limit: 100 })
        .then(response => {
          const data = normalize(response.data.accounts, arrayOf(OrgSchema));
          return normalizedToRecords(data, { 'org': OrganizationRecord });
        })
    }
  };
}

export function addOrganizationMember({ name, member }, isAdmin = false) {
  return {
    type: consts.organizations.ADD_MEMBER,
    meta: {
      promise: AccountsAPI.addOrgMember({
        name,
        member
      },
      {
        isAdmin: isAdmin,
        isPublic: true
      }).then((response) => {
        return {
          orgName: name,
          member: new Map(response.data.member)
        };
      })
    }
  };
}

export function updateOrganizationMember({ name, member, isAdmin, isPublic }) {
  return {
    type: consts.organizations.UPDATE_MEMBER,
    meta: {
      promise: AccountsAPI.updateOrgMember({
          name,
          member
        },
        {
          isAdmin: isAdmin,
          isPublic: isPublic
        }).then((response) => {

        const data = normalize(response.data, UserSchema);

        return {
          orgName: name,
          member: normalizedToRecords(data, { 'user': OrgMemberRecord })
        };
      }),
      notifications: {
        pending: 'Updating member',
        success: 'The member was updated',
        error: 'There was an error updating the member'
      }
    }
  };
}

export function createUserAndAddToOrg(data, name, member, team) {
  return dispatch => {
    dispatch({
      type: consts.organizations.CREATE_ADD_MEMBER,
      meta: {
        promise: AccountsAPI
          .createAccount({}, data)
          .then((response) => {
            // now add the member to the org
            dispatch(addOrganizationMember({ name, member }));
            if (team) {
              dispatch(addTeamMember({ name, team, member }));
            }
            return {
              orgName: name,
              member: new Map(response.data)
            };
          })
      }
    });
  };
}

export function listUserOrganizations({ name }) {
  return {
    type: consts.organizations.LIST_USER_ORGANIZATIONS,
    meta: {
      // TODO: Fix pagination
      promise: AccountsAPI
        .listUserOrganizations({ name }, { limit: 20 })
        .then(response => {
          const orgs = response.data.memberOrgs.map((o) => {
            return o.org;
          });
          const data = normalize(orgs, arrayOf(OrgSchema));
          return {
            username: name,
            orgs: normalizedToRecords(data, { 'org': OrganizationRecord })
          };
        })
    }
  };
}

/**
 * Lists all organizations if the user is an admin, or lists all orgs that
 * a user is a member of if they're a non-admin
 *
 */
export function listAdminOrUserOrganizations({ name, isAdmin = false }) {
  return (isAdmin) ? listOrganizations() : listUserOrganizations({ name });
}

/**
 * Fetch an organization's information by its name
 */
export function getOrganization(name) {
  return {
    type: consts.organizations.GET_ORGANIZATION,
    meta: {
      promiseScope: [name],
      promise: AccountsAPI.getAccount({ name })
        .then(response => {
          if (!response.data.isOrg) {
            throw new Error('Didn\'t fetch an organization');
          }
          return normalizedToRecords(
            normalize(response.data, OrgSchema),
            { 'org': OrganizationRecord }
          );
        })
    }
  };
}

export function createOrganization(params) {
  // Explicitly set the 'account' type as an organization
  params = {
    ...params,
    isOrg: true
  };

  return (dispatch) => {
    return dispatch({
      type: consts.organizations.CREATE_ORGANIZATION,
      meta: {
        promiseScope: [params.name],
        promise: AccountsAPI.createAccount({}, params)
          .then(
            (resp) => {
              dispatch(addOrganizationMember({
                name: params.name,
                member: params.user.name
              }, true));
              return {
                org: new OrganizationRecord(resp.data)
              };
            }
          )
        ,
        notifications: {
          pending: 'Creating organization',
          error: 'There was an error creating your organization'
        }
      }
    });
  };




}

export function deleteOrganization(params) {
  return {
    type: consts.organizations.DELETE_ORGANIZATION,
    meta: {
      promiseScope: [params.name],
      promise: AccountsAPI.deleteAccount(params).then(() => {
        window.location.href = '/orgs';
      }),
      notifications: {
        pending: 'Deleting organization',
        success: 'Your organization was deleted',
        error: 'There was an error deleting your organization'
      }
    }
  };
}

export function listOrganizationMembers({ orgName, limit = 50 }) {
  return {
    type: consts.organizations.LIST_ORGANIZATION_MEMBERS,
    meta: {
      promiseScope: [orgName],
      promise: AccountsAPI.listOrganizationMembers({
        name: orgName,
        limit
      }).then((response) => {
        const data = normalize(response.data.members, arrayOf(UserSchema));
        return {
          orgName,
          data: normalizedToRecords(data, { 'user': OrgMemberRecord })
        };
      })
    }
  };
}

/**
 * Check whether a user is a member of an organization. This endpoint returns
 * a 204 if the user is a member and a 404 if the user is not a member; we
 * capture the 404 error to handle this repsonse.
 *
 * @param {Object} param - an object of parameters
 * @param {string} param.name - organization name
 * @param {string} param.member - user name
 */
export function checkOrganizationMembership({ name, member }) {
  return {
    type: consts.organizations.CHECK_MEMBERSHIP,
    meta: {
      promiseScope: [name, member],
      promise: AccountsAPI.checkOrganizationMembership({ name, member })
        .then(
          () => { return { name, member, isMember: true }; },
          (resp) => {
            if (resp.status === 404 || resp.status === 403) {
              return { name, member, isMember: false };
            }
            throw resp;
          }
        )
    }
  };
}

/**
 * Delete a member of an organization from all teams
 *
 * @param {Object} param - an object of parameters
 * @param {string} param.name - organization name
 * @param {string} param.member - user name
 */
export function deleteOrganizationMember({ name, member }) {
  return {
    type: consts.organizations.DELETE_MEMBER,
    meta: {
      promiseScope: [name, member],
      promise: AccountsAPI.deleteOrganizationMember({ name, member })
        .then(() => {
          return {
            orgName: name,
            memberName: member
          };
        })
    }
  };
}
