import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';

import InvoicesView from './Invoices';
import PaymentMethodsView from './PaymentMethods';

import {
  accountFetchCurrentUser,
  accountFetchUser,
  accountFetchUserEmails,
  accountSelectNamespace,
  billingFetchInvoices,
  billingFetchPaymentMethods,
  billingFetchProfile,
  repositoryFetchOwnedNamespaces,
} from './actions';
import {
  AccountSelect,
  Card,
  FullscreenLoading,
  Tab,
  Tabs,
} from './lib/common';
import {
  PAYMENTS,
  METHODS,
} from './lib/constants';

import css from './styles.css';
const { arrayOf, func, object, string, shape, bool } = PropTypes;

const mapStateToProps = ({ account, billing }) => {
  const {
    currentUser,
    namespaceObjects,
    ownedNamespaces,
    selectedNamespace,
  } = account;
  const {
    profiles,
  } = billing;
  return {
    billingProfiles: profiles,
    currentUser,
    namespaceObjects,
    ownedNamespaces,
    selectedNamespace,
  };
};
const dispatcher = {
  accountFetchUser,
  accountFetchUserEmails,
  accountFetchCurrentUser,
  accountSelectNamespace,
  billingFetchInvoices,
  billingFetchPaymentMethods,
  billingFetchProfile,
  repositoryFetchOwnedNamespaces,
};

@connect(mapStateToProps, dispatcher)
export default class BillingProfile extends Component {
  static propTypes={
    billingProfiles: shape({
      isFetching: bool.isRequired,
      results: object.isRequired,
    }),
    currentUser: object.isRequired,
    namespaceObjects: object.isRequired,
    // List of namespace strings owned by User
    ownedNamespaces: arrayOf(string),
    selectedNamespace: string.isRequired,
    // Action Props
    accountFetchUser: func.isRequired,
    accountFetchUserEmails: func.isRequired,
    accountFetchCurrentUser: func.isRequired,
    accountSelectNamespace: func.isRequired,
    billingFetchInvoices: func.isRequired,
    billingFetchPaymentMethods: func.isRequired,
    billingFetchProfile: func.isRequired,
    repositoryFetchOwnedNamespaces: func.isRequired,
  };

  state = {
    currentView: PAYMENTS,
    isInitializing: true,
    selectedNamespace: '',
  }

  componentWillMount() {
    const {
      // actions
      accountFetchCurrentUser: fetchCurrentUser,
      accountFetchUserEmails: fetchUserEmails,
      accountSelectNamespace: selectNamespace,
      billingFetchInvoices: fetchInvoices,
      billingFetchPaymentMethods: fetchPaymentMethods,
      billingFetchProfile: fetchBillingProfiles,
      repositoryFetchOwnedNamespaces: fetchOwnedNamespaces,
    } = this.props;
    Promise.when([
      fetchCurrentUser().then((userRes) => {
        const { username: namespace, id: docker_id } = userRes.value;
        return Promise.when([
          fetchInvoices({ docker_id }),
          fetchUserEmails({ user: namespace }),
          selectNamespace({ namespace }),
          fetchBillingProfiles({ docker_id, isOrg: false }).then((res) => {
            if (res.value.profile) {
              fetchPaymentMethods({ docker_id });
            }
          }),
        ]);
      }),
      fetchOwnedNamespaces(),
    ]).then(() => {
      this.setState({ isInitializing: false });
    }).catch(() => {
      this.setState({ isInitializing: false });
    });
  }

  onSelectPage = (e, value) => {
    this.setState({ currentView: value });
  }

  onSelectNamespace = ({ value }) => {
    const {
      currentUser,
      namespaceObjects,
      // actions
      accountFetchUser: fetchUser,
      accountSelectNamespace: selectNamespace,
      billingFetchPaymentMethods: fetchPaymentMethods,
      billingFetchProfile: fetchBillingProfiles,
      billingFetchInvoices: fetchInvoices,
    } = this.props;
    this.setState({ initializing: true });
    const fetchedNamespaces = namespaceObjects.results;
    if (fetchedNamespaces[value]) {
      const userOrg = fetchedNamespaces[value];
      const { type, id: docker_id } = userOrg;
      const isOrg = type === 'Organization';
      return Promise.when([
        fetchInvoices({ docker_id }),
        selectNamespace({ namespace: value }),
        fetchBillingProfiles({ docker_id, isOrg }).then((res) => {
          if (res.value.profile) {
            fetchPaymentMethods({ docker_id });
          }
        }),
      ]).then(() => {
        this.setState({
          isInitializing: false,
        });
      }).catch(() => {
        this.setState({
          isInitializing: false,
        });
      });
    }
    const isOrg = value !== currentUser.username;
    // fetching user info vs. org info locally hits different api's which
    // requires knowing which endpoint to hit.
    // unnecessary in production since v2/users will redirect to v2/orgs
    return fetchUser({ namespace: value, isOrg }).then((userRes) => {
      const { id: docker_id } = userRes.value;
      return Promise.when([
        fetchInvoices({ docker_id }),
        selectNamespace({ namespace: value }),
        fetchBillingProfiles({ docker_id, isOrg })
          .then((billingRes) => {
            if (billingRes.value.profile) {
              fetchPaymentMethods({ namespace: value, docker_id });
            }
          }),
      ]).then(() => {
        this.setState({
          isInitializing: false,
        });
      }).catch(() => {
        this.setState({
          isInitializing: false,
        });
      });
    });
  }

  generateSelectOptions(namespaces) {
    return namespaces.map((option) => {
      return { value: option, label: option };
    });
  }

  render() {
    const {
      currentView,
      isInitializing,
    } = this.state;
    const {
      billingProfiles,
      namespaceObjects,
      ownedNamespaces,
      selectedNamespace,
    } = this.props;
    const isLoading = namespaceObjects.isFetching || billingProfiles.isFetching;
    if (isInitializing || isLoading) {
      return <FullscreenLoading />;
    }
    const selectedUser = namespaceObjects.results[selectedNamespace];
    const billingProfile = billingProfiles.results[selectedUser.id];
    // IF NAMESPACE HAS NO BILLING PROFILE RENDER EMPTY CARD
    if (!billingProfile) {
      // TODO: design - nathan 6/8/16 - MAKE PRETTIER NO CONTENT PAGE
      return (
        <div>
          <div className={css.emptyNav}>
            <AccountSelect
              className={css.select}
              options={this.generateSelectOptions(ownedNamespaces)}
              onSelectChange={this.onSelectNamespace}
              selectedNamespace={selectedNamespace}
            />
          </div>
          <div className={css.content}>
            <Card shadow>No Billing Profile to display</Card>
          </div>
        </div>
      );
    }

    let content;
    if (currentView === PAYMENTS) {
      content = <InvoicesView selectedUser={selectedUser} />;
    } else if (currentView === METHODS) {
      content = (
        <PaymentMethodsView
          selectedUser={selectedUser}
          billingProfile={billingProfile}
        />);
    }
    return (
      <div>
        <div className={css.navigation}>
          <Tabs
            className={css.tabs}
            selected={currentView}
            onSelect={this.onSelectPage}
          >
            <Tab value={PAYMENTS}>Billing History</Tab>
            <Tab value={METHODS}>Payment Methods</Tab>
          </Tabs>
          <AccountSelect
            className={css.select}
            options={this.generateSelectOptions(ownedNamespaces)}
            onSelectChange={this.onSelectNamespace}
            selectedNamespace={selectedNamespace}
          />
        </div>
        <div className={css.content}>
          {content}
        </div>
      </div>
    );
  }
}
