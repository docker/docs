'use strict';

module.exports = {
  'Can the default user log in?': (client) => {

    const username = client.globals.DTR_ADMIN;
    const password = client.globals.DTR_PASSWORD;

    const dtr = client.page.login();

    dtr.navigate()
      .login(username, password)
      .waitForElementVisible('@userMenu', 3000)
      .assert.containsText('@username', username);
  },
  'Can the logged in user log out?': (client) => {

    const dtr = client.page.login();

    dtr
      .logout()
      .waitForElementVisible('body', 3000)
      .waitForElementVisible('@loginForm', 3000)
      .assert.elementPresent('@username');

    client.end();
  }
}
