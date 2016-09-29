'use strict';

import createStore from 'fluxible/addons/createStore';
import md5 from 'md5';
var debug = require('debug')('stores: UserStore');

var UserStore = createStore({
  storeName: 'UserStore',
  handlers: {
    RECEIVE_USER: '_receiveUserFromHub',
    RECEIVE_NAMESPACES: '_receiveNamespaces',
    LOGOUT: '_logout',
    EXPIRED_SIGNATURE: '_logout'
  },
  initialize: function() {
    this.dateJoined = '';
    this.fullName = '';
    this.gravatarEmail = '';
    this.gravatarUrl = '';
    this.isActive = false;
    this.isAdmin = false;
    this.isStaff = false;
    this.profileUrl = '';
    this.company = '';
    this.id = '';
    this.location = '';
    this.userType = 'User';
    this.username = '';
    this.namespaces = [];
  },
  _getGravatarUrl: function(email) {
    return 'https://secure.gravatar.com/avatar/' + md5( email.trim().toLowerCase() );
  },
  _receiveUserFromHub: function(user) {

    // jscs:disable requireCamelCaseOrUpperCaseIdentifiers
    this.dateJoined = user.date_joined;
    this.fullName = user.full_name;
    this.gravatarEmail = user.gravatar_email;
    //TODO: the url has to be handed off from the backend
    //This fix should be temporary
    this.gravatarUrl = (user.gravatar_url === user.gravatar_email) ?
    this._getGravatarUrl(user.gravatar_email) : user.gravatar_url;
    this.isActive = user.is_active;
    this.isAdmin = user.is_admin;
    this.isStaff = user.is_staff;
    this.profileUrl = user.profile_url;
    // jscs:enable

    this.company = user.company;
    this.id = user.id;
    this.location = user.location;
    this.userType = user.type;
    this.username = user.username;

    this.emitChange();
  },
  _receiveNamespaces: function(receivedNamespaces) {
    //This is required for creating a repository
    //Namespaces are attached to a user, due to permissions/access restrictions
    //Eg. {
    //  "namespaces": [
    //    "user",
    //    "org1",
    //    "org2"
    //  ]
    //}
    this.namespaces = receivedNamespaces.namespaces;

    this.emitChange();
  },
  _logout: function() {

    this.company = '';
    this.dateJoined = '';
    this.fullName = '';
    this.gravatarEmail = '';
    this.gravatarUrl = '';
    this.id = '';
    this.isActive = false;
    this.isAdmin = false;
    this.isStaff = false;
    this.location = '';
    this.profileUrl = '';
    this.userType = 'User';
    this.username = '';
    this.namespaces = [];

    this.emitChange();
  },
  getState: function() {
    return {
      company: this.company,
      dateJoined: this.dateJoined,
      fullName: this.fullName,
      gravatarEmail: this.gravatarEmail,
      gravatarUrl: this.gravatarUrl,
      id: this.id,
      isActive: this.isActive,
      isAdmin: this.isAdmin,
      isStaff: this.isStaff,
      location: this.location,
      profileUrl: this.profileUrl,
      userType: this.userType,
      username: this.username,
      namespaces: this.namespaces
    };
  },
  getUsername: function() {
    return this.username;
  },
  getNamespaces: function() {
    return this.namespaces;
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    debug('rehydrate', state);
    this.dateJoined = state.dateJoined;
    this.fullName = state.fullName;
    this.gravatarEmail = state.gravatarEmail;
    this.gravatarUrl = state.gravatarUrl;
    this.isActive = state.isActive;
    this.isAdmin = state.isAdmin;
    this.isStaff = state.isStaff;
    this.profileUrl = state.profileUrl;
    this.company = state.company;
    this.id = state.id;
    this.location = state.location;
    this.userType = state.userType;
    this.username = state.username;
    this.namespaces = state.namespaces;
  }
});

module.exports = UserStore;
