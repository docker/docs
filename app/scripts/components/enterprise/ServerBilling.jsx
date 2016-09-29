'use strict';

import React, {
  createClass,
  PropTypes
} from 'react';
const { string, func, shape } = PropTypes;
import merge from 'lodash/object/merge';
import { Link } from 'react-router';
import connectToStores from 'fluxible-addons-react/connectToStores';
var debug = require('debug')('ServerBilling');

import EnterprisePartnerTrackingStore from 'stores/EnterprisePartnerTrackingStore';
import Card, { Block } from '@dux/element-card';
import { PageHeader } from 'dux';
import createNewLicenseProduction from '../../actions/createNewLicenseProduction';
import EnterpriseBillingForm from './EnterpriseBillingForm.jsx';
import EnterpriseLoggedOutPage from './EnterpriseLoggedOutPage.jsx';
import styles from './ServerBilling.css';

var EnterprisePaid = createClass({
  displayName: 'ServerStarter',
  propTypes: {
    JWT: string.isRequired,
    user: shape({
      username: string.isRequired
    }).isRequired,
    partnervalue: string.isRequired
  },
  contextTypes: {
    executeAction: func.isRequired
  },
  submitAction(username, values) {
    this.context.executeAction(createNewLicenseProduction, merge({},
      {
        JWT: this.props.JWT,
        username,
        package_name: 'Starter',
        partnervalue: this.props.partnervalue
      },
      values));
  },
  render() {
    if(!this.props.JWT) {
      return (<EnterpriseLoggedOutPage type='server'/>);
    } else {
      return (
        <div>
          <PageHeader title='Purchase Server Starter Subscription:' />
          <div className='row'>
            <div className={'large-8 columns ' + styles.formBody}>
              <EnterpriseBillingForm submitAction={this.submitAction}
                                     enterpriseType='server'
                                     {...this.props} />
            </div>
            <div className={'large-4 columns ' + styles.formBody}>
              <Card>
                <Block>
                  <div className={styles.pricingPanel}>
                    <h3>Starter Edition</h3>
                    <h2 className={styles.price}>$150/month</h2>
                    <hr/>
                    <p>1 Docker Trusted Registry for Server</p>
                    <hr/>
                    <p>10 Docker Engines</p>
                    <hr/>
                    <p>Email Support</p>
                    <hr/>
                    <small>Max of 1 Docker Trusted Registry for Server and 10 Engines</small>
                  </div>
                </Block>
              </Card>
            </div>
          </div>
        </div>
      );
    }
  }
});

export default connectToStores(EnterprisePaid,
  [
    EnterprisePartnerTrackingStore
  ],
  function({ getStore }, props) {
    return getStore(EnterprisePartnerTrackingStore).getState();
  });
