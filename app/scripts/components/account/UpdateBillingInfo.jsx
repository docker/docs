'use strict';
import React, { PropTypes } from 'react';
const { string, number, func, shape, object, array, oneOfType } = PropTypes;
import connectToStores from 'fluxible-addons-react/connectToStores';
import updateBillingInformation from '../../actions/updateBillingInformation.js';
import updateStripeBilling from '../../actions/updateStripeBilling.js';
import BillingInfoFormStore from '../../stores/BillingInfoFormStore.js';
import BillingInfoForm from './billingplans/BillingInfoForm.jsx';
import { Link } from 'react-router';
import classnames from 'classnames';
import styles from './UpdateBillingInfo.css';
import { PageHeader, Module } from 'dux';
var debug = require('debug')('UpdateBillingInfo:');

var updateBillingInfoPage = React.createClass({
  displayName: 'UpdateBillingInfo',
  contextTypes: {
    getStore: func.isRequired,
    executeAction: func.isRequired
  },
  propTypes: {
    JWT: string.isRequired,
    user: object.isRequired,
    billforwardId: string,
    accountInfo: shape({
      account_code: string,
      username: string,
      email: string,
      first_name: string,
      last_name: string,
      company_name: string
    }),
    billingInfo: shape({
      city: string,
      state: string,
      zip: string,
      first_name: string,
      last_name: string,
      address1: string,
      address2: string,
      country: string
    }),
    card: shape({
      number: string,
      cvv: string,
      month: oneOfType([number, string]),
      year: oneOfType([number, string]),
      type: string
    }),
    errorMessage: string,
    fieldErrors: object,
    STATUS: string
  },
  updateBillingInfoSubmit(){
    const {
      JWT,
      user,
      accountInfo,
      billingInfo,
      card,
      billforwardId,
      location
      } = this.props;
    if (billforwardId) {
      // If we have a billforwardID then we will update via stripe
      this.context.executeAction(updateStripeBilling, {
        JWT,
        user,
        accountInfo,
        billingInfo,
        billforwardId,
        card
      });
    } else {
      // If we do not have a billforward ID then we go through the original recurly update flow
      this.context.executeAction(updateBillingInformation, {
        JWT,
        user,
        accountInfo,
        billingInfo,
        card
      });
    }
  },
  render: function() {
    const {
      user,
      history,
      accountInfo,
      billingInfo,
      card,
      fieldErrors,
      STATUS,
      errorMessage
      } = this.props;
    const namespace = user.username || user.orgname;
    let isOrg = !!user.orgname;
    return (
      <div>
        <PageHeader title={'Update Billing Information: ' + namespace} />
        <div className={'row ' + styles.body}>
          <div className="columns large-6 large-offset-3 end">
            <h5>Billing information is required for changing or upgrading subscriptions</h5>
            <div className={styles.subtitle}>You may update your billing information at any time</div>
          </div>
        </div>
        <div className='row'>
          <div className='colummns large-6 large-offset-3 end'>
            <BillingInfoForm submitAction={this.updateBillingInfoSubmit}
                             accountInfo={accountInfo}
                             billingInfo={billingInfo}
                             card={card}
                             fieldErrors={fieldErrors}
                             isOrg={isOrg}
                             STATUS={STATUS}
                             errorMessage={errorMessage}
                             username={namespace}
                             history={history}
                             />
          </div>
        </div>
      </div>
    );
  }
});

export default connectToStores(updateBillingInfoPage,
  [BillingInfoFormStore],
  function({ getStore }, props) {
    return getStore(BillingInfoFormStore).getState();
  });
