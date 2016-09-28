'use strict';

module.exports = {
  'Is DTR Running?': (browser) => {
    browser
      .url(browser.globals.BASE_URL)
      .waitForElementVisible('body', 3000)
      .waitForElementVisible('#header', 3000)
      .assert.visible('#header')
      .end();
  }
};
