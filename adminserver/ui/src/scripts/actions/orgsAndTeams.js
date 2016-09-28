'use strict';

import {
  deleteOrganizationMember,
  listOrganizationMembers
} from './organizations.js';
import {
  deleteTeamMember,
  listTeamMembers
} from './teams.js';

/**
 * If given a team name this calls listTeamMembers; otherwise this calls
 * `listOrganizationMembers` from actions/organizations
 */
export const listOrgOrTeamMembers = ({ orgName, teamName, limit = 50 }) => {
  if (!teamName) {
    return listOrganizationMembers({ orgName, limit });
  } else {
    return listTeamMembers({ orgName, teamName, limit });
  }
};


/**
 * This action combines deleteOrganizationMember and deleteTeamMember; if
 * a team name exists we delete the member from the team by returning the
 * deleteTeamMember call. Otherwise, return the deleteOrganizationMember call.
 */
export const deleteOrgOrTeamMember = ({ orgName, teamName, memberName }) => {
  if (teamName === undefined) {
    return deleteOrganizationMember({
      name: orgName,
      member: memberName
    });
  }
  return deleteTeamMember({
    teamName: teamName,
    orgName,
    memberName
  });
};
