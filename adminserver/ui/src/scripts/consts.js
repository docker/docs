'use strict';

import { keyMirror as km } from 'utils';

let consts = km([
  'ADD_GROWL_NOTIFICATION',
  'REMOVE_GROWL_NOTIFICATION',
  'UPDATE_GROWL_NOTIFICATION',
  'ADD_BANNER_NOTIFICATION',
  'REMOVE_BANNER_NOTIFICATION',
  'UPDATE_BANNER_NOTIFICATION',
  'SEARCH',
  'SEARCH_NAMESPACES'
]);

consts.notification = {
    DURATION: 3000
};

consts.ui = km([
  'CHANGE_PAGE',
  'RESET_STATE',
  'RESET_UI',
  'UPDATE_UI'
]);

consts.loading = km([
  'PENDING',
  'SUCCESS',
  'FAILURE'
]);
consts.loading.NOT_IN_PROGRESS = undefined; // For convenience if status is not set it is not in progress

consts.status = km([
  'RESET_STATUS'
]);

consts.storage = km([
]);

consts.logs = km([
  'FETCH_LOGS',

  'OBSERVE_LOGS',
  'UNOBSERVE_LOGS',
  'OBSERVE_LOG_DATA'
]);

consts.auth = km([
  // Fired from login screen to inform us we're not signed in
  'NOT_SIGNED_IN',
  'AUTH',
  'LOG_OUT'
]);

consts.metrics = km([
  'FETCH_METRICS',

  'OBSERVE_METRICS',
  'UNOBSERVE_METRICS',
  'OBSERVE_METRICS_DATA',
  'SEND_ANALYTICS'
]);

consts.settings = km([
  'ALL_SETTINGS',
  'FETCH_SETTINGS',
  'GET_STORAGE_SETTINGS',
  'AUTH_SETTINGS',
  'SAVE_AUTH_SETTINGS',
  'LDAP_SETTINGS',

  // Save actions
  'YAML_STORAGE_SAVE',
  'FORM_STORAGE_SAVE',
  'GENERAL_SETTINGS_SAVE',
  'AUTH_SETTINGS_SAVE',
  'LDAP_SETTINGS_SAVE',
  'SSL_SETTINGS_SAVE',
  'LICENSE_SETTINGS_SAVE',
  'TOGGLE_LICENSE_AUTOREFRESH',
  'GET_LICENSE',

  // Action called when clearing error/success notificaitons
  'REGISTRY_CLEAR',

  // GC
  'GET_GC_SCHEDULE',
  'UPDATE_GC_SCHEDULE',
  'DELETE_GC_SCHEDULE',
  'GET_LAST_GC_SAVINGS',
  'RUN_GC',
  'GET_GC_STATUS',
  'STOP_GC',

  // GC schedule types
  'UNTIL_DONE',
  'TIMEOUT',
  'NEVER',

  'LDAP_CHECK',
  'LDAP_SYNC',

  'LDAP_SYNC_TAKING_PLACE',

  // Updates
  'UPDATES'
]);

consts.settings.GC_POLLING_INTERVAL = 500;

// Used when restarting
consts.health = km([
  'DOWN',
  'UP',
  'HEALTH_UPGRADING'
]);

consts.updates = km([
  'UPDATE_AVAILABLE',
  'UPDATE_PENDING',
  'UPDATE_FAILURE',
  'UPDATE_SUCCESS',

  'UPDATE_NOTIFICATION_ID',
  'DOCKER_VERSION_NOTIFICATION_ID'
]);

consts.repositories = km([
  /** Repository **/
  'LIST_REPOSITORIES',
  'LIST_SHARED_REPOSITORIES',
  'LIST_ALL_REPOSITORIES',
  'CREATE_REPOSITORY',
  'GET_REPOSITORY',
  'GET_REPOSITORY_WITH_USER_PERMISSIONS',
  'UPDATE_REPOSITORY',
  'DELETE_REPOSITORY',

  'GET_REPOSITORY_TAGS',
  'GET_REPOSITORY_TAG_TRUST',
  'DELETE_REPO_MANIFEST',

  /** Repository access (which teams or users have access to a repo) **/
  'LIST_REPO_USER_ACCESS',
  'LIST_REPO_TEAM_ACCESS',

  /** Team access to namespace **/
  'GET_TEAM_ACCESS_TO_REPO_NAMESPACE',
  'GRANT_TEAM_ACCESS_TO_REPO_NAMESPACE',
  'REVOKE_TEAM_ACCESS_TO_REPO_NAMESPACE',

  /** Team access to repository **/
  'LIST_TEAM_ACCESS_TO_REPO',
  'GRANT_TEAM_ACCESS_TO_REPO',
  'CHANGE_TEAM_ACCESS_TO_REPO',
  'REVOKE_TEAM_ACCESS_TO_REPO',
  'CREATE_REPO_AND_GRANT_TEAM_ACCESS',

  /** UI **/
  'UPDATE_SEARCH',
  // These are valid values for filtering your list of repos
  'SHOW_MINE',
  'SHOW_SHARED',
  'SHOW_ALL',

  'TOGGLE_UPDATE_FORM'
]);


consts.organizations = km([
  /** Org **/
  'LIST_ORGANIZATIONS',
  'LIST_USER_ORGANIZATIONS',
  'CREATE_ORGANIZATION',
  'GET_ORGANIZATION',
  'DELETE_ORGANIZATION',

  /** Org member **/
  'LIST_ORGANIZATION_MEMBERS',
  'CHECK_MEMBERSHIP',
  'DELETE_MEMBER',
  'ADD_MEMBER',
  'UPDATE_MEMBER',
  'CREATE_ADD_MEMBER'
]);

consts.teams = km([
  /** Team **/
  'LIST_TEAMS',
  'CREATE_TEAM',
  'GET_TEAM',
  'UPDATE_TEAM',
  'DELETE_TEAM',
  'GET_TEAM_SYNC',
  'UPDATE_TEAM_SYNC',

  /** Team members **/
  'LIST_MEMBERS',
  'ADD_TEAM_MEMBER',
  'ADD_TEAM_MEMBERS',
  'GET_TEAM_MEMBER',
  'DELETE_TEAM_MEMBER',

  'GET_TEAMS_FOR_USER',

  /** UI **/
  'SHOW_ADD_REPO_TO_TEAM_FORM',
  'CLOSE_ADD_REPO_TO_TEAM_FORM',
    // These toggle between adding a new and existing repo to a team
    'ADD_EXISTING_REPO',
    'ADD_NEW_REPO',
  'SHOW_ADD_TEAM_FORM'
]);

consts.users = km([
  'LIST_USERS',
  'SEARCH_USERS',
  'CREATE_USER',
  'GET_USER',
  'CHANGE_USER_PASSWORD',
  'UPDATE_ACCOUNT',
  'DELETE_USER'
]);

consts.URLS = {
  UPDATES: '/api/v0/admin/upgrade',
  VERSION: '/admin/version'
};

export default consts;
