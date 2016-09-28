/*global MktoForms2*/

'use strict';

import React, { PropTypes } from 'react';
const { string, array, object, func, shape } = PropTypes;
import { Link } from 'react-router';
import connectToStores from 'fluxible-addons-react/connectToStores';
import isEmpty from 'lodash/lang/isEmpty';

import BillingPlansStore from '../../stores/BillingPlansStore';
import updateSubscriptionPlanOrPackage from '../../actions/updateSubscriptionPlanOrPackage.js';

import PlansTable from './billingplans/Plans';
import EnterpriseSubscriptions from './billingplans/EnterpriseSubscriptions.jsx';
import BillingInfo from './billingplans/BillingInfo.jsx';
import InvoiceTables from './billingplans/InvoiceTables.jsx';
import styles from './BillingPlans.css';
import { FullSection } from '../common/Sections.jsx';
import Route404 from '../common/RouteNotFound404Page.jsx';
import { PageHeader } from 'dux';

/* Marketo constants for the marketing survey form */
const mktoFormId = 1317;
const mktoFormElemId = 'mktoForm_' + mktoFormId;
const mktoFormBaseUrl = 'https://app-sj05.marketo.com';
const mktoFormMunchkinId = '929-FJL-178';

var BillingInfoPage = React.createClass({

  contextTypes: {
    executeAction: func.isRequired
  },

  PropTypes: {
    JWT: string,
    user: object,
    currentPlan: shape({
      id: string,
      plan: string,
      package: string
    }),
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
    plansError: string,
    invoices: array,
    unsubscribing: string,
    updatePlan: string
  },

  stopSubscription(subscriptionType) {
    /**
     * UPDATE 4/6/16 "cloud_metered" is the new "free" plan instead of deleting
     * per ticket HUB-2219
     */
    return () => {
      const { JWT, user, currentPlan } = this.props;
      const namespace = user.username || user.orgname;
      let subscriptionData = {
        JWT,
        username: namespace,
        subscription_uuid: currentPlan.subscription_uuid
      };
      if (subscriptionType === 'plan') {
        subscriptionData.plan_code = 'cloud_metered';
        if (currentPlan.package) {
          // preserve package (like cloud_starter) if it exists
          subscriptionData.package_code = currentPlan.package;
        }
      } else if (subscriptionType === 'package') {
        // If you are removing a package, leave the plan alone
        subscriptionData.plan_code = currentPlan.plan;
        // Explicitly set null to remove
        subscriptionData.package_code = null;
      }
      this.context.executeAction(updateSubscriptionPlanOrPackage, subscriptionData);
    };
  },

  showSurveyModal(subscriptionType) {
    return () => {
      if (typeof MktoForms2 === 'object' &&
          typeof MktoForms2.loadForm === 'function') {
        MktoForms2.loadForm(
          mktoFormBaseUrl,
          mktoFormMunchkinId,
          mktoFormId,
          (form) => {
            form.onSubmit(this.stopSubscription(subscriptionType));
            // Don't refresh the page after a successful submission.
            // React component will re-render itself.
            form.onSuccess(() => false);
            form.vals({'Email': this.props.accountInfo.email});
            MktoForms2.lightbox(form).show();
          });
      } else {
        // If for any reason there is a problem with the Marketo script,
        // we shouldn't block the user from stopping his/her subscription.
        this.stopSubscription(subscriptionType);
      }
    };
  },

  getSurveyModalHtml() {
    return (
      <div>
        <form id={mktoFormElemId}></form>
      </div>
    );
  },

  render: function() {
    let plansIntro;
    let plansFooter;
    let errorSection;
    const {
      accountInfo,
      billingInfo,
      currentPlan,
      invoices,
      isOwner,
      JWT,
      plansError,
      unsubscribing,
      updatePlan,
      user,
      history
    } = this.props;
    if (accountInfo.newBilling) {
      plansIntro = (
        <FullSection title='Plans and Pricing'>
          <div className="columns large-10 end">
            The Docker Hub Registry is free to use for public repositories. Plans with private repositories are
            available in different sizes. All plans allow collaboration with unlimited people.
          </div>
        </FullSection>
      );
      plansFooter = (
        <FullSection title='Questions'>
          <div className={'columns large-6 ' + styles.plansQuestion}>
            <div className={styles.questionTitle}>What types of payment do you accept?</div>
            <div className={styles.questionAnswer}>Credit card (Visa, MasterCard, Discover, or American Express).</div>
          </div>
          <div className={'columns large-6 ' + styles.plansQuestion}>
            <div className={styles.questionTitle}>Do I have to pay to use your service?</div>
            <div className={styles.questionAnswer}>No, you only have to pay if you require one or more private repository.</div>
          </div>
          <div className={'columns large-6 ' + styles.plansQuestion}>
            <div className={styles.questionTitle}>Can I change my plan at a later time?</div>
            <div className={styles.questionAnswer}>Yes, you can upgrade or downgrade at any time.</div>
          </div>
          <div className={'columns large-6 ' + styles.plansQuestion}>
            <div className={styles.questionTitle}>What if I need a larger plan?</div>
            <div className={styles.questionAnswer}>Please contact our Sales team at <a href="mailto:sales@docker.com">sales@docker.com</a> or call us toll free at 888-214-4258.</div>
          </div>
        </FullSection>
      );
    } else {
      plansIntro = (<div></div>);
      plansFooter = (<div></div>);
    }
    if (plansError) {
      errorSection = (
        <FullSection>
          <div className={'columns large-12 ' + styles.error}>
            There was an error trying to update your subscription. Please contact the Docker support team.<br/>
            {plansError}
          </div>
        </FullSection>
      );
    }
    const username = user.username || user.orgname || '';
    const isOrg = !!user.orgname;

    if (isOrg && !isOwner) {
      return (
        <Route404 />
      );
    } else {
      /*
        accountInfo.newBilling / billingInfo.newBilling
        Means this is a brand new account without any information saved to the backend
      */
      return (
        <div>
          <PageHeader title={'Billing Information & Pricing Plans: ' + username }/>
          <div className={styles.body}>
            {plansIntro}
            {errorSection}
            <PlansTable JWT={JWT}
                        user={user}
                        history={history}
                        username={username}
                        isOrg={isOrg}
                        isNewBilling={accountInfo.newBilling || billingInfo.newBilling}
                        billingInfo={billingInfo}
                        updatePlan={updatePlan}
                        stopSubscription={this.showSurveyModal('plan')}/>
            {plansFooter}<br/>
            <EnterpriseSubscriptions JWT={JWT}
                                     user={user}
                                     history={history}
                                     currentPlan={currentPlan}
                                     unsubscribing={unsubscribing}
                                     stopSubscription={this.showSurveyModal('package')}/><br/>
            <BillingInfo stopSubscription={this.showSurveyModal('plan')}
                         unsubscribing={unsubscribing}
                         username={username}
                         isOrg={isOrg}
                         currentPlan={currentPlan}
                         billingInfo={billingInfo}
                         accountInfo={accountInfo}
                         history={history}/>

            <InvoiceTables JWT={JWT}
                           username={username}
                           invoices={invoices} />
          </div>
          {this.getSurveyModalHtml()}
        </div>
      );
    }
  }
});

export default connectToStores(BillingInfoPage,
  [BillingPlansStore],
  function({ getStore }, props) {
    return getStore(BillingPlansStore).getState();
  });
