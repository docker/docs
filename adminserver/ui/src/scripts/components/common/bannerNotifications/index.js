'use strict';

import React, { Component, PropTypes } from 'react';
import classnames from 'classnames';
import { connect } from 'react-redux';
import { routerActions } from 'react-router-redux';

import { mapActions } from 'utils';
import * as NotificationActions from 'actions/notifications';
import styles from './notifications.css';
import css from 'react-css-modules';

const mapState = (state) => {
  return {
    notifications: state.notifications.banner
  };
};

const {
    push
} = routerActions;
@connect(mapState, mapActions({
    ...NotificationActions,
    push
}))
@css(styles)
export default class BannerNotification extends Component {

  static propTypes = {
    actions: PropTypes.object,
    notifications: PropTypes.array.isRequired
  }

  onClick = (item) => () => {
    if (item.url) {
      this.props.actions.push(item.url);
    }
  }

  onClose = (item) => () => {
    this.removeNotification(item);
  }

  removeNotification = (item) => {
    this.props.actions.removeBannerNotification(item.id);
  }

  render() {
    return (
      <div>
        { this.props.notifications.map( item => {
          let clickableClass = {};
          clickableClass[styles.clickable] = !!item.url;
          return (
            <div
              key={ item.id }
              className={ classnames(styles.banner, styles[item.class], clickableClass) }
              onClick={ ::this.onClick(item) }>
                { item.img && <img styleName='img' src={ item.img } alt='' /> }
                <span styleName='message'>{ item.message }</span>
            </div>
          );
        }, this) }
      </div>
    );
  }
}
