'use strict';

const debug = require('debug')('hub:actions:downloadLicenseContent');
import { Billing } from 'hub-js-sdk';
// Blob is a polyfill
require('vendor/Blob');
import { saveAs } from 'vendor/FileSaver';

export default function downloadLicenseContent(actionContext, { jwt, namespace, keyId }) {
  actionContext.dispatch('ATTEMPTING_LICENSE_DOWNLOAD_START');
  Billing.getLicenseDownloadContent(jwt, { keyId, namespace }, (err, res) => {
    if (err) {
      debug('error', err);
      if(err.response.badRequest) {
        const { detail } = err.response.body;
        if(detail) {
          actionContext.dispatch('DOWNLOAD_LICENSE_CONTENT_BAD_REQUEST', detail);
        }
      } else {
        actionContext.dispatch('DOWNLOAD_LICENSE_CONTENT_FACEPALM');
      }
    } else {
      actionContext.dispatch('RECEIVE_LICENSE_DOWNLOAD_CONTENT');
      //perform download as result of clicking download button
      const blob = new Blob([res.text], {
        type: 'text/plain;charset=utf-8'
      });
      saveAs(blob, `docker_subscription.lic`);
    }
  });
}
