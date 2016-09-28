'use strict';
import _ from 'lodash';

module.exports = {
  'Create X repositories, where X is env var REPOS': (client) => {

    const username = client.globals.DTR_ADMIN;
    const password = client.globals.DTR_PASSWORD;

    const dtr = client.page.login();
    const repos = client.page.repositories();

    dtr.navigate()
      .login(username, password)
      .waitForElementVisible('@userMenu', 3000);

    repos.navigate();

    _.range(parseInt(process.env.REPOS)).map(() => {
      repos.clearValue('@repoNameInput');
      repos.clearValue('@shortDescriptionInput');
      repos.createRepo(`testrepo-${Date.now().toString(16)}`);
      client.pause(1000);
    });
  }
}
