import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { sortBy, has, find, cloneDeep, isEmpty } from 'lodash';
import TrialForm from './TrialForm';
import SignupLoginForm from './SignupLoginForm';
import PurchaseForm from './PurchaseForm';
import PurchaseDetails from './PurchaseDetails';
import {
  BackButtonArea,
  FullSectionLoading,
  DDCGraphic,
  ImageWithFallback,
} from 'common';
import { DDC_ID, DDC_TRIAL_PLAN, DDC_BUSINESS_DAY } from 'lib/constants/eusa';
import { CANCELLED } from 'lib/constants/states/subscriptions';
import { FALLBACK_IMAGE_SRC, FALLBACK_ELEMENT } from 'lib/constants/fallbacks';
import getLargestLogo from 'lib/utils/get-largest-logo';
import classNames from 'classnames';
import css from './styles.css';
import routes from 'lib/constants/routes';
const { array, arrayOf, bool, object, node, shape, func } = PropTypes;

const mapStateToProps = ({ account, billing, marketplace }, { params }) => {
  const { id } = params;
  const {
    currentUser,
    userEmails,
    ownedNamespaces,
  } = account;
  const {
    products: billingProducts,
    subscriptions,
    paymentMethods,
    profiles,
  } = billing;
  const { bundles } = marketplace;
  const formSubmitting =
    subscriptions.isSubmitting ||
    paymentMethods.isSubmitting ||
    profiles.isSubmitting;

  const billingProductDetails = billingProducts[id];
  let sortedPlans = [];
  if (billingProductDetails && has(billingProductDetails, 'rate_plans')) {
    sortedPlans = sortBy(billingProductDetails.rate_plans, (plan) => {
      return plan.name !== DDC_BUSINESS_DAY;
    });
  }
  return {
    billingProductDetails,
    currentUser,
    formSubmitting,
    ownedNamespaces,
    // details coming from the product catalog
    productDetails: bundles && bundles[id],
    sortedPlans,
    subscriptions,
    userEmails: userEmails.results,
  };
};

@connect(mapStateToProps)
export default class DDCPurchase extends Component {
  static propTypes = {
    billingProductDetails: shape({
      rate_plans: arrayOf(object),
    }),
    children: node,
    currentUser: object,
    formSubmitting: bool.isRequired,
    location: object.isRequired,
    history: object.isRequired,
    params: object.isRequired,
    ownedNamespaces: array,
    productDetails: object,
    userEmails: array,
    sortedPlans: arrayOf(object),
    subscriptions: shape({
      results: object,
      isFetching: bool.isRequired,
      isSubmitting: bool.isRequired,
    }),
  }

  static contextTypes = {
    router: shape({
      push: func.isRequired,
    }).isRequired,
  }

  constructor(props) {
    super(props);
    const query = props.location.query;
    this.state = {
      selectedPlan: query.plan,
      components: {},
    };
  }

  componentDidMount() {
    const { location, billingProductDetails } = this.props;
    if (!location.query.plan && has(billingProductDetails, 'rate_plans')) {
      // Get the Business Day ddc plan
      // NOTE: only DDC may be purchased right now. This logic
      // will have to be changed once we support other purchasable products
      const ddc_bd_plan = find(billingProductDetails.rate_plans, (plan) => {
        return plan.name === DDC_BUSINESS_DAY;
      });

      const plan = ddc_bd_plan && ddc_bd_plan.name;
      this.context.router.replace({
        pathname: location.pathname,
        query: { plan },
      });
    }
  }

  getNamespaces() {
    const { ownedNamespaces } = this.props;
    return ownedNamespaces.map(
      account => ({ label: account, value: account })
    );
  }

  getEmails() {
    // ORGS HAVE NO EMAILS - ONLY GETS CURRENT USER EMAILS
    const { userEmails } = this.props;
    // sortBy in ascending order puts false first and true last. ! to reverse
    const emailList = sortBy(userEmails, o => !o.primary);
    return emailList.map((emailObject) => {
      return emailObject.email;
    });
  }

  changePlan = (planName) => {
    this.setState({
      selectedPlan: planName,
    });
  }

  updateComponents = (componentName) => (value) => {
    const newState = cloneDeep(this.state.components);
    newState[componentName] = value;
    this.setState({ components: newState });
  }

  renderDetails({ selectedPlanObject }) {
    const {
      billingProductDetails,
      formSubmitting,
      productDetails, // details coming from the product api
    } = this.props;

    let details = (
      <PurchaseDetails
        {...this.state}
        billingProductDetails={billingProductDetails}
        formSubmitting={formSubmitting}
        location={location}
        onSubmit={this.submitForm}
        selectedPlanObject={selectedPlanObject}
        productDetails={productDetails}
      />
    );
    // If the product has a trial plan - render trial flow
    if (selectedPlanObject.name === DDC_TRIAL_PLAN) {
      const logoUrl = getLargestLogo(productDetails.logo_url);
      const logo = logoUrl ? (
        <ImageWithFallback
          src={logoUrl}
          className={css.logo}
          fallbackImage={FALLBACK_IMAGE_SRC}
          fallbackElement={FALLBACK_ELEMENT}
        />
      ) : (
        <DDCGraphic
          size="xlarge"
          className={css.ddcgraphic}
          variant="dull"
        />
      );

      // current trial flow hardcoded to DDC graphic icon && descriptions
      details = (
        <div className={css.details}>
          <div className={css.header}>
            <div className={css.logoWrapper}>
              {logo}
            </div>
            <div className={css.headerText}>
              <div className={css.title}>Docker Datacenter</div>
              <div className={css.tagline}>
                Docker Datacenter delivers container management and deployment
                services to the enterprise via a Containers as a Service
                solution that is supported by Docker and hosted locally behind
                the enterprise firewall.
              </div>
            </div>
          </div>
          <div className={css.value}>
            <div className={css.title}>Benefits</div>
            <ul
              className={classNames({
                [css.bullets]: true,
                [css.check]: true,
              })}
            >
              <li>
                <strong>Simple:</strong> Install quickly and utilize a
                powerful user-friendly web based GUI
              </li>
              <li>
                <strong>Secure:</strong> Set granular role based access, sign
                images and deploy on-premises
              </li>
              <li>
                <strong>Scalable:</strong> Orchestrate highly available
                Dockerized workloads anywhere, at scale
              </li>
              <li>
                <strong>Supported:</strong> Support from Docker helps you
                focus on building and deploying apps versus patching and
                fixing bugs
              </li>
            </ul>
          </div>
          <div className={css.value}>
            <div className={css.title}>Includes</div>
            <ul className={css.bullets}>
              <li>
                <strong>Universal Control Plane:</strong> Manage, deploy and
                scale dockerized applications across multiple environments,
                on-premises
              </li>
              <li>
                <strong>Docker Trusted Registry:</strong> Store and secure
                image content on-premises
              </li>
              <li>
                <strong>CS Engine:</strong> Receive enterprise-grade support
                for the Docker engine from our Support team
              </li>
            </ul>
          </div>
          <div className={css.value}>
            <div className={css.title}>Have More Questions?</div>
            <div className={css.contact}>
              <a className={css.contact}
                href="https://goto.docker.com/contact-us.html"
              >
              Contact our Sales Team
              </a>
            </div>
          </div>
        </div>
      );
    }
    return details;
  }

  renderForm({ selectedPlanObject }) {
    const {
      billingProductDetails,
      productDetails,
      subscriptions,
      location,
      currentUser,
     } = this.props;
     /*
      NOTE:productDetails comes from product API.
      BillingProductDetails comes from billing
    */

    const {
      components, // selected pricing component values;
    } = this.state;

    let form = (
      <PurchaseForm
        accountEmails={this.getEmails()}
        accountOptions={this.getNamespaces()}
        billingProductDetails={billingProductDetails}
        location={location}
        productDetails={productDetails}
        changePlan={this.changePlan}
        components={components}
        updateComponents={this.updateComponents}
        selectedPlanObject={selectedPlanObject}
      />
    );

    // If the product has a trial plan - render trial flow
    if (selectedPlanObject.name === DDC_TRIAL_PLAN) {
      const hasDDCTrial = find(subscriptions.results, (sub) => {
        return sub.product_id === DDC_ID &&
          sub.product_rate_plan === DDC_TRIAL_PLAN &&
          sub.state !== CANCELLED;
      });
      form = (
        <TrialForm
          accountEmails={this.getEmails()}
          accountOptions={this.getNamespaces()}
          billingProductDetails={billingProductDetails}
          location={location}
          productDetails={productDetails}
          DDCTrialExists={!!hasDDCTrial}
        />
      );
    }

    // If the user is not logged in then render the sign up form
    if (!currentUser || isEmpty(currentUser)) {
      form = (
        <SignupLoginForm
          location={location}
          plan={selectedPlanObject.name}
        />
      );
    }

    return form;
  }

  render() {
    const {
      formSubmitting,
      productDetails,
      params,
      sortedPlans,
    } = this.props;

    const { id } = params;
    if (!productDetails) {
      return <div>Fetching ...</div>;
    } else if (productDetails.error) {
      return <div>Product does not exist</div>;
    }

    let selectedPlanObject = {};
    /*
      TODO: design - 6/1/2016 nathan - Clean up 'plan does not exist'
    */
    if (sortedPlans.length > 0) {
      const { selectedPlan } = this.state;
      selectedPlanObject =
        find(sortedPlans, (ratePlan) => {
          return ratePlan.name === selectedPlan;
        }) || sortedPlans[0];
    }
    const { name = 'Docker Datacenter' } = productDetails;
    const formLoading = formSubmitting ?
      <FullSectionLoading
        title={'Submitting...'}
      /> : null;
    const pageClass = selectedPlanObject.name === DDC_TRIAL_PLAN ?
      css.trialPage : css.purchasePage;
    return (
      <div>
        <BackButtonArea
          pathname={routes.bundleDetail({ id })}
          text={name}
        />
        <div className={pageClass}>
          <div className={css.form}>
            {this.renderForm({ selectedPlanObject })}
            {formLoading}
          </div>
          {this.renderDetails({ selectedPlanObject })}
        </div>
      </div>
    );
  }
}
