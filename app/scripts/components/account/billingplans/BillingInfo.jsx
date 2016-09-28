'use strict';

import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import Button from '@dux/element-button';
import { SplitSection } from '../../common/Sections.jsx';
import FA from '../../common/FontAwesome.jsx';
import styles from './BillingInfo.css';

var BillingInfo = React.createClass({
  propTypes: {
    currentPlan: PropTypes.shape({
      id: PropTypes.string,
      plan: PropTypes.string
    }),
    accountInfo: PropTypes.shape({
      account_code: PropTypes.string,
      username: PropTypes.string,
      email: PropTypes.string,
      first_name: PropTypes.string,
      last_name: PropTypes.string,
      company_name: PropTypes.string
    }),
    billingInfo: PropTypes.shape({
      city: PropTypes.string,
      state: PropTypes.string,
      zip: PropTypes.string,
      first_name: PropTypes.string,
      last_name: PropTypes.string,
      address1: PropTypes.string,
      address2: PropTypes.string,
      country: PropTypes.string
    }),
    invoices: PropTypes.array,
    isOrg: PropTypes.bool.isRequired,
    username: PropTypes.string.isRequired
  },
  _updateClick: function(e) {
    e.preventDefault();
    if (this.props.isOrg) {
      this.props.history.pushState(null, `/u/${this.props.username}/dashboard/billing/update-info/`);
    } else {
      this.props.history.pushState(null, '/account/billing-plans/update/');
    }
  },
  render: function() {
    const {
      accountInfo,
      billingInfo,
      currentPlan
    } = this.props;
    if (accountInfo.newBilling) {
      return (
        <div></div>
      );
    }
    let subtitle = (
      <div>
        <div className="account-section-text">
          <p>* fields required to complete a billing transaction</p>
        </div>
        <div>
          <div className="columns large-6">
            <Button size='tiny'
                    onClick={this._updateClick}>Update Billing Info</Button>
          </div>
        </div>
      </div>
    );
    let cardInfo;
    if (billingInfo.card_type && billingInfo.last_four) {
      cardInfo = (
        <div className="row">
          <div className="columns large-4">Card Info:</div>
          <div className={'columns large-8 ' + styles.infoContent}>
            {billingInfo.first_name} {billingInfo.last_name}<br/>
            {billingInfo.card_type} card ending with x{billingInfo.last_four} <br/>
            Expiration: {billingInfo.month}/{billingInfo.year}
          </div>
        </div>
      );
    }
    let addressInfo;
    if (billingInfo.address1 && billingInfo.country) {
      addressInfo = (
        <div className="row">
          <div className="columns large-4">Billing Address:</div>
          <div className={'columns large-8 ' + styles.infoContent}>
            {billingInfo.address1} <br/>
            {billingInfo.address2} <br/>
            {billingInfo.city} {billingInfo.state} {billingInfo.zip} <br/>
            {billingInfo.country}
          </div>
        </div>
      );
    }
    return (
      <SplitSection title='Account Billing Information'
                    subtitle={subtitle}>
          <div className="row">
            <div className="columns large-4">
              <div>Name:</div>
            </div>
            <div className={'columns large-8 ' + styles.infoContent}>
              <div>{accountInfo.first_name} {accountInfo.last_name}</div>
            </div>
          </div>
          <div className="row">
            <div className="columns large-4">
              <div>Email:</div>
            </div>
            <div className={'columns large-8 ' + styles.infoContent}>
              <div>{accountInfo.email}</div>
            </div>
          </div>
          <div className="row">
            <div className="columns large-4">
              <div>Company:</div>
            </div>
            <div className={'columns large-8 ' + styles.infoContent}>
              <div>{accountInfo.company_name}</div>
            </div>
          </div>
          { cardInfo }
          { addressInfo }
      </SplitSection>
    );
  }
});

module.exports = BillingInfo;
