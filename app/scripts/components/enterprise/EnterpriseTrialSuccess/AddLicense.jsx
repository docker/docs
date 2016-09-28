'use strict';

import React, { PropTypes, Component } from 'react';
import connectToStores from 'fluxible-addons-react/connectToStores';
import AddTrialLicenseStore from 'stores/AddTrialLicenseStore';
import downloadLicenseContent from 'actions/downloadLicenseContent';
import Button from '@dux/element-button';
import { ATTEMPTING_DOWNLOAD,
         BAD_REQUEST,
         DEFAULT,
         FACEPALM,
         SUCCESSFUL_DOWNLOAD } from 'stores/addtriallicensestore/Constants';
const { func, object, string } = PropTypes;
import moment from 'moment';
import FA from 'common/FontAwesome';
import styles from '../EnterpriseTrialSuccess.css';
const debug = require('debug')('AddLicense');

class AddLicense extends Component {
  static propTypes = {
    error: string,
    JWT: string.isRequired,
    license: object.isRequired,
    STATUS: string
  }

  static contextTypes = {
    executeAction: func.isRequired
  }

  onClick = (e) => {
    const { JWT: jwt } = this.props;
    const { keyId, namespace } = this.props.license;
    this.context.executeAction(downloadLicenseContent, {
      jwt,
      keyId,
      namespace
    });
  }

  render() {
    const { error, license, STATUS } = this.props;
    const { expiration } = license;
    debug(license);

    let buttonText = 'Download License';
    let icon = <FA icon='fa-cloud-download' />;
    let variant = 'primary';
    let maybeError;
    if (STATUS === ATTEMPTING_DOWNLOAD) {
      buttonText = 'Downloading License';
      icon = <FA icon='fa-spinner fa-spin' />;
    } else if (STATUS === SUCCESSFUL_DOWNLOAD) {
      buttonText = 'Successfully Downloaded License';
      variant = 'success';
      icon = <FA icon='fa-check' />;
    } else if (STATUS === BAD_REQUEST || STATUS === FACEPALM) {
      maybeError = <div className={styles.error}>{ error }</div>;
      variant = 'alert';
    }

    const button = (
      <div className={styles.fullWidth}>
        <Button variant={variant}
                disabled={STATUS === FACEPALM}
                onClick={this.onClick}
                ghost>
         { icon } { buttonText }
        </Button>
      </div>
    );
    return (
      <div>
        <div className={styles.section}>
          <h4>Licensing UCP</h4>
          <div>
            You may <span className={styles.emphasized}>{ 'upload the license file, ' +
            'provided below, in '} <i> Settings > Licenses</i>.</span>
            <i> Save and Restart</i> UCP.
          </div>
        </div>
        <div className={styles.section}>
          <h4>Licensing DTR</h4>
          <div>
            {'Before proceeding, you must set the domain name to the full ' +
            'host-name of your Trusted Registry server (this is under '}<i>Settings </i>
            {' in Trusted Registry). Once you’ve saved and restarted the Trusted Registry, you may now '}
            <span className={styles.emphasized}>{ 'upload the license file, ' +
            'provided below, in '} <i> Settings > Licenses</i>.</span>
            <i> Save and Restart</i> once more.
          </div>
        </div>
        <div>
          {maybeError}
          <div className={styles.downloadBlockWrapper}>
            <div className={styles.downloadBlock}>
              <FA icon='fa-file-o' size='3x' />
            </div>
            <div className={styles.downloadBlock2}>
              <div>docker_subscription.lic</div>
              <div>Expires on {moment(expiration).format('l')}</div>
            </div>
            <br />
            { button }
          </div>
        </div>
      </div>
    );
  }
}

export default connectToStores(AddLicense,
	[ AddTrialLicenseStore ],
  function({ getStore }, props) {
    return getStore(AddTrialLicenseStore).getState();
  });
