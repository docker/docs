import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import {
  FullscreenLoading,
} from 'common';
import forEach from 'lodash/forEach';
import get from 'lodash/get';
import isStaging from 'lib/utils/isStaging';
import isDev from 'lib/utils/isDevelopment';
import { DOCKER, DOCKER_GUID } from 'lib/constants/defaults';
import SubscriptionList from './SubscriptionList';
import SubscriptionDetail from './SubscriptionDetail';
import DDCInstructions from './DDCInstructions';
import { setIsOwnerParam } from './helpers';
import css from './styles.css';
import { isProductBundle } from 'lib/utils/product-utils';
import {
  CANCELLED,
  DDC_INSTRUCTIONS,
  SUBSCRIPTION_DETAIL,
  SUBSCRIPTION_LIST,
} from 'lib/constants/states/subscriptions';
import {
  accountFetchCurrentUser,
  accountFetchUser,
  accountFetchUserOrgs,
  accountSelectNamespace,
} from 'actions/account';
import {
  billingFetchLicenseDetail,
  billingFetchProduct,
  billingFetchProfileSubscriptionsAndProducts,
} from 'actions/billing';
import { repositoryFetchOwnedNamespaces } from 'actions/repository';
import {
  marketplaceFetchBundleDetail,
  marketplaceFetchRepositoryDetail,
} from 'actions/marketplace';
const {
  array,
  arrayOf,
  bool,
  func,
  object,
  objectOf,
  string,
  shape,
} = PropTypes;

const mapStateToProps = ({ account, billing, marketplace }) => {
  const {
    currentUser,
    // object of all fetched namespace objects (owned & orgs)
    namespaceObjects,
    ownedNamespaces,
    selectedNamespace,
  } = account;
  const currentUserNamespace = currentUser.username || '';
  // Set the isOwner property for each user / org object in results
  namespaceObjects.results =
    setIsOwnerParam(namespaceObjects.results, ownedNamespaces);
  // Is this user an owner for the _currently selected_ namespace?
  const isOwner =
    get(namespaceObjects, ['results', selectedNamespace, 'isOwner']);
  const { subscriptions, products: billingProducts } = billing;
  const { bundles, images } = marketplace;
  const {
    delete: deleteInfo,
    error,
    isFetching,
    results,
    update: updateInfo,
  } = subscriptions;
  const dockerSubscriptions = [];
  const partnerSubscriptions = [];
  const subscriptionsByProductId = {};
  if (!isFetching) {
    let dockerPublisherID = DOCKER_GUID;
    if (isStaging() || isDev()) {
      dockerPublisherID = DOCKER;
    }
    // Group each subscription by product
    forEach(results, (subscription) => {
      // Look up the product id associated with this subscription id
      const { product_id, state } = subscription;
      if (!isOwner && state === CANCELLED) {
        // Non owners cannot view cancelled subscriptions
        return;
      }
      // Get the billing product info for that product (rate_plans, etc.)
      const billingProductInfo = billingProducts[product_id];
      // Get the product catalog information for that product (ex. logo, name)
      let productCatalogInfo;
      if (isProductBundle(product_id)) {
        productCatalogInfo = bundles[product_id] || {};
      } else {
        productCatalogInfo = images.certified[product_id] || {};
      }

      // If we CANNOT get the billing product or product catalog information,
      // do NOT display this subscription
      const productErr = productCatalogInfo && productCatalogInfo.error;
      const billingErr = billingProductInfo && billingProductInfo.error;
      if (productErr || billingErr) {
        return;
      }
      // Combine the subscription info and product info collision on `name`)
      const allInfo = {
        // Subscription info and billing product info both have name
        subscription_name: subscription.name,
        // Product catalog & subscription both have eusa (we want subscription)
        ...productCatalogInfo,
        ...billingProductInfo,
        ...subscription,
        // TODO Kristie 8/17/16 Separate this info better so that there are no
        // collisions
        name: productCatalogInfo.name,
      };
      // Check if this product exists in the subscriptions already
      const existingSub = subscriptionsByProductId[product_id];
      // Add this subscription + product info to the list of subscriptions for
      // this product, or create a list if none exists
      const newSubs = existingSub ? [...existingSub, allInfo] : [allInfo];
      subscriptionsByProductId[product_id] = newSubs;
    });
    // Separate each product (and its array of subscriptions) by vendor
    forEach(subscriptionsByProductId, (productSubs) => {
      // Guaranteed to be at least one subscription for this product
      const { publisher_id } = productSubs[0];
      if (publisher_id === dockerPublisherID) {
        dockerSubscriptions.push(productSubs);
      } else {
        partnerSubscriptions.push(productSubs);
      }
    });
  }
  return {
    currentUserNamespace,
    deleteInfo,
    dockerSubscriptions,
    error,
    isFetching,
    licenses: subscriptions && subscriptions.licenses,
    namespaceObjects,
    partnerSubscriptions,
    selectedNamespace,
    subscriptionsByProductId,
    updateInfo,
  };
};

const dispatcher = {
  accountFetchCurrentUser,
  accountFetchUser,
  accountFetchUserOrgs,
  accountSelectNamespace,
  billingFetchLicenseDetail,
  billingFetchProduct,
  billingFetchProfileSubscriptionsAndProducts,
  marketplaceFetchBundleDetail,
  marketplaceFetchRepositoryDetail,
  repositoryFetchOwnedNamespaces,
};

/*
 * Subscriptions is the connected component that handles showing either
 * a SubscriptionList or a SubscriptionDetail page. Subscriptions is only
 * shown within the context of the MagicCarpet by clicking on the Subscriptions
 * link in the Nav dropdown. There is no route within the app.
 */
@connect(mapStateToProps, dispatcher)
export default class Subscriptions extends Component {
  static propTypes = {
    accountFetchCurrentUser: func.isRequired,
    accountFetchUser: func.isRequired,
    accountFetchUserOrgs: func.isRequired,
    accountSelectNamespace: func.isRequired,
    billingFetchLicenseDetail: func.isRequired,
    billingFetchProduct: func.isRequired,
    billingFetchProfileSubscriptionsAndProducts: func.isRequired,
    marketplaceFetchBundleDetail: func.isRequired,
    marketplaceFetchRepositoryDetail: func.isRequired,
    repositoryFetchOwnedNamespaces: func.isRequired,

    currentUserNamespace: string,
    error: string,
    // Object of all fetched namespace objects keyed off of namespace
    namespaceObjects: shape({
      isFetching: bool,
      results: object,
      error: string,
    }).isRequired,
    isFetching: bool,
    selectedNamespace: string,
    // Values are arrays containing subscription(s) for a product (key)
    subscriptionsByProductId: objectOf(array),

    // The following props are not being used in this component
    // but being passed down through to it's children
    deleteInfo: object,
    // Array of arrays containing subscription(s) for a docker product
    dockerSubscriptions: arrayOf(array),
    licenses: shape({
      isFetching: bool,
      // Licenses are keyed off of subscription_id
      results: object,
      error: string,
    }),
    // Array of arrays containing subscription(s) for a product
    partnerSubscriptions: arrayOf(array),
    updateInfo: object,
  }

  state = {
    // View list of subscriptions (grouped by product) or a detail view of all
    // subscriptions for a selected product
    currentView: SUBSCRIPTION_LIST,
    isInitializing: true,
    selectedProductId: '',
  }

  componentWillMount() {
    const {
      accountFetchCurrentUser: fetchCurrentUser,
      accountFetchUserOrgs: fetchUserOrgs,
      accountSelectNamespace: selectNamespace,
      billingFetchProfileSubscriptionsAndProducts: fetchProfileSubscriptions,
      repositoryFetchOwnedNamespaces: fetchOwnedNamespaces,
    } = this.props;
    Promise.when([
      fetchOwnedNamespaces(),
      fetchUserOrgs(),
      fetchCurrentUser().then((res) => {
        const { username: namespace, id: docker_id } = res.value;
        Promise.all([
          fetchProfileSubscriptions({ docker_id }),
          selectNamespace({ namespace }),
        ]);
      }).catch(() => { this.setState({ isInitializing: false }); }),
    ]).then(() => {
      this.setState({
        isInitializing: false,
      });
    }).catch(() => {
      this.setState({ isInitializing: false });
    });
  }

  onSelectNamespace = ({ value: namespace }) => {
    const {
      accountFetchUser: fetchUser,
      accountSelectNamespace: selectNamespace,
      currentUserNamespace,
      billingFetchProfileSubscriptionsAndProducts: fetchProfileSubscriptions,
      namespaceObjects,
    } = this.props;
    // NOTE
    // currently all options are being populated from already fetched namespaces
    // BUT just in case - this will fetch the namespaceObject
    if (!namespaceObjects.results[namespace]) {
      const isOrg = namespace !== currentUserNamespace;
      return fetchUser({ namespace, isOrg }).then((userRes) => {
        const { id: docker_id } = userRes.value;
        Promise.all([
          selectNamespace({ namespace }),
          fetchProfileSubscriptions({ docker_id }),
        ]);
      });
    }
    const { id } = namespaceObjects.results[namespace];
    return Promise.all([
      selectNamespace({ namespace }),
      fetchProfileSubscriptions({ docker_id: id }),
    ]);
  }

  showSubscriptionDetail = (productId) => {
    this.setState({
      currentView: SUBSCRIPTION_DETAIL,
      selectedProductId: productId,
    });
  }

  showSubscriptionList = () => {
    this.setState({ currentView: SUBSCRIPTION_LIST });
  }

  showDDCInstructions = () => {
    this.setState({ currentView: DDC_INSTRUCTIONS });
  }

  generateSelectOptions(namespaceObjects) {
    const namespaceArray = [];
    forEach(namespaceObjects, (val, key) => {
      namespaceArray.push({ value: key, label: key });
    });
    return namespaceArray;
  }

  render() {
    const {
      currentView,
      isInitializing,
      selectedProductId,
    } = this.state;
    const {
      currentUserNamespace,
      deleteInfo,
      dockerSubscriptions,
      error,
      isFetching,
      licenses,
      namespaceObjects,
      partnerSubscriptions,
      selectedNamespace,
      subscriptionsByProductId,
      updateInfo,
    } = this.props;
    const selectedUserOrOrg = namespaceObjects.results[selectedNamespace];
    const selectedProductSubs = subscriptionsByProductId[selectedProductId];
    let content;
    if (isInitializing || isFetching) {
      content = <FullscreenLoading />;
    } else if (currentView === SUBSCRIPTION_LIST) {
      content = (
        <SubscriptionList
          currentUserNamespace={currentUserNamespace}
          selectedUserOrOrg={selectedUserOrOrg}
          dockerSubscriptions={dockerSubscriptions}
          namespaceObjects={namespaceObjects}
          onSelectNamespace={this.onSelectNamespace}
          partnerSubscriptions={partnerSubscriptions}
          selectedNamespace={selectedNamespace}
          showSubscriptionDetail={this.showSubscriptionDetail}
          updateInfo={updateInfo}
          error={error}
        />
      );
    } else if (currentView === DDC_INSTRUCTIONS) {
      content = (
        <DDCInstructions
          showSubscriptionDetail={this.showSubscriptionDetail}
        />
      );
    } else {
      content = (
        <SubscriptionDetail
          selectedUserOrOrg={selectedUserOrOrg}
          deleteInfo={deleteInfo}
          error={error}
          licenses={licenses}
          selectedNamespace={selectedNamespace}
          selectedProductSubscriptions={selectedProductSubs}
          showDDCInstructions={this.showDDCInstructions}
          showSubscriptionList={this.showSubscriptionList}
        />
      );
    }
    return (
      <div className={css.content}>
        {content}
      </div>
    );
  }
}
