'use strict';

import React, { PropTypes } from 'react';
const { string, number, func, shape, object, array, oneOfType } = PropTypes;
import map from 'lodash/collection/map';
import find from 'lodash/collection/find';
import includes from 'lodash/collection/includes';
import has from 'lodash/object/has';
import merge from 'lodash/object/merge';
import classnames from 'classnames';
import { Link } from 'react-router';
import connectToStores from 'fluxible-addons-react/connectToStores';

import BillingInfoForm from './billingplans/BillingInfoForm.jsx';
import BillingInfoFormStore from '../../stores/BillingInfoFormStore.js';
import PlansStore from '../../stores/PlansStore.js';
import BillingPlansStore from '../../stores/BillingPlansStore.js';
import DUXInput from '../common/DUXInput.jsx';
import createBillingSubscription from '../../actions/createSubscription.js';
import updateSubscriptionPlanOrPackage from '../../actions/updateSubscriptionPlanOrPackage.js';
import updateBillingInfoFormField from '../../actions/updateBillingInfoFormField.js';
import validateBillingInfo from '../../actions/common/validateBillingInfo.js';
import validateCouponCode from '../../actions/validateCouponCode.js';
import EnterpriseLoggedOutPage from '../enterprise/EnterpriseLoggedOutPage.jsx';
import { Button, PageHeader } from 'dux';
import styles from './CreateBillingSubscription.css';

var debug = require('debug')('createBillingSubscription');

function mkShortPlanIntervalUnit(unit) {
  if (unit === 'months'){
    return 'mo';
  } else if (unit === 'years') {
    return 'yr';
  } else {
    return 'mo';
  }
}

function _mkOptions(list_item) {
  return (
    <option key={list_item.plan_code} value={list_item.plan_code}>{list_item.slug}</option>
  );
}

var CreateBillingSubscription = React.createClass({
  displayName: 'CreateBillingSubscription',
  contextTypes: {
    getStore: func.isRequired,
    executeAction: func.isRequired
  },
  getInitialState: function() {
    return {
      couponCode: '',
      selectedPlan: this.props.location.query.plan
    };
  },
  propTypes: {
    JWT: string.isRequired,
    user: object.isRequired,
    billingInfoForm: shape({
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
        type: string,
        coupon_code: string,
        coupon: number
      }),
      errorMessage: string,
      fieldErrors: object,
      STATUS: string
    }),
    plans: shape({
      currentPlan: shape({
        subscription_uuid: string,
        package: string
      }),
      plansList: array
    })
  },
  createBillingSubscription(){
    const { JWT, user } = this.props;
    const { billforwardId, accountInfo, billingInfo, card } = this.props.billingInfoForm;
    const { currentPlan } = this.props.plans;
    const { package: currentPackage, subscription_uuid } = currentPlan;
    const { selectedPlan } = this.state;
    const { executeAction } = this.context;
    if (selectedPlan) {
      // User has no billing profile account OR no billing payment information
      var subscriptionData = {
        JWT,
        user,
        accountInfo,
        billingInfo,
        card,
        billforwardId,
        isNewBillingAccount: accountInfo.newBilling,
        plan_code: selectedPlan
      };

      return executeAction(createBillingSubscription, subscriptionData);
    } else if (currentPackage) {
      // User HAS a CLOUD plan - upgrade to hub plan WITH cloud subscription
      // SHOULD NEVER REACH HERE ANYMORE - PACKAGES HAVE BEEN REMOVED!!
      const namespace = user.username || user.orgname;
      let updatePlanInfo = {
        JWT,
        username: namespace,
        subscription_uuid: subscription_uuid,
        plan_code: selectedPlan,
        package_code: currentPackage,
        coupon_code: card.coupon_code
      };
      return executeAction(updateSubscriptionPlanOrPackage, updatePlanInfo);
    }
    return executeAction(validateBillingInfo({storePrefix: 'BILLING'}));
  },
  getPlan(selectedPlan) {
    const plans = [
      'index_personal_micro',
      'index_personal_small',
      'index_personal_medium',
      'index_personal_large',
      'index_personal_xlarge',
      'index_personal_xxlarge'
    ];
    if (selectedPlan && includes(plans, selectedPlan)) {
      const { plansList } = this.props.plans;
      const planObject = find(plansList, {plan_code: selectedPlan});
      const interval = has(planObject, 'plan_interval_unit') ? mkShortPlanIntervalUnit(planObject.plan_interval_unit) : 'mo';
      const price = parseInt(planObject.price_in_cents, 10) / 100;
      return {planCode: planObject.plan_code, price: price, interval};
    } else {
      return {price: 0, interval: 'mo', planCode: null};
    }
  },
  _onChange(field, fieldKey) {
    return (e) => {
      this.context.executeAction(updateBillingInfoFormField, {
        field,
        fieldKey,
        fieldValue: e.target.value
      });
    };
  },
  _updateSelectPlan(e) {
    e.preventDefault();
    let planCode = e.target.value;
    const { card } = this.props.billingInfoForm;
    this.setState({selectedPlan: planCode});
    this.props.history.pushState(null, this.props.location.pathname, {plan: planCode});
    if (card.coupon > 0 || this.state.couponCode) {
      this.context.executeAction(validateCouponCode, {coupon_code: this.state.couponCode, plan: planCode});
    }
  },
  _updateCoupon(e) {
    e.preventDefault();
    this.setState({couponCode: e.target.value});
  },
  validateCoupon(plan) {
    return (e) => {
      e.preventDefault();
      this.context.executeAction(validateCouponCode, {coupon_code: this.state.couponCode, plan: plan});
      this.setState({couponCode: ''});
    };
  },
  render: function() {
    const { plansList } = this.props.plans;
    const { selectedPlan } = this.state;
    const { JWT, user, history } = this.props;

    if(!JWT) {
      //get plan type from end of '?plan=index_personal_PLANTYPE'
      //default to micro plan if no query selected
      const planType = selectedPlan ? selectedPlan.substring(15) : 'micro';
      return <EnterpriseLoggedOutPage type={planType} />;
    } else {
      const { accountInfo, billingInfo, card, fieldErrors, STATUS, errorMessage } = this.props.billingInfoForm;
      const { price, interval, planCode } = this.getPlan(selectedPlan);

      var discount;
      if (card.coupon > 0) {
        discount = '- $' + card.coupon;
      } else {
        discount = '-----';
      }
      let titleClass = classnames({
        [styles.title]: true,
        [styles.error]: true
      });
      let title = (
        <div className={titleClass}>
          * Please select a plan:&nbsp;
          <select className={styles.select}
                  value='---'
                  onChange={this._updateSelectPlan}>
            <option value='none'>----</option>
            {map(plansList, _mkOptions)}
          </select>
        </div>
      );

      if (selectedPlan) {
        title = (
          <div className={styles.title}>
            You are subscribing to the &nbsp;
            <select className={styles.select}
                    value={planCode}
                    onChange={this._updateSelectPlan}>
              {map(plansList, _mkOptions)}
            </select>
            &nbsp; plan at ${price}/{interval}
          </div>
        );
      }
      const namespace = user.username || user.orgname;
      const isOrg = !!user.orgname;

      return (
        <div>
          <PageHeader title={'Create Plan Subscription (' + namespace + ')'} />

            <div className="row">
              <div className="columns large-8 large-offset-2 end">
                { title }
                <div>Once your billing information has been processed, your account will be immediately upgraded.</div>
                <div>Thank you for subscribing!</div>
                <div className={styles.subtitle}>Your billing information can be updated at any time</div>
              </div>
            </div>
            <div className="row">
              <div className='columns large-6 large-offset-2 end'>
                <BillingInfoForm submitAction={this.createBillingSubscription}
                                 accountInfo={accountInfo}
                                 billingInfo={billingInfo}
                                 card={card}
                                 fieldErrors={fieldErrors}
                                 isOrg={isOrg}
                                 STATUS={STATUS}
                                 username={namespace}
                                 errorMessage={errorMessage}
                                 history={history}
                  />
              </div>
              <div className={'columns large-3 left ' + styles.preview}>
                <div className="row">
                  <div className="columns large-5">
                    Plan Cost:
                  </div>
                  <div className={'columns large-4 end ' + styles.price}>
                    ${price}
                  </div>
                </div>
                <div className="row">
                  <div className="columns large-5">
                    Coupon:
                  </div>
                  <div className={'columns large-4 end ' + styles.price}>
                    {discount}
                  </div>
                </div>
                <hr />
                <div className={'row ' + styles.total}>
                  <div className="columns large-5">
                    Total Charge:
                  </div>
                  <div className={'columns large-4 end ' + styles.price}>
                    ${price - card.coupon}
                  </div>
                </div>
                <div className="row">
                  <form className={styles.couponCode} onSubmit={this.validateCoupon(selectedPlan)}>
                    <div className="columns large-8">
                      <DUXInput label='Coupon Code'
                            hasError={fieldErrors.coupon_code}
                            error='Invalid Coupon'
                            onChange={this._updateCoupon}
                            value={this.state.couponCode}/>
                    </div>
                    <div className="columns large-4">
                      <Button type="submit" size='small'>Add</Button>
                    </div>
                  </form>
                </div>
              </div>
            </div>
        </div>
      );
    }
  }
});

export default connectToStores(CreateBillingSubscription,
  [BillingInfoFormStore, PlansStore],
  function({ getStore }, props) {
    return merge({},
      {plans: getStore(PlansStore).getState(), billingInfoForm: getStore(BillingInfoFormStore).getState()});
  });
