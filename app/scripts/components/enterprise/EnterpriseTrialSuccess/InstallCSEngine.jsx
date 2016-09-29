'use strict';

import React, { Component } from 'react';
import CopyCodeBox from 'common/CopyCodeBox';
import styles from '../EnterpriseTrialSuccess.css';
import {DEB, RPM} from 'common/data/csEngineInstructions';
const debug = require('debug')('InstallCSEngine');

export default class InstallCSEngine extends Component {
 render() {
    return (
      <div>
        <div className={styles.section}>
          To install Docker Datacenter, you must first install
          the commercially supported Docker engine.
        </div>
        <div className={styles.smallSection}>
          <span className={styles.emphasized}>
            Copy and run either the RPM or the DEB specific instructions in your terminal.
          </span>
        </div>
        <div className={styles.copyLabel}>RPM</div>
        <div className={styles.copyBox}>
          <CopyCodeBox content={RPM}
                       dollar={true}
                       lines={4} />
        </div>
        <div className={styles.copyLabel}>DEB</div>
        <div className={styles.copyBox}>
          <CopyCodeBox content={DEB}
                       dollar={true}
                       lines={5} />
        </div>

        <div className={styles.copyLabel}>
          For more details about this installation,&nbsp;
          <a href='http://docs.docker.com/docker-trusted-registry/install/install-csengine/' target='_blank'>
            view our documentation
          </a>
        </div>
      </div>
    );
  }
}
