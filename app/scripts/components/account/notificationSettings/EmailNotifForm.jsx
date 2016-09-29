'use strict';
import React, { PropTypes } from 'react';
import _ from 'lodash';
import classnames from 'classnames';
import connectToStores from 'fluxible-addons-react/connectToStores';

import styles from './EmailNotifForm.css';

import EmailNotifStore from '../../../stores/EmailNotifStore.js';
import deleteEmailNotifs from '../../../actions/deleteEmailNotifs';
import resetNotifications from '../../../actions/resetNotifications.js';
import saveEmailNotifs from '../../../actions/saveEmailNotifs';
import updateNotifCheckbox from '../../../actions/updateNotifCheckbox.js';
import { Button } from 'dux';
import { SplitSection } from '../../common/Sections';
var debug = require('debug')('EmailNotifForm');

var EmailNotifForm = React.createClass({
  displayName: 'EmailNotifForm',
  contextTypes: {
    executeAction: PropTypes.func.isRequired,
    getStore: PropTypes.func.isRequired
  },
  propTypes: {
    user: PropTypes.object.isRequired,
    JWT: PropTypes.string.isRequired
  },
  onNotifCheckboxClick: function(cboxType, evt) {
    this.context.executeAction(updateNotifCheckbox, cboxType);
  },
  //TODO: Fix this when we have time to change the backend API. This sucks atm.
  _getNotificationPayload: function(type) {
    var nType = '';
    switch (type) {
      case 'auto':
        nType = 'trusted_build_fail';
        break;
      case 'star':
        nType = 'new_repo_star';
        break;
      case 'comment':
        nType = 'new_repo_comment';
        break;
    }
    return {
      'user': this.props.user.username,
      'notification': nType,
      'last_occurrence': '1970-01-01T00:00:00Z'
    };
  },
  _getNotificationDeletePayload: function(type) {
    switch (type) {
      case 'auto':
        return this.props.autoBuildNotificationID;
      case 'star':
        return this.props.starNotificationID;
      case 'comment':
        return this.props.imgCommentNotificationID;
      default:
        break;
    }
  },
  onNotifSubmit: function(e) {
    e.preventDefault();
    var store = this.context.getStore(EmailNotifStore);
    var changed = store.hasChanged.bind(store);
    //TODO: this NEEDS refactoring once we discuss a better way to handle notification settings! eek.
    if (changed('star')) {
      if (this.props.starNotification) {
        var nStar = this._getNotificationPayload('star');
        this.context.executeAction(saveEmailNotifs, {jwt: this.props.JWT, notification: nStar});
      } else {
        var nStarDel = this._getNotificationDeletePayload('star');
        if(nStarDel > 0) {
          this.context.executeAction(deleteEmailNotifs, {jwt: this.props.JWT, notificationID: nStarDel});
        }
      }
    }
    if (changed('comment')) {
      if (this.props.imgCommentNotification) {
        var nComment = this._getNotificationPayload('comment');
        this.context.executeAction(saveEmailNotifs, {jwt: this.props.JWT, notification: nComment});
      } else {
        var nCommentDel = this._getNotificationDeletePayload('comment');
        if(nCommentDel > 0) {
          this.context.executeAction(deleteEmailNotifs, {jwt: this.props.JWT, notificationID: nCommentDel});
        }
      }
    }
    if (changed('auto')) {
      if (this.props.autoBuildNotification) {
        var nAuto = this._getNotificationPayload('auto');
        this.context.executeAction(saveEmailNotifs, {jwt: this.props.JWT, notification: nAuto});
      } else {
        var nAutoDel = this._getNotificationDeletePayload('auto');
        if(nAutoDel > 0) {
          this.context.executeAction(deleteEmailNotifs, {jwt: this.props.JWT, notificationID: nAutoDel});
        }
      }
    }
  },
  onNotifCancel: function(e) {
    debug('resetting email notifications');
    this.context.executeAction(resetNotifications, ['notifications']);
  },
  render: function() {
    return (
            <SplitSection title="Event Notification by Email"
                          subtitle={<p>The following settings will control how email is sent based on the occurrence of specific events.</p>}>
      <form onSubmit={this.onNotifSubmit}>
        <div className={'row ' + styles.notification} onClick={this.onNotifCheckboxClick.bind(null, 'starNotification')}>
          <div className={'columns large-2 ' + styles.checkbox}>
            <label>
              <input type="checkbox" name="starNotification" checked={this.props.starNotification} readOnly/>
            </label>
          </div>
          <div className="columns large-10">
            Notify me when my Repositories get starred
          </div>
        </div>
        <div className={'row ' + styles.notification} onClick={this.onNotifCheckboxClick.bind(null, 'imgCommentNotification')}>
          <div className={'columns large-2 ' + styles.checkbox}>
            <label>
              <input type="checkbox" name="imgCommentNotification" checked={this.props.imgCommentNotification} readOnly/>
            </label>
          </div>
          <div className="columns large-10">
            Notify me when a comment is posted on my Repositories
          </div>
        </div>
        <div className={'row ' + styles.notification} onClick={this.onNotifCheckboxClick.bind(null, 'autoBuildNotification')}>
          <div className={'columns large-2 ' + styles.checkbox}>
            <label>
              <input type="checkbox" name="autoBuildNotification" checked={this.props.autoBuildNotification} readOnly/>
            </label>
          </div>
          <div className="columns large-10">
            Notify me when an automated build fails
          </div>
        </div>
        <div className="row">
          <div className={'columns large-2 end right ' + styles.button}>
            <Button size="small" intent='secondary' onClick={this.onNotifCancel}>Reset</Button>
          </div>
          <div className={'columns large-2 end right ' + styles.button}>
            <Button type="submit" size="small">Save</Button>
          </div>
        </div>
      </form>
      </SplitSection>
    );
  }
});

export default connectToStores(EmailNotifForm,
  [EmailNotifStore], function({ getStore }, props) {
    return getStore(EmailNotifStore).getState();
  });
