import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import {
  Button,
  Card,
  CheckIcon,
  ImageWithFallback,
  FullscreenLoading,
  FetchingError,
  Select,
} from 'common';
import { FALLBACK_IMAGE_SRC, FALLBACK_ELEMENT } from 'lib/constants/fallbacks';
import css from './styles.css';
import {
  filter,
  find,
  get,
  map,
  sortBy,
} from 'lodash';
import getLogo from 'lib/utils/get-largest-logo';
import { SMALL } from 'lib/constants/sizes';
import { PRIMARY, PANIC } from 'lib/constants/variants';
import isStaging from 'lib/utils/isStaging';
import isDev from 'lib/utils/isDevelopment';
import { DOCKER, DOCKER_GUID } from 'lib/constants/defaults';
import { ACTIVE, CANCELLED } from 'lib/constants/states/subscriptions';
import { EUSA_LINK } from 'lib/constants/eusa';
import {
  billingFetchProfileSubscriptions,
  billingUpdateSubscription,

  BILLING_BASE_URL,
} from 'actions/billing';
const noOp = () => {};
const { array, bool, func, shape, string, object } = PropTypes;

const dispatcher = {
  fetchSubscriptions: billingFetchProfileSubscriptions,
  updateSubscription: billingUpdateSubscription,
};
@connect(null, dispatcher)
export default class SubscriptionList extends Component {
  static propTypes = {
    currentUserNamespace: string,
    dockerSubscriptions: array.isRequired,
    isFetching: bool,
    fetchSubscriptions: func.isRequired,
    namespaceObjects: shape({
      isFetching: bool,
      results: object,
      error: string,
    }).isRequired,
    onSelectNamespace: func.isRequired,
    partnerSubscriptions: array.isRequired,
    selectedNamespace: string.isRequired,
    showSubscriptionDetail: func.isRequired,
    updateInfo: object,
    updateSubscription: func.isRequired,
    selectedUserOrOrg: object,
    error: string,
  }

  onSubscriptionClick = (productId) => () => {
    this.props.showSubscriptionDetail(productId);
  }

  acceptEusa = ({ subscription_id, product_id }) => () => {
    const payload = { eusa: { accepted: true } };
    const { id: docker_id } = this.props.selectedUserOrOrg;
    this.props.updateSubscription({ subscription_id, body: payload })
      .then(() => {
        // Fetch all the subscriptions so you have the most updated
        this.props.fetchSubscriptions({ docker_id }).then(() => {
          // Transition to the detail page for this product to download license
          this.props.showSubscriptionDetail(product_id);
        });
      });
  }

  // Returns the EUSA url associated with the given subscription.
  determineEUSA(subscription) {
    const { origin } = subscription;

    // Use a partner specific EUSA where applicable.
    // TODO - This will be obtained from the product catalog in the near future.
    if (origin === 'hpe') {
      return `${BILLING_BASE_URL}/eusa/hpe/eusa.pdf`;
    }

    // Use the standard Docker software EUSA
    return EUSA_LINK;
  }

  renderEusaAcceptance(subs) {
    const {
      namespaceObjects,
      selectedNamespace,
      updateInfo,
    } = this.props;
    // All subscriptions that require legal acceptance need to have the legal
    // terms accepted
    const { eusa_acceptance_required, product_id } = subs[0];
    const isOwner =
      get(namespaceObjects, ['results', selectedNamespace, 'isOwner']);
    if (!eusa_acceptance_required || !isOwner) {
      return null;
    }
    // See if you have to acccept terms for any subscriptions
    const subscriptionToAccept = find(subs, ({ eusa, state }) => {
      return eusa && !eusa.accepted && state !== CANCELLED;
    });
    // All subscriptions have accepted terms --> no action needed
    if (!subscriptionToAccept) {
      return null;
    }
    const { subscription_id } = subscriptionToAccept;
    // Prevent clicking on the buttons here from clicking on the parent div
    const stopProp = (e) => {
      e.stopPropagation();
    };
    // Grab that subscription's isUpdating and error information
    const { error, isUpdating } = updateInfo[subscription_id] || {};
    let buttonText = 'Accept Terms';
    if (isUpdating && !error) {
      buttonText = 'Accepting Terms...';
    } else if (error) {
      buttonText = 'Error! Please Try Again';
    }

    const eusaLink = this.determineEUSA(subscriptionToAccept);

    // We only have to cover the EUSA link because all trials will be created
    // through the online form and will have accepted eval terms
    return (
      <div className={css.noClick} onClick={stopProp}>
        To access your License keys, please accept the
        <a href={eusaLink} target="_blank" className={css.terms}>
          Terms of Service
        </a>
        <Button
          onClick={this.acceptEusa({ product_id, subscription_id })}
          variant={error ? PANIC : PRIMARY}
        >
          {buttonText}
        </Button>
      </div>
    );
  }

  renderProductSubscriptionSummary = (subs) => {
    if (!subs || !subs.length) {
      return null;
    }
    const {
      publisher_id,
      publisher = {},
      label,
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
    return (
      <div>
        <div className={css.product}>
          <div className={css.productName}>{label}</div>
          {partnerName}
        </div>
        <div className={css.nameAndLicense}>
          {licenses}
        </div>
      </div>
    );
  }

  renderProductSubscriptions = (subs) => {
    // Subscription is an array containing one or more subscriptions for the
    // same product id
    const { product_id, logo_url } = subs[0];
    const src = getLogo(logo_url);
    const maybeEusaAcceptance = this.renderEusaAcceptance(subs);
    return (
      <Card
        key={product_id}
        className={css.subscriptionCard}
        hover
        shadow
        onClick={this.onSubscriptionClick(product_id)}
      >
        <div className={css.subscriptionSummary}>
          <ImageWithFallback
            className={css.icon}
            fallbackElement={FALLBACK_ELEMENT}
            fallbackImage={FALLBACK_IMAGE_SRC}
            src={src}
          />
          {this.renderProductSubscriptionSummary(subs)}
        </div>
        <div>{maybeEusaAcceptance}</div>
      </Card>
    );
  }

  renderNoSubscriptions(type) {
    return (
      <Card shadow>You have no {type} subscriptions.</Card>
    );
  }

  renderPartnerSubscriptionList() {
    const { partnerSubscriptions } = this.props;
    let subs = map(partnerSubscriptions, this.renderProductSubscriptions);
    if (!partnerSubscriptions.length) {
      return null;
    }
    return (
      <div>
        <div className={css.sectionTitle}>Docker Partner Services</div>
        <div>{subs}</div>
      </div>
    );
  }

  renderDockerSubscriptionList() {
    const { dockerSubscriptions } = this.props;
    let subs = map(dockerSubscriptions, this.renderProductSubscriptions);
    if (!dockerSubscriptions.length) {
      subs = this.renderNoSubscriptions('Docker');
    }
    return (
      <div>
        <div className={css.dockerSectionTitle}>
          <div className={css.sectionTitle}>Docker Services</div>
          <div>{this.renderNamespaceSelect()}</div>
        </div>
        <div>{subs}</div>
      </div>
    );
  }

  renderNamespaceSelect() {
    const { namespaceObjects, currentUserNamespace } = this.props;
    let options = map(namespaceObjects.results, (userObj) => {
      const namespace = userObj.username || userObj.orgname;
      return { value: namespace, label: namespace };
    });
    if (!namespaceObjects.results[currentUserNamespace]) {
      // current user should be included in namespace objects
      // BUT if it's not, we should add it to the front
      options.unshift(
        { value: currentUserNamespace, label: currentUserNamespace }
      );
    }
    options = sortBy(options, (option) => {
      return option.value !== currentUserNamespace;
    });
    return (
      <div className={css.namespaceSelectWrapper}>
        Account
        <Select
          className={css.namespaceSelect}
          clearable={false}
          ignoreCase
          onBlur={noOp}
          onChange={this.props.onSelectNamespace}
          options={options}
          placeholder="Account"
          value={this.props.selectedNamespace}
        />
      </div>
    );
  }


  render() {
    const {
      isFetching,
      error,
    } = this.props;
    if (isFetching) {
      return <FullscreenLoading />;
    } else if (error) {
      return (
        <div>
          <div className={css.empty}>
            {this.renderNamespaceSelect()}
          </div>
          <div className={css.fetchingError}>
            <FetchingError resource="your subscriptions" />
          </div>
        </div>
      );
    }
    return (
      <div>
        {this.renderDockerSubscriptionList()}
        {this.renderPartnerSubscriptionList()}
      </div>
    );
  }
}
