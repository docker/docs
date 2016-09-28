'use strict';

import React, { PropTypes } from 'react';
const {
  string,
  bool,
  element,
  object,
  func
} = PropTypes;
import { connect } from 'react-redux';
import ui from 'redux-ui';
import autoaction from 'autoaction';

// Components
import Nav from 'components/common/nav';
import GrowlNotifications from 'components/common/growlNotifications';
import BannerNotifications from 'components/common/bannerNotifications';
import LeftNav from 'components/common/leftNav';
import { Login } from 'components/scenes';
import { Modal } from 'components/common/modal';
import * as notificationActions from 'actions/notifications';
import * as updateActions from 'actions/updates';
import * as settingsActions from 'actions/settings';
import * as metricActions from 'actions/metrics';

import { createStructuredSelector } from 'reselect';
import { isAdminSelector, isLoggedInSelector } from 'selectors/users';
import { getAuthMethod, getEnabledUpgrades } from 'selectors/settings';
import { mapActions } from 'utils';

import styles from './app.css';
import css from 'react-css-modules';

const mapState = createStructuredSelector({
  isAdmin: isAdminSelector,
  isLoggedIn: isLoggedInSelector,
    authMethod: getAuthMethod,
    enabledUpgrades: getEnabledUpgrades
});

@connect(mapState, mapActions({
  updates: updateActions,
  settings: settingsActions,
  notifications: notificationActions
}))

@autoaction({
  sendClientAnalytics: []
}, metricActions)

@ui({
  key: 'app',
  state: {
    navExpanded: true
  }
})
@css(styles)
export default class App extends React.Component {
  static displayName = 'App'

  static propTypes = {
    actions: object,
    settings: object,
    children: element.isRequired,

    isAdmin: bool,
    isLoggedIn: bool,
    authMethod: string,
    enabledUpgrades: bool,

    ui: object,
    updateUI: func
  }

  componentWillMount() {
    if (this.props.isAdmin) {
      // Check settings so that we know whether to poll for updates across the
      // entire app.
      this.props.actions.settings.getSettings();
    }
  }

  componentWillReceiveProps(next) {
    // Only start polling for updates if we haven't disabled upgrades.
    // We query settings in componentWillMount() - when we receive settings
    // properties this will be called.
    if (next.enabledUpgrades) {
      this.props.actions.updates.pollForUpdates();
    }

    if (next.authMethod === 'none') {
      this.props.actions.notifications.addBannerNotification({
        id: 'auth',
        message: 'Warning: you have no authentication set up',
        class: 'alert',
        url: '/admin/settings/auth'
      });
    }
  }

  render() {
    const {
      isLoggedIn,
      ui: {
        navExpanded
      }
    } = this.props;

    return (
      <div className={ styles.wrapper }>
        <Modal>
          <Nav isLoggedIn={ isLoggedIn } navExpanded={ navExpanded } />
          <div styleName={ navExpanded ? 'notificationBannerExpanded' : 'notificationBanner' } >
            <BannerNotifications />
          </div>
          <div styleName={ navExpanded ? 'contentExpanded' : 'content' }>
            { isLoggedIn ? <LeftNav isAdmin={ this.props.isAdmin } navExpanded={ navExpanded } /> : undefined }
            { isLoggedIn ? this.props.children : <Login /> }
          </div>
          <GrowlNotifications />
        </Modal>
      </div>
    );
  }
}
