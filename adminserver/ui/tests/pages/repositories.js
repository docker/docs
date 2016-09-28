'use strict';

const testRepoName = `testrepo-${Date.now().toString(16)}`;
const testRepoDescription = 'Test repository.';

const commands = {
  createRepo: function (repoName = testRepoName) {
    return this
      .waitForElementVisible('@newRepoButton', 3000)
      .click('@newRepoButton')
      .waitForElementVisible('@repoNameInput', 3000)
      .setValue('@repoNameInput', repoName)
      .setValue('@shortDescriptionInput', testRepoDescription)
      .click('@createRepoSubmit')
      .waitForElementVisible('@newRepoButton', 3000);
  },
  getTestRepoName: function () {
    return testRepoName;
  }
};

module.exports = {
  commands: [commands],
  url: function () {
    return `${this.api.globals.BASE_URL}/repositories`;
  },
  elements: {
    newRepoButton: '#new-repo-button button',
    repoNameInput: 'input[name=name]',
    shortDescriptionInput: 'input[name=shortDescription]',
    createRepoSubmit: '#create-repo-form button[type=submit]',
    repoSettingsTab: '#repo-settings-tab a',
    deleteRepoButton: '#delete-repo-button button',
    deleteModalInput: '#delete-modal input',
    deleteModalDeleteButton: '#delete-modal button[type=submit]'
  }
};
