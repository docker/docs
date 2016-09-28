import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import {
  BackButtonArea,
  Button,
  Card,
  CheckIcon,
  CloseIcon,
  CopyPullCommand,
  DownloadIcon,
  ImageWithFallback,
  Markdown,
  Modal,
} from 'common';
import { FALLBACK_IMAGE_SRC, FALLBACK_ELEMENT } from 'lib/constants/fallbacks';
import filter from 'lodash/filter';
import find from 'lodash/find';
import map from 'lodash/map';
import startCase from 'lodash/startCase';
// Blob is a polyfill to support older browsers
require('lib/utils/blob');
import { saveAs } from 'lib/utils/file-saver';
import moment from 'moment';
import classnames from 'classnames';
import getLogo from 'lib/utils/get-largest-logo';
import { SMALL } from 'lib/constants/sizes';
import isStaging from 'lib/utils/isStaging';
import isDev from 'lib/utils/isDevelopment';
import { isPullableProduct } from 'lib/utils/product-utils';
import { DOCKER, DOCKER_GUID } from 'lib/constants/defaults';
import { DDC_TRIAL_PLAN, DDC_ID } from 'lib/constants/eusa';
import { ACTIVE, CANCELLED } from 'lib/constants/states/subscriptions';
import {
  billingDeleteSubscription,
  billingFetchLicenseFile,
  billingFetchProfileSubscriptions,
} from 'actions/billing';
import css from './styles.css';
const { arrayOf, bool, func, object, shape, string } = PropTypes;

// Has a subscription been created within the past day?
const isNewSubscription = (subscription) => {
  const YESTERDAY = moment().subtract(1, 'days').startOf('day');
  return moment(subscription.initial_period_start).isAfter(YESTERDAY);
};

const isExpiredLicense = (expiration) => {
  return moment().isAfter(expiration);
};

const dispatcher = {
  deleteSubscription: billingDeleteSubscription,
  fetchLicenseFile: billingFetchLicenseFile,
  fetchSubscriptions: billingFetchProfileSubscriptions,
};

@connect(null, dispatcher)
export default class SubscriptionDetail extends Component {
  static propTypes = {
    deleteInfo: shape({
      error: string,
      isDeleting: bool,
      subscription_id: string,
    }),
    deleteSubscription: func.isRequired,
    error: string,
    fetchLicenseFile: func.isRequired,
    fetchSubscriptions: func.isRequired,
    licenses: shape({
      isFetching: bool,
      // Licenses are keyed off of subscription_id
      results: object,
      error: string,
    }),
    selectedNamespace: string.isRequired,
    selectedProductSubscriptions: arrayOf(object).isRequired,
    showDDCInstructions: func.isRequired,
    showSubscriptionList: func.isRequired,
    selectedUserOrOrg: object.isRequired,
  }

  state = {
    cancelSubscriptionId: '',
    cancelSubscriptionName: '',
    showCancelConfirmation: false,
    downloadLicenseSubscriptionId: '',
    downloadLicenseError: false,
  }

  hideCancelConfirmation = () => {
    this.setState({
      cancelSubscriptionId: '',
      cancelSubscriptionName: '',
      showCancelConfirmation: false,
    });
  }

  showCancelConfirmation = ({ subscription_id, subscription_name }) => () => {
    this.setState({
      cancelSubscriptionId: subscription_id,
      cancelSubscriptionName: subscription_name,
      showCancelConfirmation: true,
    });
  }

  cancelSubscription = () => {
    const { cancelSubscriptionId: subscription_id } = this.state;
    const {
      deleteSubscription,
      fetchSubscriptions,
      selectedUserOrOrg,
    } = this.props;
    deleteSubscription({ subscription_id }).then(() => {
      fetchSubscriptions({ docker_id: selectedUserOrOrg.id });
    });
  }

  downloadLicense = ({
    subscription_id,
    subscription_name,
    product_rate_plan,
  }) => () => {
    const { fetchLicenseFile } = this.props;
    this.setState({
      downloadLicenseSubscriptionId: subscription_id,
      downloadLicenseError: false,
    }, () => {
      fetchLicenseFile({ subscription_id })
        .then((res) => {
          // value contains the API response
          const downloadContent = JSON.stringify(res.value);
          this.setState({ downloadLicenseSubscriptionId: '' });
          const blob = new Blob([downloadContent], {
            type: 'text/plain;charset=utf-8',
          });
          saveAs(blob, `${subscription_name}.lic`, true);
          analytics.track('download_license', {
            license_tier: product_rate_plan,
          });
        })
        .catch(() => {
          this.setState({ downloadLicenseError: true });
        });
    });
  }

  // Returns true if the given subscription can be cancelled, false otherwise.
  //
  // A subscription may not be cancelled if:
  // 1. The initiator does not own the subscription OR
  // 2. The subscription is already in a cancelled state OR
  // 3. This is an 'offline' subscription.
  // Note that the presence of a subscription origin implies that this
  // is an offline subscription.
  canCancelSubscription(subscription) {
    const { isOwner } = this.props.selectedUserOrOrg;
    const { state, origin } = subscription;

    if (!isOwner) {
      return false;
    }

    if (state === CANCELLED) {
      return false;
    }

    if (origin) {
      return false;
    }

    return true;
  }

  renderProductSubscriptionSummary = (subs) => {
    // Subs is an array containing one or more subscriptions for the
    // same product id. When rendering the *product* information (and
    // aggregated subscription info), we must look at at least one subscription
    // to understand the product information. There is guaranteed to be at least
    // one sub, so we can get the information out of the first one.
    if (!subs || !subs.length) {
      return null;
    }
    const {
      label: productName,
      links,
      publisher_id,
      publisher,
    } = subs[0];
    const numActiveSubscriptions =
      filter(subs, ({ state }) => state === ACTIVE).length;
    let partnerName;
    let dockerPublisherID = DOCKER_GUID;
    if (isStaging() || isDev()) {
      dockerPublisherID = DOCKER;
    }
    if (publisher_id !== dockerPublisherID) {
      partnerName = (
        <div className={css.partnerName}>
          <CheckIcon className={css.check} size={SMALL} />
          {publisher.name}
        </div>
      );
    }
    const addS = numActiveSubscriptions === 1 ? '' : 's';
    const licenses = (
      <span>{`${numActiveSubscriptions} active subscription${addS}`}</span>
    );
    let licenseAgreement;
    // Add license agreement if it exists in the links
    if (links) {
      const link = find(links, ({ label }) => label === 'License Agreement');
      if (link) {
        const { label: text, url } = link;
        licenseAgreement =
          <span> | <a href={url} target="_blank">{text}</a></span>;
      }
    }
    return (
      <div>
        <div className={css.product}>
          <div className={css.productName}>{productName}</div>
          {partnerName}
        </div>
        <div className={css.nameAndLicense}>
          {this.props.selectedNamespace} | {licenses} {licenseAgreement}
        </div>
      </div>
    );
  }

  renderCancelSubscriptionLine(subscription) {
    const { isDeleting, subscription_id: deletingId } = this.props.deleteInfo;
    const { subscription_name, subscription_id } = subscription;

    if (!this.canCancelSubscription(subscription)) {
      return null;
    }

    const isDeletingThisSub = isDeleting && deletingId === subscription_id;
    if (isDeletingThisSub) {
      return <div>Canceling Subscription...</div>;
    }
    const subInfo = { subscription_id, subscription_name };
    return (
      <div
        className={css.underlined}
        onClick={this.showCancelConfirmation(subInfo)}
      >
        Cancel Subscription
      </div>
    );
  }

  renderPullCommand(subscription) {
    // TODO Kristie 8/17/19 Fix this when we actually set up the new
    // subscriptions page (and per sub pull commands)
    // This is a known issue until then
    const {
      plans: productPlans,
      product_rate_plan: currentRatePlanID,
    } = subscription;
    // TODO Kristie 8/17/19 This needs to be synced with billing
    let currentPlanDetail = find(productPlans, ({ id }) => {
      return id === currentRatePlanID;
    });
    // TODO Kristie 8/17/19 Remove this when billing plan ids are in sync with
    // the product plan ids
    if (!currentPlanDetail) {
      currentPlanDetail = productPlans[0];
    }
    const { repositories, default_version } = currentPlanDetail;
    // Assumption: there is only one repo per rate plan - this may change in the
    // future
    if (!repositories || !repositories[0]) {
      return null;
    }
    const { namespace, reponame } = repositories[0];
    return (
      <CopyPullCommand
        codeClassName={css.pullText}
        namespace={namespace}
        reponame={reponame}
        tag={default_version && default_version.linux}
      />
    );
  }

  renderPricingComponents(subscription) {
    const {
      pricing_components,
      // The specific rate plan that this subscription is subscribed to
      product_rate_plan: sub_rate_plan,
      // All possible rate plans and their pricing components for this product
      rate_plans,
    } = subscription;
    // Find the rate plan that matches this subscription's rate plan
    const plan = find(rate_plans, (rp) => rp.name === sub_rate_plan);
    // Do not show this user cancelled subscriptions if they are not an owner
    if (!plan || !pricing_components || !pricing_components.length) {
      return null;
    }
    const comps = pricing_components.map((pc) => {
      // Find the pricing component (with full information) from the billing
      // plan which matches this subscription's pricing component
      const priceCompInfo = find(plan.pricing_components, comp => {
        // Comp is the full information from billing api, pc is this sub's
        // data with the amount purchased (value)
        return comp === pc.name;
      });
      // Fallback to this subscription's component "name"
      const label = priceCompInfo && priceCompInfo.label || pc.name;
      return `${pc.value} ${label}`;
    });
    return (
      <div>
        <div>{plan.label}</div>
        <div>{comps.join(', ')}</div>
      </div>
    );
  }

  renderDownloadLicense(subscription) {
    const { downloadLicenseError, downloadLicenseSubscriptionId } = this.state;
    const { subscription_id } = subscription;
    const isDownloadingThisLicense =
      subscription_id === downloadLicenseSubscriptionId;
    const downloadError = isDownloadingThisLicense && downloadLicenseError;
    let text = 'License Key';
    if (isDownloadingThisLicense && !downloadError) {
      text = 'Downloading License...';
    } else if (isDownloadingThisLicense && downloadError) {
      text = 'Error Downloading License';
    }
    const icon = <DownloadIcon size={SMALL} />;
    const classes = classnames({
      [css.downloadLicense]: true,
      [css.downloadLicenseError]: isDownloadingThisLicense && downloadError,
    });
    return (
      <div
        className={classes}
        onClick={this.downloadLicense(subscription)}
      >
        {icon} {text}
      </div>
    );
  }

  renderLicense(subscription) {
    const { isOwner } = this.props.selectedUserOrOrg;
    const { subscription_id, eusa, state, product_rate_plan } = subscription;
    const { results } = this.props.licenses;
    // Find this license in results
    const license = results[subscription_id];
    // Do not render a license for a cancelled subscription
    if (!license || state === CANCELLED) {
      return null;
    }
    const { expiration } = license;
    const isLicenseExpired = isExpiredLicense(expiration);
    const isLegalAccepted = eusa && eusa.accepted;
    let licenseLine;

    if (!isLegalAccepted) {
      const text = isOwner ? 'Accept terms to download'
        : 'Admin must accept terms to download';
      licenseLine = <div>{text}</div>;
    } else if (isLicenseExpired) {
      licenseLine = <div>Expired License</div>;
    } else {
      licenseLine = this.renderDownloadLicense(subscription);
    }
    let expirationDate;
    // TODO Kristie 6/10/16 Incorporate other non-recurring subscriptions
    if (product_rate_plan === DDC_TRIAL_PLAN) {
      expirationDate = <div>Expires {moment(expiration).format('l')}</div>;
    }
    return (
      <div>
        {expirationDate}
        {licenseLine}
      </div>
    );
  }

  // A single subscription for this product
  renderSubscriptionLine = (subscription) => {
    const { subscription_name, subscription_id, state } = subscription;
    const { isOwner } = this.props.selectedUserOrOrg;
    if (!isOwner && state === CANCELLED) {
      return null;
    }
    // TODO Kristie 6/1/16 Include edit sub name
    // TODO Kristie 6/1/16 Include purchased by or registered
    return (
      <div key={subscription_id}>
        <hr />
        <div className={css.subscriptionLine} >
          <div>
            <div className={css.subscriptionName}>{subscription_name}</div>
            {this.renderPricingComponents(subscription)}
            {this.renderCancelSubscriptionLine(subscription)}
          </div>
          <div className={css.licenseInformation}>
            <div className={css.licenseStatus}>{startCase(state)}</div>
            {this.renderLicense(subscription)}
          </div>
        </div>
      </div>
      );
  }

  renderNewSubscriptionNotification() {
    const { selectedProductSubscriptions: subs } = this.props;
    if (!subs || !subs.length) {
      return null;
    }
    const { name, label } = subs[0];
    // TODO Get hint from API for paid subscription or not
    const isProductPaid = name === 'devddc';
    const newSubscriptions = filter(subs, (s) => isNewSubscription(s));
    if (!newSubscriptions.length) {
      return null;
    }
    let header;
    const details = 'Your subscription details are below.';
    if (isProductPaid) {
      header = `Payment Successful! Thank you for purchasing ${label}`;
    } else {
      header = `Thank you for subscribing to ${label}`;
    }
    return (
      <Card className={css.newSubscription} shadow>
        <div className={css.notificationHeader}>{header}</div>
        {details}
      </Card>
    );
  }

  renderErrorNotification() {
    // TODO Kristie 6/2/16 Error dismissal
    const { error } = this.props.deleteInfo;
    if (!error) {
      return null;
    }
    return (
      <Card className={css.error} shadow>
        <div className={css.notificationHeader}>Error</div>
        {['Sorry, we could not cancel your subscription.',
          'Please try again.'].join(' ')}
      </Card>
    );
  }

  renderCancelSubscriptionModal = () => {
    const {
      cancelSubscriptionId,
      cancelSubscriptionName,
      showCancelConfirmation,
    } = this.state;
    const { isDeleting, subscription_id: deletingId } = this.props.deleteInfo;
    const isDeletingThisSub = cancelSubscriptionId === deletingId && isDeleting;
    return (
      <Modal
        isOpen={showCancelConfirmation}
        className={css.modal}
        onRequestClose={this.hideCancelConfirmation}
      >
        <div
          className={css.closeModal}
          onClick={this.hideCancelConfirmation}
        >
          <CloseIcon size={SMALL} />
        </div>
        <div className={css.cancelHeader}>
          {`Cancel ${cancelSubscriptionName}?`}
        </div>
        <div className={css.cancelDetail}>
          This action cannot be undone.
        </div>
        <Button onClick={this.cancelSubscription}>
          {isDeletingThisSub ? 'Cancelling...' : 'Cancel Subscription'}
        </Button>
      </Modal>
    );
  }

  renderDDCInstructions() {
    const setupText = 'To install DDC on your machine, VM, or cloud instance ';
    const instructionsLink = (
      <div className={css.link} onClick={this.props.showDDCInstructions}>
        follow these instructions
      </div>
    );
    const supportText =
      'If you need assistance with your subscription please contact';
    const mailTo = 'mailto:support@docker.com?subject=Docker Datacenter';
    const supportLink = (
      <a href={mailTo} className={css.link}>{'support@docker.com'}</a>
    );
    return (
      <div>
        <div className={css.instructionsSection}>
          <div className={css.instructionsHeader}>Setup Guide</div>
          <div className={css.instructionsDetail}>
            {setupText} {instructionsLink}
          </div>
        </div>
        <div className={css.instructionsSection}>
          <div className={css.instructionsHeader}>Customer Support</div>
          <div className={css.instructionsDetail}>
            {supportText} {supportLink}
          </div>
        </div>
      </div>
    );
  }

  renderInstructions() {
    const { selectedProductSubscriptions: subs } = this.props;
    // At least one subscription otherwise render would have returned null
    const { product_id, instructions } = subs[0];
    if (product_id === DDC_ID) {
      return this.renderDDCInstructions();
    }
    if (!instructions) {
      return null;
    }
    return (
      <div>
        <Markdown rawMarkdown={instructions} />
      </div>
    );
  }

  renderSubscriptionHeader = (subs) => {
    const sampleSub = subs[0];
    const { logo_url, product_id } = sampleSub;
    const src = getLogo(logo_url);
    let pullCommand;
    if (isPullableProduct(product_id)) {
      pullCommand = this.renderPullCommand(sampleSub);
    }
    return (
      <div className={css.summaryAndPullCommand}>
        <div className={css.subscriptionSummary}>
          <ImageWithFallback
            src={src}
            className={css.icon}
            fallbackImage={FALLBACK_IMAGE_SRC}
            fallbackElement={FALLBACK_ELEMENT}
          />
          {this.renderProductSubscriptionSummary(subs)}
        </div>
        {pullCommand}
      </div>
    );
  }

  render() {
    const {
      selectedProductSubscriptions: subs,
      showSubscriptionList,
    } = this.props;
    if (!subs || !subs.length) {
      return null;
    }
    const { product_id } = subs[0];
    const maybeNewSubNotification = this.renderNewSubscriptionNotification();
    const maybeErrorNotication = this.renderErrorNotification();
    const maybeCancelSubscriptionModal = this.renderCancelSubscriptionModal();
    const maybeInstructions = this.renderInstructions();
    return (
      <div>
        <BackButtonArea onClick={showSubscriptionList} text="Subscriptions" />
        <Card key={product_id} className={css.subscriptionDetailCard} shadow>
          {maybeNewSubNotification}
          {maybeErrorNotication}
          {this.renderSubscriptionHeader(subs)}
          {map(subs, this.renderSubscriptionLine)}
          {maybeCancelSubscriptionModal}
        </Card>
        {maybeInstructions}
      </div>
    );
  }
}
