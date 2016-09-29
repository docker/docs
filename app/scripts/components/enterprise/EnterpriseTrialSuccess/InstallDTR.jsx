'use strict';

import React, { Component } from 'react';
import styles from '../EnterpriseTrialSuccess.css';
import Button from '@dux/element-button';
import classnames from 'classnames';
import CopyCodeBox from 'common/CopyCodeBox';
const debug = require('debug')('InstallDTR');

export default class InstallDTR extends Component {
  render() {
    const copyBoxClasses = classnames({
      [styles.section]: true,
      [styles.noBottomMargin]: true
    });
    return (
      <div>
        <div className={styles.section}>
          <h4>To Install Universal Control Plane</h4>
          <div>To install UCP, please see the instructions <a href='http://docs.docker.com/ucp/'>here</a>.</div>
        </div>
        <div className={styles.section}>
          <h4>To Install Docker Trusted Registry</h4>
          <div>
            {'Once the engine is installed, '}
            <span className={styles.emphasized}> install Docker Trusted Registry
            by running the <span className={styles.code}>docker/dtr
            </span> container below.</span>
            {' This command pulls and runs Docker Trusted Registry on a container.'}
          </div>
        </div>
        <div className={copyBoxClasses}>
          <CopyCodeBox content={`docker run -it --rm docker/dtr install --ucp-insecure-tls`}
                       dollar={true} />
        </div>
        <div className={styles.section}>
          Then, <span className={styles.emphasized}>point your browser to
          <a className={styles.code}>{` https://<host-ip>/`}</a>.</span>
        </div>
        <div className={styles.section}>
          <i>
            {`NOTE: Your browser will warn you that this is an unsafe site, ` +
            `with a self-signed, untrusted certificate. This is normal and ` +
            `expected; please allow this connection temporarily.`}
          </i>
        </div>
      </div>
    );
  }
}
