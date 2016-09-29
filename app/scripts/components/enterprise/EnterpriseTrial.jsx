'use strict';

import React, { Component, PropTypes } from 'react';
const debug = require('debug')('EnterpriseTrial');
import EnterpriseTrialForm from './EnterpriseTrialForm';
import EnterpriseLoggedOutPage from './EnterpriseLoggedOutPage';
import { PageHeader } from 'dux';
import Card, { Block } from '@dux/element-card';
const { string } = PropTypes;
import styles from './EnterpriseTrial.css';

class TrialPanel extends Component {
  render() {
    return (
      <Card>
        <Block>
          <div className={styles.pricingPanel}>
            <img src="/public/images/logos/docker-logo-text.svg" alt='Docker Logo' className={styles.logo} />
            <h3>What's Included</h3>
            <p>Docker Trusted Registry</p>
            <hr />
            <p>Docker Universal Control Plane</p>
            <hr />
            <p>Commercially Supported Docker Engine</p>
          </div>
        </Block>
      </Card>
    );
  }
}

export default class EnterpriseTrial extends Component {
  static propTypes = {
    JWT: string
  }

  render() {
    const { JWT } = this.props;

    if(!JWT) {
      return <EnterpriseLoggedOutPage type='trial' />;
    } else {
      return (
        <div>
          <PageHeader title='Register for trial of Docker Datacenter' />
          <div className={styles.trialPageWrapper}>
            <div className='row'>
              <div className='large-8 columns'>
                <EnterpriseTrialForm JWT={JWT} />
              </div>
              <div className='large-4 columns'>
                <TrialPanel />
              </div>
            </div>
          </div>
        </div>
      );
    }
  }
}
