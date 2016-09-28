'use strict';

module.exports = {
  'Open DTR and log in as the default user': (client) => {

    const username = client.globals.DTR_ADMIN;
    const password = client.globals.DTR_PASSWORD;

    const dtr = client.page.login();

    dtr.navigate()
      .login(username, password);
  }
}
