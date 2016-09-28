import React, { Component, PropTypes } from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';
import css from './styles.css';
import { Button } from 'common';
import { accountToggleMagicCarpet } from 'actions/account';
import routes from 'lib/constants/routes';
import { SUBSCRIPTIONS } from 'lib/constants/overlays';
import { DDC } from 'lib/constants/landingPage';
import { DDC_ID, DDC_TRIAL_PLAN } from 'lib/constants/eusa';
import { ACTIVE } from 'lib/constants/states/subscriptions';
import ProductDetailPage from 'marketplace/ProductDetailPage';
import filter from 'lodash/filter';
const { arrayOf, array, bool, func, number, object, shape, string } = PropTypes;

const mapStateToProps = (state, { params }) => {
  const { account, billing, marketplace, root } = state;
  const { bundles } = marketplace;
  const { products, subscriptions } = billing;
  const { currentUser, isCurrentUserWhitelisted } = account;
  const { id } = params;
  let activeProductSubscriptions = [];
  // Check if the user has any subscriptions to this product
  if (subscriptions && subscriptions.results) {
    activeProductSubscriptions = filter(subscriptions.results, (sub) => {
      return sub.product_id === id && sub.state === ACTIVE;
    });
  }
  return {
    isLoggedIn: currentUser && !!currentUser.id || false,
    isCurrentUserWhitelisted,
    isPageTransitioning: root && root.isPageTransitioning,
    subscriptions: {
      isFetching: subscriptions && subscriptions.isFetching,
      activeProductSubscriptions,
    },
    product: products && products[id] || {},
    bundle: bundles && bundles[id] || {},
  };
};

@connect(mapStateToProps, { accountToggleMagicCarpet })
export default class DDCDetail extends Component {
  static propTypes = {
    bundle: shape({
      categories: arrayOf(shape({
        name: string,
        label: string,
      })),
      default_version: shape({
        linux: string,
        windows: string,
      }),
      name: string,
      error: string,
      isFetching: bool,
      last_updated: string,
      logo_url: object,
      namespace: string,
      platforms: arrayOf(shape({
        name: string,
        label: string,
      })),
      popularity: number,
      publisher: shape({
        id: string,
        name: string,
      }),
      reponame: string,
      short_description: string,
      source: string,
      tags: shape({
        count: number,
        isFetching: bool,
        results: array,
      }),
    }),
    isLoggedIn: bool,
    isCurrentUserWhitelisted: bool,
    isPageTransitioning: bool,
    location: object.isRequired,
    params: object.isRequired,
    product: shape({
      id: string,
      isFetching: bool,
      rate_plans: array,
    }),
    subscriptions: shape({
      isFetching: bool,
      activeProductSubscriptions: array,
    }),
    accountToggleMagicCarpet: func.isRequired,
  }

  static contextTypes = {
    router: shape({
      push: func.isRequired,
    }).isRequired,
  }

  shouldComponentUpdate = (nextProps) => {
    const { isPageTransitioning } = nextProps;
    if (isPageTransitioning) {
      return false;
    }
    return true;
  }


  openMagicCarpet = () => {
    this.props.accountToggleMagicCarpet({ magicCarpet: SUBSCRIPTIONS });
  };

  linkTo = (location) => () => {
    this.context.router.push(location);
  }

  renderWhatsIncluded() {
    return (
      <div className={css.whatsIncluded}>
        <div className={css.subSectionTitle}>What's Included</div>
        {
          DDC.whatsIncluded.map((line) => {
            return <div key={line} className={css.included}>{line}</div>;
          })
        }
      </div>
    );
  }

  renderSubscriptionDetail() {
    const { subscriptions, location } = this.props;
    const { activeProductSubscriptions } = subscriptions;
    const numSubs = activeProductSubscriptions.length;
    if (numSubs) {
      const subs = numSubs === 1 ? 'subscription' : 'subscriptions';
      // add ?overlay=subscriptions to trigger opening the subscriptions MC
      // TODO Kristie 6/8/16 Make this work!
      const currentPage = `${location.pathname}?overlay=subscriptions`;
      return (
        <div className={css.subscriptionDetail}>
          <div className={css.subSectionTitle}>
            You have {numSubs} active {subs}
          </div>
          <div className={css.subSectionHelpText}>
            <Link
              to={currentPage}
              onClick={this.openMagicCarpet}
            >View Subscriptions</Link>
          </div>
        </div>
      );
    }
    return (
      <div className={css.subscriptionDetail}>
        <div className={css.subSectionTitle}>
          Starting at $150 per engine per month
        </div>
        <div className={css.subSectionHelpText}>
          Business Day and Business Critical Support available
        </div>
      </div>
    );
  }

  renderPurchaseOptions() {
    const ddcPurchase = routes.bundleDetailPurchase({ id: DDC_ID });
    return (
      <div>
        <Button
          className={css.button}
          onClick={
            this.linkTo({
              pathname: ddcPurchase,
              query: { plan: DDC_TRIAL_PLAN },
            })
          }
        >
          FREE 30-day evaluation
        </Button>
        <Button
          className={css.button}
          onClick={this.linkTo({ pathname: ddcPurchase })}
        >
          Buy Subscription
        </Button>
      </div>
    );
  }

  render() {
    const {
      bundle,
      isCurrentUserWhitelisted,
      isLoggedIn,
      product,
      subscriptions,
    } = this.props;
    const { isFetching: isFetchingSubs } = subscriptions;
    const isProductDoneFetching = !product.isFetching &&
      (product.id || product.error);
    const isSubscriptionDoneFetching = !isFetchingSubs;
    // Don't render until both tags and product are done so that there
    // is no rendering flash
    const isFetching = !isProductDoneFetching && !isSubscriptionDoneFetching;
    let displayElement = <div></div>;
    // If we cannot fetch billing product information ==> error
    const error = !bundle.isFetching && !isFetching && product.error;
    if (!isFetching) {
      displayElement = (
        <div>
          {this.renderWhatsIncluded()}
          {this.renderSubscriptionDetail()}
          {this.renderPurchaseOptions()}
        </div>
      );
    }
    const versionsAndSubscriptions = { displayElement, error, isFetching };
    return (
      <ProductDetailPage
        isLoggedInAndWhitelisted={isLoggedIn && isCurrentUserWhitelisted}
        imageOrBundle={bundle}
        versionsAndSubscriptions={versionsAndSubscriptions}
      />
    );
  }
}
