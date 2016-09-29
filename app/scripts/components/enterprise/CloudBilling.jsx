'use strict';

import React, {
  createClass,
  PropTypes
  } from 'react';
const { string, number: propNum, object, shape, bool, func } = PropTypes;
import has from 'lodash/object/has';
import merge from 'lodash/object/merge';
import { Link } from 'react-router';
import connectToStores from 'fluxible-addons-react/connectToStores';
var debug = require('debug')('CloudBilling');

import { PageHeader } from 'dux';
import Card, { Block } from '@dux/element-card';
import CloudBillingStore from 'stores/CloudBillingStore'; // Clone of BillingPlansStore for now
import CloudCouponStore from 'stores/CloudCouponStore';
import cloudNamespaceChange from '../../actions/cloudNamespaceChange.js';
import createBillingSubscription from '../../actions/createBillingSubscription.js';
import updateCloudBillingSubscription from '../../actions/updateCloudBillingSubscription.js';
import CouponForm from './CouponForm.jsx';
import EnterpriseBillingForm from './EnterpriseBillingForm.jsx';
import EnterpriseLoggedOutPage from './EnterpriseLoggedOutPage.jsx';
import styles from './CloudBilling.css';

var CloudStarter = createClass({
  displayName: 'CloudStarter',
  propTypes: {
    JWT: string.isRequired,
    user: shape({username: string.isRequired}).isRequired,
    currentPlan: object,
    billingInfo: object,
    accountInfo: object,
    coupon: shape({
      discountValue: propNum,
      couponCode: string,
      hasError: bool
    })
  },
  contextTypes: {
    executeAction: func.isRequired
  },
  orgAction(namespace) {
    // This action updates the BillingPlansStore with billing info given user/org namespace
    this.context.executeAction(cloudNamespaceChange, {JWT: this.props.JWT, namespace});
  },
  submitAction(username, values) {
    const { JWT, currentPlan } = this.props;

    let userObject = {};
    let isOrg = false;
    if (username === this.props.user.username) {
      userObject.username = username;
    } else {
      userObject.orgname = username;
      isOrg = true;
    }

    if (has(currentPlan, 'plan') && currentPlan.plan) {
      // User HAS a hub plan - upgrade to hub plan WITH cloud subscription
      let updatePlanInfo = {
        JWT,
        username,
        subscription_uuid: currentPlan.subscription_uuid,
        package_code: 'cloud_starter',
        coupon_code: this.props.coupon.couponCode,
        isOrg
      };
      this.context.executeAction(updateCloudBillingSubscription, updatePlanInfo);
    } else {
      // User has NO hub plan or account - create account & subscription
      const {
        first_name,
        last_name,
        postal_code,
        number,
        month,
        year,
        cvv,
        address1,
        address2,
        city,
        state,
        country,
        last_four,
        card_type,
        account_first,
        account_last,
        company_name,
        email
        } = values;
      const billingInfo = {
        first_name,
        last_name,
        zip: postal_code,
        address1,
        address2,
        city,
        state,
        country
      };
      const card = {
        number,
        cvv,
        month,
        year,
        coupon_code: this.props.coupon.couponCode
      };
      const accountInfo = {
        first_name: account_first,
        last_name: account_last,
        company_name,
        email
      };
      const planInfo = {
        JWT,
        user: userObject,
        accountInfo,
        billingInfo,
        card,
        package_code: 'cloud_starter'
      };
      this.context.executeAction(createBillingSubscription, planInfo);
    }
  },
  render() {
    if(!this.props.JWT) {
      /**
       * IF not logged in - should redirect to login page.
       * OR have public plans page.
       */
      return (<EnterpriseLoggedOutPage type='cloud'/>);
    } else {
      return (
        <div>
          <PageHeader title='Purchase Cloud Starter Subscription:' />
          <div className='row'>
            <div className={'large-8 columns ' + styles.formBody}>
              <EnterpriseBillingForm submitAction={this.submitAction}
                                     orgAction={this.orgAction}
                                     enterpriseType='cloud'
                                     {...this.props} />
            </div>
            <div className={'large-4 columns'}>
              <div className={styles.formBody}>
                <Card>
                  <Block>
                    <div className={styles.pricingPanel}>
                      <h3>Starter Edition: Cloud</h3>
                      <h2 className={styles.price}>$150/month</h2>
                      <hr/>
                      <p>20 Private Repos (Medium Plan)</p>
                      <hr/>
                      <p>10 Docker Engines</p>
                      <hr/>
                      <p>Email Support</p>
                      <hr/>
                      <small>Max of 10 Docker Engines</small>
                    </div>
                  </Block>
                </Card>
              </div>
              <div className={styles.formBody}>
                <CouponForm coupon={this.props.coupon}/>
              </div>
            </div>
          </div>
        </div>
      );
    }
  }
});
/**
 * To BillingPlansStore or create separate CloudBillingStore...
 */
export default connectToStores(CloudStarter,
  [
    CloudBillingStore,
    CloudCouponStore
  ],
  function({ getStore }, props) {
    return merge({},
      getStore(CloudBillingStore).getState(),
      {coupon: getStore(CloudCouponStore).getState()}
    );
  });
