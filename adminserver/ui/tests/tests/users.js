'use strict';

module.exports = {
  'Can a new user be created?': (client) => {

    const username = client.globals.DTR_ADMIN;
    const password = client.globals.DTR_PASSWORD;

    const dtr = client.page.login();
    const users = client.page.users();

    const newUsername = users.getUsername();

    dtr.navigate()
      .login(username, password)
      .waitForElementVisible('@userMenu', 3000);

    users
      .navigate()
      .waitForElementVisible('@newUserButton', 3000)
      .click('@newUserButton')
      .createNewUser();

    client.url(`${client.globals.BASE_URL}/users/${newUsername}`);

    users
      .waitForElementVisible('@usernameHeader', 3000);

    // yuk
    // render delay so we wait for half a second
    client.pause(500);

    users
      .assert.containsText('@usernameHeader', newUsername);
  },
  'Can a user be deleted?': (client) => {
    const users = client.page.users();
    users
      .click('@userSettingsTab')
      .waitForElementVisible('@deleteUserButton', 3000)
      .click('@deleteUserButton')
      .waitForElementVisible('@deleteModalInput', 3000)
      .setValue('@deleteModalInput', 'DELETE')
      .click('@deleteModalDeleteButton')
      .waitForElementVisible('@newUserButton', 3000);

    client.assert.urlEquals(`${client.globals.BASE_URL}/users`);
    client.end();
  }
};
