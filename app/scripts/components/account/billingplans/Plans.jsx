'use strict';
/*eslint-disable camelcase*/
import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import PlansStore from '../../../stores/PlansStore';
import updateSubscriptionPlanOrPackage from '../../../actions/updateSubscriptionPlanOrPackage.js';
import defaultPlans from 'common/data/plans.js';
/**
* Used for the public pricing page (/billing-plans/)
* We should either find a workaround for defaultPlans or remove this route altogether
*/
import FA from '../../common/FontAwesome.jsx';
import _ from 'lodash';
import connectToStores from 'fluxible-addons-react/connectToStores';
import { FullSection } from '../../common/Sections.jsx';
import { FlexTable, FlexRow, FlexHeader, FlexItem } from '../../common/FlexTable.jsx';
import styles from './Plans.css';
const CLOUD_METERED = 'cloud_metered';
const debug = require('debug')('Plans:');
const NAUTILUS = 'nautilus';
import Tooltip from 'rc-tooltip';
import classnames from 'classnames';


function mkShortPlanIntervalUnit(unit) {
  if (unit === 'months'){
    return 'mo';
  } else if (unit === 'years') {
    return 'yr';
  } else {
    return 'mo';
  }
}

let mkPricingElement = function(plan) {
  var action;
  var shortInterval = mkShortPlanIntervalUnit(plan.plan_interval_unit);
  var price = parseInt(plan.price_in_cents, 10) / 100;
  var pricePerInterval = {price: price, interval: shortInterval};

  var currentPlanName = this.props.currentPlan.plan || 'free';
  var currentSize = this.getPlanSize(currentPlanName);
  debug(currentSize);
  if (this.props.JWT) {
    if (!!this.props.currentPlan && currentPlanName === plan.plan_code) {
      action = 'Current Plan';
    } else if (this.state.confirmPlan === plan.plan_code) {
      action = (
        <div>
          <a onClick={this.updateSubscriptionsPlan(plan)}>Confirm</a> or&nbsp;
          <a className={styles.cancel} onClick={this.cancelSelectPlan}>Cancel</a>
        </div>
      );
    } else if (this.props.updatePlan === plan.plan_code) {
      action = (
        <div>
          Processing <FA icon='fa-spinner fa-spin'/>
        </div>
      );
    } else if (!!this.props.currentPlan && currentSize > plan.num_private_repos) {
      action = (<a href='#' onClick={this.selectConfirmPlan(plan)}>Downgrade Plan</a>);
    } else {
      action = (<a href='#' onClick={this.selectConfirmPlan(plan)}>Upgrade Plan</a>);
    }
  } else {
    action = (
      <div>
        <Link to='/'>Sign up</Link> or <Link to='/login/'>Log in</Link>
      </div>
    );
  }

  return (
    <PlanRow {...plan}
    pricePerInterval={pricePerInterval}
    action={action}
    currentPlanName={currentPlanName}
    key={plan.slug}/>
  );
};

function isNautilusEnabled(currentPlan) {
  const { add_ons } = currentPlan;
  if (!add_ons) {
    return false;
  }
  return _.indexOf(add_ons, NAUTILUS) >= 0;
}


var PlanRow = React.createClass({
  displayName: 'PlanRow',
  propTypes: {
    is_active: PropTypes.bool.isRequired,
    num_private_repos: PropTypes.number.isRequired,
    name: PropTypes.string.isRequired,

    pricePerInterval: PropTypes.shape({
      price: PropTypes.number.isRequired,
      interval: PropTypes.string.isRequired
    }),
    action: PropTypes.oneOfType([
      PropTypes.string.isRequired,
      PropTypes.object.isRequired
    ])
  },
  render() {
    if (!this.props.is_active) {
      return null;
    } else {
      var ppi = this.props.pricePerInterval;
      return (
        <FlexRow>
          <FlexItem>{this.props.name}</FlexItem>
          <FlexItem>${ppi.price}/{ppi.interval}</FlexItem>
          <FlexItem>{this.props.num_private_repos}</FlexItem>
          <FlexItem>{this.props.num_private_repos}</FlexItem>
          <FlexItem>{this.props.action}</FlexItem>
        </FlexRow>
      );
    }
  }
});

var PlansTable = React.createClass({
  displayName: 'PlansTable',

  contextTypes: {
    getStore: React.PropTypes.func.isRequired,
    executeAction: React.PropTypes.func.isRequired
  },

  propTypes: {
    JWT: PropTypes.string,
    username: PropTypes.oneOfType([
      PropTypes.string, PropTypes.bool
    ]),
    history: PropTypes.object.isRequired,
    user: PropTypes.object,
    plansList: PropTypes.array.isRequired,
    currentPlan: PropTypes.shape({
      subscription_uuid: PropTypes.string,
      plan: PropTypes.string,
      package: PropTypes.string
    }),
    stopSubscription: PropTypes.func.isRequired,
    billingInfo: PropTypes.object.isRequired,
    isNewBilling: PropTypes.bool.isRequired,
    updatePlan: PropTypes.string
  },

  getInitialState() {
    return {
      confirmPlan: '',
      hasNautilus: isNautilusEnabled(this.props.currentPlan)
    };
  },

  getPlanSize(plan_code) {
    return (_.result(_.find(this.props.plansList, {plan_code: plan_code}), 'num_private_repos'));
  },

  selectConfirmPlan(plan) {
    return (e) => {
      e.preventDefault();
      this.setState({confirmPlan: plan.plan_code});
    };
  },

  cancelSelectPlan: function(e) {
    e.preventDefault();
    this.setState({confirmPlan: ''});
  },

  _stopSubscription: function(e) {
    e.preventDefault();
    this.setState({confirmPlan: ''});
    this.props.stopSubscription();
  },

  updateSubscriptionsPlan(plan) {
    return (e) => {
      e.preventDefault();
      this.setState({confirmPlan: ''});
      if (this.props.isNewBilling) {
        // IF USER HAS NO BILLING ACCOUNT OR NO BILLING PAYMENT INFORMATION
        // GO TO CREATION PAGE
        if (_.has(this.props.user, 'username')) {
          this.props.history.pushState(null, '/account/billing-plans/create-subscription/', {plan: plan.plan_code});
        } else if (_.has(this.props.user, 'orgname')) {
          this.props.history.pushState(null, `/u/${this.props.user.orgname}/dashboard/billing/create-subscription/`, {plan: plan.plan_code});
        }
      } else {
        // IF USER HAS BILLING INFORMATION - THEN WE CAN JUST UPGRADE
        this.context.executeAction(updateSubscriptionPlanOrPackage,
          {JWT: this.props.JWT,
            username: this.props.username,
            subscription_uuid: this.props.currentPlan.subscription_uuid,
            plan_code: plan.plan_code,
            package_code: this.props.currentPlan.package
          });
      }
    };
  },

  toggleNautilus() {
    const { currentPlan, username, JWT } = this.props;
    const { plan, subscription_uuid } = currentPlan;
    const { hasNautilus } = this.state;
    const add_ons = hasNautilus ? [] : [NAUTILUS];
    const data = {
      plan_code: plan,
      username,
      add_ons,
      subscription_uuid,
      JWT
    };

    this.context.executeAction(updateSubscriptionPlanOrPackage, data);
    this.setState({ hasNautilus: !hasNautilus });
  },

  renderNautilusUpsell() {
    const { currentPlan } = this.props;
    const isFreePlan = !currentPlan || !currentPlan.plan
      || currentPlan.plan === CLOUD_METERED;
    let text;
    const shield = (
      <span className={styles.shield}>
        <svg width="17px" height="17px" viewBox="1080 274 30 30" version="1.1">
            <g id="ic-security-black-24-px" stroke="none" strokeWidth="1" fill="none" fill-rule="evenodd" transform="translate(1080.000000, 274.000000)">
                <path d="M14.5,1.20833333 L3.625,6.04166667 L3.625,13.2916667 C3.625,19.9979167 8.265,26.2691667 14.5,27.7916667 C20.735,26.2691667 25.375,19.9979167 25.375,13.2916667 L25.375,6.04166667 L14.5,1.20833333 L14.5,1.20833333 Z M14.5,14.4879167 L22.9583333,14.4879167 C22.3179167,19.46625 18.995,23.9008333 14.5,25.2904167 L14.5,14.5 L6.04166667,14.5 L6.04166667,7.6125 L14.5,3.85458333 L14.5,14.4879167 L14.5,14.4879167 Z" id="Shape" fill="#8F9EA8"></path>
                <polygon id="Shape" points="0 0 29 0 29 29 0 29"></polygon>
            </g>
        </svg>
      </span>
    );
    if (isFreePlan) {
      text = 'Enable security scanning when you upgrade your plan.';
      const url = 'https://docs.docker.com/docker-cloud/builds/image-scan/';
      const link = (
        <a href={url} target="_blank" className={styles.nautilusLink}>
          Learn more about Docker Security Scanning
        </a>
      );
      return (
        <div className={styles.nautilusUpsell}>
          <div className={styles.monitoredWithNautilus}>
            {shield}
            {text}&nbsp;{link}
          </div>
        </div>
      );
    }

    var nautilusUpsellClasses = classnames({
      [styles.nautilusUpsell]: true,
      [styles.nautilusEnabled]: this.state.hasNautilus
    });

    text = [
      'Monitor with Docker Security Scanning - available for your private ',
      'repositories for free while in preview.'
    ].join('');
    return (
      <div className={nautilusUpsellClasses}>
        <label>
          <div className={styles.enableNautilus}>
            <input
              type="checkbox"
              checked={this.state.hasNautilus}
              onChange={this.toggleNautilus}
            />
          </div>
          <div className={styles.monitoredWithNautilus}>
            {shield}
            {text}
          </div>
        </label>
      </div>
    );
  },

  render() {
    const {JWT, currentPlan, updatePlan, plansList, isOrg} = this.props;
    var action;
    if (!JWT) {
      action = (
        <div>
          <Link to='/'>Sign up</Link> or <Link to='/login/'>Log in</Link>
        </div>
      );
    } else if (!currentPlan.plan || currentPlan.plan === CLOUD_METERED) {
      action = 'Current Plan';
    } else if (updatePlan === CLOUD_METERED) {
      action = (
        <div>Removing Plan <FA icon='fa-spinner fa-spin'/></div>
      );
    } else if (this.state.confirmPlan === 'free') {
      action = (
        <div>
          <a onClick={this._stopSubscription}>Confirm</a> or&nbsp;
          <a className={styles.cancel} onClick={this.cancelSelectPlan}>Cancel</a>
        </div>
      );
    } else {
      action = <a href='#' onClick={this.selectConfirmPlan({plan_code: 'free'})}>Downgrade Plan</a>;
    }
    let allPlans = defaultPlans;
    if (!_.isEmpty(plansList)) {
      allPlans = _.clone(plansList);
    }
    const privateRepos = isOrg ? '0' : '1';
    return (
      <FullSection title='Choose the Hub private repo plan that works for you.'>
        <div className="columns large-12">
          <FlexTable>
            <FlexHeader>
              <FlexItem>
                Plan
              </FlexItem>
              <FlexItem>
                Price
              </FlexItem>
              <FlexItem>
                Private Repositories
              </FlexItem>
              <FlexItem>
                Parallel Builds
              </FlexItem>
              <FlexItem/>
            </FlexHeader>
            <FlexRow>
              <FlexItem>Free</FlexItem>
              <FlexItem>$0/mo</FlexItem>
              <FlexItem>{privateRepos}</FlexItem>
              <FlexItem>{privateRepos}</FlexItem>
              <FlexItem>{action}</FlexItem>
            </FlexRow>
            {allPlans.map(mkPricingElement, this)}
            <div className={styles.nautilusUpsellRow}>
              {this.renderNautilusUpsell()}
            </div>
          </FlexTable>
        </div>
      </FullSection>
    );
  }
});

export default connectToStores(PlansTable,
                               [PlansStore],
                               function({ getStore }, props) {
                                 return getStore(PlansStore).getState();
                               });
