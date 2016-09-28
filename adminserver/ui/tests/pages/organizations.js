'use strict';

const names = {
  org: `testorg-${Date.now().toString(16)}`,
  team: `testteam-${Date.now().toString(16)}`,
  user: `testuser-${Date.now().toString(16)}`
};

const patience = 2000;

const commands = {
  createOrg: function () {
    this
      .click('@newOrgButton')
      .waitForElementVisible('@orgNameInput', patience)
      .setValue('@orgNameInput', [names.org, this.api.Keys.ENTER]);

    this.api.pause(1000);

    this.api.url(`${this.api.globals.BASE_URL}/orgs/${names.org}/users`);
  },
  createTeam: function () {
    this
      .waitForElementPresent('@addTeamButton', patience)
      .click('@addTeamButton')
      .waitForElementPresent('@teamNameInput', patience)
      .setValue('@teamNameInput', names.team)
      .waitForElementPresent('@createTeamSubmit', patience)
      .click('@createTeamSubmit');
  },
  deleteOrg: function () {
    this
      .waitForElementPresent('@deleteOrgButton', patience)
      .click('@deleteOrgButton')
      .waitForElementPresent('@deleteOrgModalInput', patience)
      .setValue('@deleteOrgModalInput', 'DELETE')
      .click('@deleteOrgModalSubmit')
      .api.pause(1000);
  }
};

module.exports = {
  commands: [commands],
  url: function () {
    return `${this.api.globals.BASE_URL}/orgs`;
  },
  elements: {
    newOrgButton: '#new-organization-button',
    orgNameInput: '#new-org-form input[name=name]',
    orgNameForm: '#new-org-form',
    newOrgSubmit: '#new-org-form button[type=submit]',
    orgSettingsTab: '#org-settings-tab a',
    deleteOrgButton: '#delete-org-button',
    deleteOrgModalInput: '#delete-org-modal input[type=text]',
    deleteOrgModalSubmit: '#delete-org-modal button[type=submit]',
    addTeamButton: '#add-team-button',
    teamNameInput: '#create-team-form input[name=name]',
    createTeamForm: '#create-team-form',
    createTeamSubmit: '#create-team-form button[type=submit]',
    addUserButton: '#add-user-button',
    addUserUI: '#add-user-ui',
    userCreateToggle: '#user-create-toggle .choice:nth-child(2)',
    newUsernameInput: '#new-user-form input[name=username]',
    newUserPasswordInput: '#new-user-form input[name=password]',
    newUserFullNameInput: '#new-user-form input[name=fullName]',
    newUserForm: '#new-user-form'
  },
  names: names
};
