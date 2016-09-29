'use strict';
// Blob is a polyfill
require('vendor/Blob');
import { saveAs } from 'vendor/FileSaver';

const debug = require('debug')('hub:actions:downloadInvoice');

module.exports = function(actionContext, { JWT, username, invoiceId }) {
  const xhr = new XMLHttpRequest();
  xhr.open('GET', process.env.REGISTRY_API_BASE_URL + '/api/billing/v3/account/' + username + '/invoices/' + invoiceId + '/');
  xhr.setRequestHeader('Authorization', 'JWT ' + JWT);
  xhr.responseType = 'blob';
  xhr.onload = function() {
    if (xhr.status === 200) {
      const blob = new Blob([xhr.response], {
        type: 'application/pdf'
      });
      saveAs(blob, 'docker_invoice_' + invoiceId + '.pdf');
    } else {
      debug('error');
    }
  };
  xhr.send();
};
