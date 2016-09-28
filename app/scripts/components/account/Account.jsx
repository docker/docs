'use strict';

import React, { PropTypes, cloneElement } from 'react';
import { Link } from 'react-router';
import FluxibleMixin from 'fluxible-addons-react/FluxibleMixin';
import { SecondaryNav } from 'dux';
import FA from 'common/FontAwesome';
import LiLink from '../common/LiLink';
import AccountSettingsLicensesStore from '../../stores/AccountSettingsLicensesStore';
import EmailNotifStore from '../../stores/EmailNotifStore';
import Route404 from '../common/RouteNotFound404Page.jsx';

let SettingNav = React.createClass({
  displayName: 'SettingsSecondaryNav',
  mixins: [FluxibleMixin],
  contextTypes: {
    getStore: PropTypes.func.isRequired
  },
  statics: {
    storeListeners: {
      onLicensesStoreChange: [AccountSettingsLicensesStore],
      onNotifStoreChange: [EmailNotifStore]
    }
  },
  onLicensesStoreChange: function() {
    let store = this.context.getStore(AccountSettingsLicensesStore);
    this.setState({
      licenseAttempt: store.getAttempt()
    });
  },
  onNotifStoreChange: function() {
    let store = this.context.getStore(EmailNotifStore);
    this.setState({
      notificationAttempt: store.getAttempt()
    });
  },
  getInitialState() {
    return {
      licenseAttempt: false,
      notificationAttempt: false
    };
  },
  _showLicensesLoader() {
    this.context.getStore(AccountSettingsLicensesStore).setAttempt(true);
    this.setState({
      licenseAttempt: true
    });
  },
  _showNotificationsLoader() {
    this.context.getStore(EmailNotifStore).setAttempt(true);
    this.setState({
      notificationAttempt: true
    });
  },
  renderLicenses() {
    var licensesElement;
    if (!this.state.licenseAttempt) {
      licensesElement = 'Licenses';
    } else {
      licensesElement = (<span>Licenses&nbsp;<FA icon='fa-spinner fa-spin'></FA></span>);
    }
    return <LiLink to="/account/licenses/" onClick={this._showLicensesLoader}>{licensesElement}</LiLink>;
  },
  renderNotifications() {
    var notificationsElement;
    if (!this.state.notificationAttempt) {
      notificationsElement = 'Notifications';
    } else {
      notificationsElement = (<span>Notifications&nbsp;<FA icon='fa-spinner fa-spin'></FA></span>);
    }
    return <LiLink to="/account/notifications/" onClick={this._showNotificationsLoader}>{notificationsElement}</LiLink>;
  },
  render() {
    return (
      <SecondaryNav>
        <ul>
          <LiLink to="/account/settings/">Account Settings</LiLink>
          <LiLink to="/account/billing-plans/">Billing & Plans</LiLink>
          <LiLink to="/account/authorized-services/">Linked Accounts & Services</LiLink>
          {this.renderNotifications()}
          {this.renderLicenses()}
        </ul>
      </SecondaryNav>
    );
  }
});

var Account = React.createClass({
  mixins: [FluxibleMixin],
  propTypes: {
    loggedOutElement: PropTypes.element
  },
  contextTypes: {
    getStore: PropTypes.func.isRequired
  },
  render: function() {

    const { JWT, user, location } = this.props;
    const path = location.pathname;
    if (!JWT && path === '/account/billing-plans/create-subscription/') {
      /* Handles the www.docker.com/pricing links to buy hub plans at URLs
       * /account/billing-plans/create-subscription/?plan=index_personal_PLANSIZE
       * when the user is logged out (has special redirects)
       */
       return cloneElement(this.props.children, {
           user: this.props.user,
           JWT: this.props.JWT
         });
    } else if (!JWT) {
      return (
        <Route404 />
      );
    } else {
      return (
        <div className="settings-wrapper">
          <SettingNav />
          {this.props.children && cloneElement(this.props.children, {
            user: this.props.user,
            JWT: this.props.JWT
          })}
        </div>
      );
    }
  }
});

module.exports = Account;
