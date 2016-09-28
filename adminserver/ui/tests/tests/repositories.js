'use strict';



module.exports = {
  'Can the user create a repo under their own namespace?': (client) => {

    const username = client.globals.DTR_ADMIN;
    const password = client.globals.DTR_PASSWORD;

    const loginPage = client.page.login();
    const repoPage = client.page.repositories();

    loginPage.navigate()
      .login(username, password)
      .waitForElementVisible('@userMenu', 3000);

    // now use the repository page object to see the list
    repoPage.navigate()
      .createRepo();

    // navigate to the newly created repo
    client.url(`${repoPage.url}/${client.globals.DTR_ADMIN}/${repoPage.getTestRepoName()}`);

  },

  'Can the user delete a repo under their own namespace?': (client) => {
    const repoPage = client.page.repositories();

    // a little time for the UI to catch up
    client.pause(500);

    repoPage
      .waitForElementVisible('@repoSettingsTab', 3000)
      .click('@repoSettingsTab')
      .waitForElementVisible('@deleteRepoButton', 3000)
      .click('@deleteRepoButton')
      .waitForElementVisible('@deleteModalInput', 3000)
      .setValue('@deleteModalInput', 'DELETE')
      .click('@deleteModalDeleteButton');

    client.end();
  }
};
