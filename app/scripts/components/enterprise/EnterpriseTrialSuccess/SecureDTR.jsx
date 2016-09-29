'use strict';

import React, { Component } from 'react';
import styles from '../EnterpriseTrialSuccess.css';

export default class SecureDTR extends Component {
  render() {
    return (
      <div>
        <div className={styles.section}>
          <h4>Configure UCP</h4>
          <div>Once you’ve configured, you’re ready to start deploying and managing your container environment.</div>
          <div>
            <ul>
              <li>Configure UCP using the docs <a href='https://docs.docker.com/ucp/'>here</a>.</li>
              <li>Get started with UCP with this <a href='https://github.com/docker/ucp_lab'>hands-on lab</a>.</li>
              <li><a href='https://forums.docker.com/c/commercial-products/ucp'>Visit the forums</a> for more help.</li>
            </ul>
          </div>
        </div>
        <div className={styles.section}>
          <h4>Secure DTR</h4>
          <div>
            { 'You’re almost ready to push and pull images! You need to ' }
            <span className={styles.emphasized}>secure your Trusted Registry </span>
            first. Navigate to <i> Settings > Security</i>{', and enter your ' +
            'data in the required fields. At this time, you may also want to ' +
            'configure additional settings such as ports, storage, authentication, and so forth.' }
          </div>
          <div>
            <ul>
              <li>Configure DTR using the docs <a href='https://docs.docker.com/docker-trusted-registry/configure/configuration/'>here</a>.</li>
              <li>Get started with DTR with this <a href='https://docs.docker.com/docker-trusted-registry/quick-start/'>Quickstart Guide</a>.</li>
              <li><a href='https://forums.docker.com/c/commercial-products/dtr'>Visit the forums</a> for more help.</li>
            </ul>
          </div>
        </div>
      </div>
    );
  }
}
