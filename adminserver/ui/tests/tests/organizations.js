'use strict';

import { names } from '../pages/organizations';

const patience = 3000;

module.exports = {
  'Can an org be created?': (client) => {

    const username = client.globals.DTR_ADMIN;
    const password = client.globals.DTR_PASSWORD;
    const orgList = client.page.organizations();

    client.page.login()
      .navigate()
      .login(username, password)
      .waitForElementPresent('@userMenu', patience);

    orgList
      .navigate()
      .waitForElementPresent('@newOrgButton', patience)
      .createOrg();

    client.pause(500).assert.urlEquals(`${client.globals.BASE_URL}/orgs/${names.org}/users`);

  },
  'Can a user be added to an org?': (client) => {

    const orgList = client.page.organizations();

    orgList
      .waitForElementPresent('@addUserButton', patience);

    client.pause(500);

    orgList
      .click('@addUserButton')
      .waitForElementPresent('@addUserUI', patience)
      .waitForElementPresent('@userCreateToggle', patience)
      .click('@userCreateToggle')
      .waitForElementPresent('@newUserForm', patience)
      .waitForElementPresent('@newUsernameInput', patience)
      .setValue('@newUsernameInput', names.user)
      .setValue('@newUserPasswordInput', 'password')
      .setValue('@newUserFullNameInput', 'Test User')
      .submitForm('@newUserForm');
  },

  'Can a team be created?': (client) => {
    client.page.organizations().createTeam();
    client.url(`${client.globals.BASE_URL}/orgs/${names.org}/teams/${names.team}/settings`);
    client.pause(1000);
  },
  'Can a team be deleted?': (client) => {
    // go to the team settings page
    client.url(`${client.globals.BASE_URL}/orgs/${names.org}/teams/${names.team}/settings`);
    // delete the team; the 'delete' button has ID deleteOrgButton, so the deleteOrg()
    // func works here.
    client.page.organizations().deleteOrg();
  },
  'Can an org be deleted?': (client) => {
    // go to the org settings page
    client.url(`${client.globals.BASE_URL}/orgs/${names.org}/settings`);

    // delete the org
    client.page.organizations().deleteOrg();

    client.end();
  }
};
