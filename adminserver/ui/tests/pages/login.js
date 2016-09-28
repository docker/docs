'use strict';

const commands = {
  login: function (username, password) {
    return this
      .waitForElementVisible('body', 3000)
      .waitForElementVisible('@username', 1000)
      .setValue('@username', username)
      .setValue('@password', password)
      .click('@submitButton');
  },
  logout: function () {
    return this
      .waitForElementVisible('@userMenu', 1000)
      .moveToElement('@userMenu', 0, 0)
      .waitForElementVisible('@logoutLink', 1000)
      .click('@logoutLink')
      .waitForElementVisible('body', 3000)
      .assert.elementPresent('@logoutForm')
      .click('@submitButton');
  }
};

module.exports = {
  commands: [commands],
  url: function () {
    return this.api.globals.BASE_URL;
  },
  elements: {
    username: {
      selector: '#username'
    },
    password: {
      selector: '#password'
    },
    userMenu: '#user-menu',
    submitButton: '#submit_button',
    logoutLink: '#user-menu ul li:nth-child(4) a',
    logoutForm: '#logout_form',
    loginForm: '#login_form'
  }
};
