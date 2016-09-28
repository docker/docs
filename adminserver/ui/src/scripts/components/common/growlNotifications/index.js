'use strict';

import React, { Component, PropTypes } from 'react';
import cn from 'classnames';
import { connect } from 'react-redux';
import { routerActions } from 'react-router-redux';

import { mapActions } from 'utils';
import * as NotificationActions from 'actions/notifications';
import styles from './notifications.css';

const mapState = (state) => {
  return {
    notifications: state.notifications.growl
  };
};

const {
    push
} = routerActions;
@connect(mapState, mapActions({
    ...NotificationActions,
    push
}))
export default class GrowlNotification extends Component {

  static propTypes = {
    actions: PropTypes.object,
    notifications: PropTypes.array.isRequired
  }

  onClick = (item) => () => {
    if (item.url) {
      this.props.actions.push(item.url);
    }
    this.removeNotification(item);
  }

  onClose = (item) => () => {
    this.removeNotification(item);
  }

  removeNotification = (item) => {
    this.props.actions.removeGrowlNotification(item.id);
  }

  render() {
    return (
      <div className={ styles.container }>
        { this.props.notifications.map( item => {
          // Call the onclose action by default when rendering. This is needed
          // from server-side notifications that we only want to show once (they
          // include an onclose property).
          if (item.onclose) {
            this.props.actions.onClose(item.onclose);
          }

          let key = item.id || item.defaultId;

          return (
            <div
              key={ key }
              className={ cn([styles[item.status], styles.item, item.class]) }
              onClick={ ::this.onClick(item) }>
                <span className={ styles.close } onClick={ ::this.onClose }>&#x2716;</span>
                {
                    item.img ? (
                        <img className={ styles.img } src={ item.img } />
                    ) : null
                }
                <div className={ styles.text }>
                    <div className={ styles.title }>{ item.title }</div>
                    <div>{ item.message }</div>
                </div>
            </div>
          );
        }, this) }
      </div>
    );
  }
}
