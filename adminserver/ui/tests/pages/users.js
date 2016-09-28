'use strict';

const testUserName = `user-${Date.now().toString(16)}`;
const password = 'password';
const fullName = 'Test User';

const commands = {
  createNewUser: function () {
    return this
      .waitForElementVisible('@usernameInput', 3000)
      .setValue('@usernameInput', testUserName)
      .setValue('@passwordInput', password)
      .setValue('@fullNameInput', fullName)
      .click('@newUserFormSubmit');
  },
  getUsername: () => {
    return testUserName;
  }
};

module.exports = {
  commands: [commands],
  url: function () {
    return `${this.api.globals.BASE_URL}/users`;
  },
  elements: {
    newUserButton: '#new-user-button',
    usernameInput: 'input[name=username]',
    passwordInput: 'input[name=password]',
    fullNameInput: 'input[name=fullName]',
    newUserFormSubmit: '#new-user-form button[type=submit]',
    usernameHeader: '#user-detail-username',
    userSettingsTab: '#user-settings-tab',
    deleteUserButton: '#delete-user button',
    deleteModalInput: '#delete-modal input[type=text]',
    deleteModalDeleteButton: '#delete-modal button[type=submit]'
  }
};
