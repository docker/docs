'use strict';

const licensePath = require('path').resolve(__dirname, '../util/license.lic');

const commands = {
  applyLicense: function () {

    console.log('trying: ' + licensePath);

    this
      .setValue('@licenseInput', licensePath)
      .click('@applyLicenseButton');
    this.api.pause(1000);
  }
};

module.exports = {
  commands: [commands],
  url: function () {
    return `${this.api.globals.BASE_URL}/admin/settings/general`;
  },
  elements: {
    licenseInput: 'input#newLicense',
    applyLicenseButton: '#apply-license-button'
  }
};
