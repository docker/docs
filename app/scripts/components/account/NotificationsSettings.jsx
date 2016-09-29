'use strict';

import React, { PropTypes } from 'react';
import _ from 'lodash';
import connectToStores from 'fluxible-addons-react/connectToStores';
import classnames from 'classnames';

import { Button, PageHeader } from 'dux';
import OutboundCommunicationStore from '../../stores/OutboundCommunicationStore';
import EmailsStore from '../../stores/EmailsStore';
import resetNotifications from '../../actions/resetNotifications.js';
import saveOutbound from '../../actions/saveOutbound';
import UpdateOutbound from '../../actions/updateOutbound.js';
import EmailNotifForm from './notificationSettings/EmailNotifForm.jsx';

import { SplitSection } from './../common/Sections.jsx';
import styles from './NotificationsSettings.css';
import DocumentTitle from 'react-document-title';

var debug = require('debug')('NotificationSettings');

var Notifications = React.createClass({
  displayName: 'Notifications',
  contextTypes: {
    executeAction: PropTypes.func.isRequired
  },
  propTypes: {
    user: PropTypes.shape({
      username: PropTypes.string
    }),
    JWT: PropTypes.string
  },
  getDefaultProps: function() {
    return {
      user: {},
      JWT: ''
    };
  },
  onOutboundClick: function(e) {
    debug('onOutboundClick');
    var email = e.currentTarget.getAttribute('data-email');
    var list = e.currentTarget.getAttribute('data-list');
    var unsubscribedIndex;
    var subscribedIndex;
    var newSubscribed;
    var newUnsubscribed;
    var newList;
    if (list === 'weekly') {
      /*eslint-disable camelcase */
      unsubscribedIndex = this.props.weeklyDigest.unsubscribed_emails.indexOf(email);
      subscribedIndex = this.props.weeklyDigest.subscribed_emails.indexOf(email);
      newUnsubscribed = _.clone(this.props.weeklyDigest.unsubscribed_emails);
      newSubscribed = _.clone(this.props.weeklyDigest.subscribed_emails);
      if (unsubscribedIndex > -1 && subscribedIndex === -1) {
        newUnsubscribed.splice(unsubscribedIndex, 1);
        newSubscribed.push(email);
      } else {
        newSubscribed.splice(subscribedIndex, 1);
        newUnsubscribed.push(email);
      }
      newList = {
        subscribed_emails: newSubscribed,
        unsubscribed_emails: newUnsubscribed
      };
      this.context.executeAction(UpdateOutbound, {list: 'weekly', data: newList});
      /*eslint-enable camelcase */
    } else if (list === 'beta') {
      unsubscribedIndex = this.props.betaGroup.unsubscribed_emails.indexOf(email);
      subscribedIndex = this.props.betaGroup.subscribed_emails.indexOf(email);
      newUnsubscribed = _.clone(this.props.betaGroup.unsubscribed_emails);
      newSubscribed = _.clone(this.props.betaGroup.subscribed_emails);
      if (unsubscribedIndex > -1 && subscribedIndex === -1) {
        newUnsubscribed.splice(unsubscribedIndex, 1);
        newSubscribed.push(email);
      } else {
        newSubscribed.splice(subscribedIndex, 1);
        newUnsubscribed.push(email);
      }
      newList = {
        unsubscribed_emails: newUnsubscribed,
        subscribed_emails: newSubscribed
      };
      this.context.executeAction(UpdateOutbound, {list: 'beta', data: newList});
    }
  },
  onOutboundSubmit: function(e) {
    e.preventDefault();
    let _this = this;
    let weeklyUns = [];
    let betaUns = [];
    let weekly = _.clone(this.props.weeklyDigest);
    let beta = _.clone(this.props.betaGroup);
    /* eslint-disable camelcase */
    weekly.unsubscribed_emails.forEach(function(email) {
      if (_this.isVerified(email)) {
        weeklyUns.push(email);
      }
    });
    weekly.unsubscribed_emails = weeklyUns;
    beta.unsubscribed_emails.forEach(function(email) {
      if (_this.isVerified(email)) {
        betaUns.push(email);
      }
    });
    beta.unsubscribed_emails = betaUns;
    /* eslint-enable camelcase */
    var outboundData = {
      JWT: this.props.JWT,
      username: this.props.user.username,
      weeklyDigest: weekly,
      betaGroup: beta
    };
    this.context.executeAction(saveOutbound, outboundData);
  },
  onOutboundCancel: function(e) {
    debug('resetting outbound communication');
    this.context.executeAction(resetNotifications, ['outbound']);
  },
  sortEmails: function(emailArray) {
    var _this = this;
    return (_.sortBy(emailArray,
      function(email) {
        return !(_this.isVerified(email));
      })
    );
  },
  isVerified: function(email) {
    return _.pluck(_.filter(this.props.emails, {email: email}), 'verified')[0];
  },
  render: function() {
    var digestEmails = this.sortEmails(this.props.digestEmails);
    var betaEmails = this.sortEmails(this.props.betaEmails);
    return (
      <DocumentTitle title='Notifications - Docker Hub'>
        <div>
          <PageHeader title='Notifications' />
          <div className={'row ' + styles.body}>
            <div className="columns large-12">
              <EmailNotifForm user={this.props.user}
                              JWT={this.props.JWT}/>

              <SplitSection title="Docker Outbound Communication"
                            subtitle={<p>The following settings will control how Docker communicates news, new products, new features, and more.</p>}>
                <form onSubmit={this.onOutboundSubmit}>
                  <OutboundList emailGroups={this.props.weeklyDigest}
                                emailList={digestEmails}
                                type='weekly'
                                isVerified={this.isVerified}
                                onOutboundClick={this.onOutboundClick}>
                    <h5>Docker Weekly</h5>
                    <p>Docker weekly is a newsletter which contains the latest news, updates, and information on releases and other exciting stuff. Sign up!</p>
                  </OutboundList>
                  <OutboundList emailGroups={this.props.betaGroup}
                                emailList={betaEmails}
                                type='beta'
                                isVerified={this.isVerified}
                                onOutboundClick={this.onOutboundClick}>
                    <h5>Docker Beta Group</h5>
                    <p>Become part of our Docker Beta Group and get early access to some Docker products and/or features.</p>
                    <p>Please select which email address (or addresses) you'd like to subscribe:</p>
                  </OutboundList>

                  <div className="row">
                    <div className={'columns large-2 end right ' + styles.button}>
                      <Button size="small" intent='secondary' onClick={this.onOutboundCancel}>Reset</Button>
                    </div>
                    <div className={'columns large-2 end right ' + styles.button}>
                      <Button type="submit" size="small">Save</Button>
                    </div>
                  </div>
                </form>
              </SplitSection>

            </div>
          </div>
        </div>
      </DocumentTitle>
    );
  }
});

var OutboundList = React.createClass({
  contextTypes: {
    executeAction: PropTypes.func.isRequired
  },
  propTypes: {
    emailGroups: PropTypes.shape({
      subscribed_emails: PropTypes.array,
      unsubscribed_emails: PropTypes.array
    }),
    emailList: PropTypes.array,
    emails: PropTypes.array,
    onOutboundClick: PropTypes.func.isRequired,
    type: PropTypes.string.isRequired
  },
  render: function() {
    let _this = this;
    return (
      <div className="Outbound">
        <div className="row">
          <div className="columns large-12">
            {this.props.children}
          </div>
        </div>
        {this.props.emailList.map(function(email) {
          if (_this.props.isVerified(email)) {
            return (
              <div className={'row ' + styles.notification} data-email={email} data-list={_this.props.type} key={email} onClick={_this.props.onOutboundClick}>
                <div className={'columns large-1 ' + styles.checkbox}>
                  <label>
                    <input type="checkbox" name={email} checked={_.includes(_this.props.emailGroups.subscribed_emails, email)} readOnly/>
                  </label>
                </div>
                <div className="columns large-10 large-offset-1">
                  {email}
                </div>
              </div>
            );
          } else {
            return (
              <div className={'row ' + styles.unverifiedNotif} data-email={email} data-list="weekly" key={email}>
                <div className={'columns large-1 ' + styles.checkbox}>
                  <input type="checkbox" name={email} checked={false} readOnly disabled/>
                </div>
                <div className='columns large-10 large-offset-1'>
                  {email} - unverified
                </div>
              </div>
            );
          }
        })}
      </div>
    );
  }
});

export default connectToStores(Notifications,
  [OutboundCommunicationStore, EmailsStore], function({ getStore }, props) {
    return _.merge({}, getStore(OutboundCommunicationStore).getState(), getStore(EmailsStore).getEmails());
  });
