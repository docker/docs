'use strict';

const patience = 3000;

module.exports = {
  'Apply license': (client) => {

    const username = client.globals.DTR_ADMIN;
    const password = client.globals.DTR_PASSWORD;
    const settingsPage = client.page.settings();

    client.page.login()
      .navigate()
      .login(username, password)
      .waitForElementPresent('@userMenu', patience);

    settingsPage.navigate();

    client.pause(500);

    settingsPage.applyLicense();
  }
};
