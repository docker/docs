'use strict';

import Immutable from 'immutable';

export const OrganizationRecord = Immutable.Record({
  id: undefined,
  name: undefined,
  type: 'organization'
});

export const TeamRecord = Immutable.Record({
  id: 0,
  orgID: 0,

  name: '',
  description: '',

  // Is the current user that's making the API request a member of this org?
  clientUserIsMember: false,

  type: '', // 'managed' or 'ldap'

  // When listing teams in a repository this has an access level
  accessLevel: undefined,
  // ldap sync member options
  enableSync: false,
  selectGroupMembers: false,
  groupDN: '',
  groupMemberAttr: '',
  searchBaseDN: '',
  searchScopeSubtree: false,
  searchFilter: ''
});

export const UserRecord = Immutable.Record({
  id: undefined,
  fullName: '',
  type: 'user',
  name: '',
  ldapLogin: undefined, // string containing ldap search string
  isActive: false
});

export const OrgMemberRecord = Immutable.Record({
  isAdmin: false,
  member: {
    id: undefined,
    fullName: '',
    type: 'user',
    name: '',
    ldapLogin: undefined, // string containing ldap search string
    isActive: false
  }
});

export const RepositoryRecord = Immutable.Record({
  id: undefined,
  name: '',
  namespace: '',
  namespaceType: '', // 'user' or 'organization'
  shortDescription: '',
  longDescription: '',
  visibility: 'public', // 'private' or 'public'
  accessLevel: '' // ['read-write', 'read-only', 'admin']
});

// NamespaceRecord is used when searching namesapces via the index API.
// It represents either an org or a user
export const NamespaceRecord = Immutable.Record({
  id: undefined,
  type: '', // enum - 'user' or 'organization'
  name: '',
  // Below fields are applicable for a user only.
  isActive: true,
  ldapLogin: ''
});

