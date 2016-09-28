import React, { Component, PropTypes } from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';
import css from './styles.css';
import {
  Button,
  CopyPullCommand,
  // ShieldIcon,
} from 'components/common';
import { accountToggleMagicCarpet } from 'actions/account';
import sanitize from 'lib/utils/remove-undefined';
import {
  getPriceStringForRatePlan,
  isAnonymousDownloadProduct,
} from 'lib/utils/billing-plan-utils';
// import { SMALL } from 'lib/constants/sizes';
import { PANIC, PRIMARY } from 'lib/constants/variants';
import { SUBSCRIPTIONS } from 'lib/constants/overlays';
import { EUSA_LINK, DOWNLOAD_ATTRIBUTES } from 'lib/constants/eusa';
import { ACTIVE } from 'lib/constants/states/subscriptions';
import ProductDetailPage from 'marketplace/ProductDetailPage';
import {
  billingCreateProfile,
  billingCreateSubscription,
  billingFetchProfileSubscriptions,
} from 'actions/billing';
import find from 'lodash/find';
import filter from 'lodash/filter';
const { arrayOf, array, bool, func, number, object, shape, string } = PropTypes;

const mapStateToProps = (state, { params }) => {
  const { account, billing, marketplace, root } = state;
  const { certified } = marketplace && marketplace.images;
  const { products, profiles, subscriptions } = billing;
  const {
    currentUser,
    isCurrentUserWhitelisted,
    namespaceObjects,
    selectedNamespace,
    userEmails,
  } = account;
  const selectedUser = namespaceObjects.results[selectedNamespace] || {};
  const billingProfile = profiles.results[selectedUser.id];
  const billingProfileError = profiles.error;
  const { id } = params;
  let activeProductSubscriptions = [];
  // Check if the user has any subscriptions to this product
  if (subscriptions && subscriptions.results) {
    activeProductSubscriptions = filter(subscriptions.results, (sub) => {
      return sub.product_id === id && sub.state === ACTIVE;
    });
  }
  const image = certified[id] || {};
  const initiallySelectedPlan = image.plans && image.plans[0] || {};
  return {
    billingProfile,
    billingProfileError,
    selectedNamespace,
    fetchedNamespaces: namespaceObjects.results,
    userEmails: userEmails.results,
    isLoggedIn: currentUser && !!currentUser.id || false,
    isCurrentUserWhitelisted,
    isPageTransitioning: root && root.isPageTransitioning,
    subscriptions: {
      isFetching: subscriptions && subscriptions.isFetching,
      error: subscriptions && subscriptions.error,
      activeProductSubscriptions,
    },
    product: products[id] || {},
    image,
    initiallySelectedPlan,
  };
};

const dispatcher = {
  createProfile: billingCreateProfile,
  createSubscription: billingCreateSubscription,
  fetchSubscriptions: billingFetchProfileSubscriptions,
  toggleMagicCarpet: accountToggleMagicCarpet,
};

@connect(mapStateToProps, dispatcher)
export default class PartnerImageDetail extends Component {
  static propTypes = {
    billingProfile: object,
    billingProfileError: string,
    userEmails: array,
    fetchedNamespaces: object.isRequired,
    selectedNamespace: string.isRequired,
    image: shape({
      categories: arrayOf(shape({
        name: string,
        label: string,
      })),
      name: string.isRequired,
      error: string,
      isFetching: bool,
      last_updated: string,
      links: arrayOf(shape({
        label: string,
        url: string,
      })),
      logo_url: object.isRequired,
      platforms: arrayOf(shape({
        name: string,
        label: string,
      })),
      popularity: number,
      plans: arrayOf(shape({
        default_version: shape({
          linux: string,
          windows: string,
        }),
      })),
      publisher: shape({
        id: string,
        name: string,
      }),
      screenshots: arrayOf(shape({
        label: string,
        url: string,
      })),
      short_description: string,
      source: string,
    }),
    initiallySelectedPlan: object,
    location: shape({
      pathname: string,
    }).isRequired,
    params: object.isRequired,

    isLoggedIn: bool,
    isCurrentUserWhitelisted: bool,
    isPageTransitioning: bool,
    // Billing information
    product: shape({
      id: string,
      isFetching: bool,
      error: string,
      rate_plans: array,
    }),
    subscriptions: shape({
      isFetching: bool,
      activeProductSubscriptions: array,
    }),

    createProfile: func.isRequired,
    createSubscription: func.isRequired,
    fetchSubscriptions: func.isRequired,
    toggleMagicCarpet: func.isRequired,
  }

  constructor(props) {
    super(props);
    this.state = {
      // Use this to track state so that there is no flash in between
      // creating and fetching subscriptions
      isCreatingSubscription: false,
      // Automatically select the first plan
      currentlySelectedPlan: props.initiallySelectedPlan,
    };
  }

  // When you navigate from one PartnerImageDetail to another without changing
  // the route in between, React will not call the constructor because the
  // initial component will not be unmounted. Thus, the state will not be
  // properly initialized for the new PartnerImageDetail and it will hold
  // the old PartnerImageDetail's last selected plan. To fix this, we set the
  // initial state (as it would have happened in the constructor) if we are
  // transitioning to a new PartnerImageDetail page
  componentWillReceiveProps = (nextProps) => {
    const currentID = this.props.image.id;
    const nextID = nextProps.image.id;
    if (currentID !== nextID) {
      this.setState({
        currentlySelectedPlan: nextProps.initiallySelectedPlan,
      });
    }
  }

  // Prevent the page from flashing / updating when transitioning to another
  // PartnerImageDetail page (same issue with component not unmounting as
  // described above)
  shouldComponentUpdate = (nextProps) => {
    const { isPageTransitioning } = nextProps;
    if (isPageTransitioning) {
      return false;
    }
    return true;
  }

  onChangeProductTiers = (planObj) => () => {
    this.setState({ currentlySelectedPlan: planObj });
  }

  getPricingForPlan = (rate_plan_id) => {
    const { rate_plans } = this.props.product;
    // Find this plan
    const plan = find(rate_plans, ({ id }) => id === rate_plan_id);
    return getPriceStringForRatePlan(plan);
  }

  openMagicCarpet = () => {
    this.props.toggleMagicCarpet({ magicCarpet: SUBSCRIPTIONS });
  }

  subscribeToImage = () => {
    const {
      billingProfile,
      billingProfileError,
      createProfile,
      createSubscription,
      fetchedNamespaces,
      fetchSubscriptions,
      selectedNamespace,
      userEmails,
    } = this.props;
    // TODO Kristie 8/17/16 Update / remove this when we switch to integrate
    // the purchase page
    const userObject = fetchedNamespaces[selectedNamespace];
    const auth = { docker_id: userObject.id };
    const { id, label, name, rate_plans } = this.props.product;
    const plan = rate_plans && rate_plans.length && rate_plans[0]
      && rate_plans[0].name;
    const { publisher } = this.props.image;
    const date = new Date();
    const prodName = label || name;
    const subData = {
      name: `${prodName} Subscription ${date.toDateString()}`,
      pricing_components: [],
      product_id: id,
      product_rate_plan: plan,
      eusa: {
        accepted: true,
      },
    };

    // =========================================================================
    // Track subscription data (omit any 'undefined' fields)
    const trackedSubData = {
      docker_id: auth.docker_id,
      product_id: id,
      product_name: prodName,
      publisher_id: publisher && publisher.id,
      publisher_name: publisher && publisher.name,
    };
    analytics.track('create_subscription', sanitize(trackedSubData));
    // =========================================================================


    this.setState({ isCreatingSubscription: true }, () => {
      if (!billingProfile || (billingProfileError)) {
        const primary = find(userEmails, (emailObj) => {
          return emailObj.primary;
        });
        const profile = {
          first_name: '',
          last_name: '',
          email: primary && primary.email,
        };
        const data = { ...auth, profile };
        createProfile(data).then(() => {
          createSubscription({ ...subData, ...auth })
            .then(() => {
              fetchSubscriptions(auth)
                .then(() => {
                  this.setState({ isCreatingSubscription: false });
                })
                .catch(() => {
                  this.setState({ isCreatingSubscription: false });
                });
            }).catch(() => {
              this.setState({ isCreatingSubscription: false });
            });
        });
      } else {
        createSubscription({
          ...subData,
          ...auth,
        }).then(() => {
          fetchSubscriptions(auth)
            .then(() => { this.setState({ isCreatingSubscription: false }); })
            .catch(() => { this.setState({ isCreatingSubscription: false }); });
        }).catch(() => { this.setState({ isCreatingSubscription: false }); });
      }
    });
  }

  /* Render the number of subscriptions for this user to for this product, if
   * any.
   */
  renderSubscriptions() {
    const { subscriptions, location } = this.props;
    const { activeProductSubscriptions } = subscriptions;
    const numSubs = activeProductSubscriptions.length;
    if (!numSubs) { return null; }
    const subs = numSubs === 1 ? 'subscription' : 'subscriptions';
    // add ?overlay=subscriptions to trigger opening the subscriptions MC
    const currentPage = `${location.pathname}?overlay=subscriptions`;
    const link = (
      <Link to={currentPage} onClick={this.openMagicCarpet}>
        {numSubs} active {subs}
      </Link>
    );
    return (
      <div className={css.subscriptionDetail}>
        You have {link}
      </div>
    );
  }

  /* Render the product tiers available for this product, which is a minimum of
   * one (may be free or paid) and (currently) a maximum of 2.
   */
  renderProductTiers() {
    const { plans = [] } = this.props.image;
    const { id, description } = this.state.currentlySelectedPlan;
    let productTierTabs;
    if (plans && plans.length) {
      const productTiers = plans.map(plan => {
        const isSelected = plan.id === id;
        const tabClasses = isSelected ? css.selectedTab : css.tab;
        return (
          <div
            key={plan.id}
            onClick={this.onChangeProductTiers(plan)}
            className={tabClasses}
          >
            <div className={css.name}>{plan.name}</div>
            <div className={css.price}>{this.getPricingForPlan(plan.id)}</div>
          </div>
        );
      });
      productTierTabs = (
        <div className={css.tabs}>
          {productTiers}
        </div>
      );
    }
    const rpDescription = <div className={css.description}>{description}</div>;
    return (
      <div>
        {productTierTabs}
        {rpDescription}
      </div>
    );
  }

  /* The purchase options or pull command area can be in one of the following
   * states, depending on the currentlySelectedPlan
   * 1. currentlySelectedPlan is an anonymous download
   *   > Show pull command for that plan
   * 2. currentlySelectedPlan requires login and the user is logged in
   *   > Show pull command for that plan
   * 3. currentlySelectedPlan requires login and the user is NOT logged in
   *   > Show "please login" text
   * 4. currentlySelectedPlan requires a subscription
   *   > Show purchase options for that plan (button to subscribe)
   *   > Pull command can be viewed in the subscriptions pane
   */
  renderPurchaseOptionsOrPullCommand() {
    const { isLoggedIn } = this.props;
    const { download_attribute } = this.state.currentlySelectedPlan;
    const requiresLogin = download_attribute === DOWNLOAD_ATTRIBUTES.LOGGED_IN;
    const isAnonymous = download_attribute === DOWNLOAD_ATTRIBUTES.ANONYMOUS;
    // Case 1 or 2
    if (isAnonymous || (requiresLogin && isLoggedIn)) {
      return this.renderPullCommand();
    }
    return this.renderPurchaseOptions();
  }

  renderPurchaseOptions() {
    const { isLoggedIn, billingProfileError } = this.props;
    const { eusa } = this.state.currentlySelectedPlan;
    const { isCreatingSubscription } = this.state;
    const { error } = this.props.subscriptions;
    let text = 'Subscribe';
    const termsOfService = eusa || EUSA_LINK;
    // Error setting up subscription or setting up account
    const subscriptionError = error || billingProfileError;
    let isDisabled = false;
    let legalText = (
      <div className={css.legalText}>
        By clicking Subscribe you accept the
        <a href={termsOfService} target="_blank" className={css.terms}>
          Publisher Terms
        </a>
      </div>
    );
    if (subscriptionError) {
      text = 'Unable to subscribe';
      isDisabled = false;
    } else if (isCreatingSubscription) {
      text = 'Subscribing...';
      isDisabled = true;
    } else if (!isLoggedIn) {
      // Must always be logged in to subscribe
      text = 'Please Login to Subscribe';
      legalText = '';
      isDisabled = true;
    }
    return (
      <div>
        {legalText}
        <Button
          className={css.button}
          disabled={isDisabled}
          fullWidth
          onClick={this.subscribeToImage}
          variant={subscriptionError ? PANIC : PRIMARY}
        >
          {text}
        </Button>
      </div>
    );
  }

  renderPullCommand() {
    const { repositories, default_version } = this.state.currentlySelectedPlan;
    // Assumption: there is only one repo per rate plan - this may change in the
    // future
    if (!repositories || !repositories[0]) {
      return null;
    }
    const { namespace, reponame } = repositories[0];
    return (
      <CopyPullCommand
        namespace={namespace}
        reponame={reponame}
        tag={default_version && default_version.linux}
      />
    );
  }

  render() {
    const { image, product, isLoggedIn, isCurrentUserWhitelisted } = this.props;
    const { tags = {} } = image;
    const isTagsDoneFetching = !tags.isFetching && (tags.results || tags.error);
    const isProductDoneFetching = !product.isFetching &&
      (product.id || product.error);
    // Don't render until both tags and product are done so that there
    // is no rendering flash
    const isFetching = !isTagsDoneFetching && !isProductDoneFetching;
    let displayElement = <div></div>;
    // If we cannot fetch tag information or billing
    // product information, and it is not an anonymous download ==> error
    const isAnonymous = isAnonymousDownloadProduct(image);
    const fetchError = tags.error && product.error;
    const error = !image.isFetching && !isFetching &&
      fetchError && !isAnonymous;
    if (!isFetching) {
      displayElement = (
        <div className={css.flexColumn}>
          {this.renderProductTiers()}
          {this.renderSubscriptions()}
          {this.renderPurchaseOptionsOrPullCommand()}
        </div>
      );
    }
    const versionsAndSubscriptions = { displayElement, error, isFetching };
    return (
      <ProductDetailPage
        isLoggedInAndWhitelisted={isLoggedIn && isCurrentUserWhitelisted}
        imageOrBundle={image}
        versionsAndSubscriptions={versionsAndSubscriptions}
      />
    );
  }
}
