'use strict';

import React, { PropTypes } from 'react';
const { string, shape, object, func } = PropTypes;
import { Link } from 'react-router';
import { FullSection } from '../../common/Sections.jsx';
import FA from '../../common/FontAwesome';
import styles from './EnterpriseSubscriptions.css';
import CopyCodeBox from 'common/CopyCodeBox';
import classnames from 'classnames';
import {DEB, RPM} from 'common/data/csEngineInstructions';
const debug = require('debug')('EnterpriseSubscriptions');

var EnterpriseSubscriptions = React.createClass({
  propTypes: {
    currentPlan: shape({
      id: string,
      plan: string
    }),
    stopSubscription: func.isRequired,
    unsubscribing: string,
    user: object,
    history: object.isRequired
  },
  getInitialState() {
    return {
      confirmAction: ''
    };
  },
  selectAction(plan) {
    return (e) => {
      e.preventDefault();
      this.setState({confirmAction: plan});
    };
  },
  cancelSelectPlan: function(e) {
    e.preventDefault();
    this.setState({confirmAction: ''});
  },
  moreInfo: function(e) {
    e.preventDefault();
    this.props.history.pushState(null, '/enterprise/');
  },
  purchaseCloud: function(e) {
    e.preventDefault();
    this.props.history.pushState(null, '/enterprise/cloud-starter/');
  },
  render: function() {
    let owned;
    let price;
    let subInfo;
    let cloudActionButton;
    let curlBody;

    // Hide Cloud Starter for anyone who has not already purchased it
    if (this.props.currentPlan.package !== 'cloud_starter') {
      return null;
    }

    owned = (
      <h5>Currently Subscribed&nbsp;<span className={styles.check}><FA icon="fa-check-circle-o"/></span></h5>
    );
    if (this.state.confirmAction === 'cloud') {
      cloudActionButton = (
        <div className={styles.flexItem}>
          <a onClick={this.props.stopSubscription}>Confirm</a>&nbsp;or&nbsp;
          <a className={styles.cancel} onClick={this.cancelSelectPlan}>Cancel</a>
        </div>
      );
    } else {
      cloudActionButton = (
        <div className={styles.flexItem}>
          <a className={styles.download} onClick={this.selectAction('download')}>Download</a>
          <a className={styles.cancel} onClick={this.selectAction('cloud')}>Remove Subscription</a>
        </div>
      );
    }
    if (this.props.unsubscribing === 'package' || this.props.unsubscribing === 'subscription') {
      cloudActionButton = (<div className={styles.flexItem}>Removing Subscription <FA icon='fa-spinner fa-spin'/></div>);
    }

    let cloudStarter = (
      <div className={styles.flexItem}>
        <div className={styles.cloudStarter}>
          <h4 className={styles.title} onClick={this.moreInfo}><FA icon='fa-cloud'/>&nbsp;Cloud Starter</h4>
          {owned}
        </div>
      </div>
    );

    if (this.state.confirmAction === 'download') {
      cloudActionButton = null;
      cloudStarter = null;
      curlBody = (
        <div className={styles.downloadFlexItem}>
          <div className={styles.curlWrap}>
            <div className='row'>
              <div className={'columns large-12 ' + styles.curlHelp }>
                * Copy and run either the RPM or the DEB specific instructions in your terminal.
                <a className={styles.cancel} onClick={this.cancelSelectPlan}><FA icon='fa-times'/></a>
              </div>
            </div>
            <div className='row'>
              <div className='columns large-12'>
                RPM
                <CopyCodeBox content={RPM}
                             dollar={true}
                             lines={4} />
              </div>
            </div>
            <div className='row'>
              <div className='columns large-12'>
                DEB
                <CopyCodeBox content={DEB}
                             dollar={true}
                             lines={5} />

                <span className={styles.curlHelp}>
                  For more details about this installation,&nbsp;
                  <a href='http://docs.docker.com/docker-trusted-registry/install/install-csengine/'
                     target='_blank'>
                    view our documentation
                  </a>
                </span>
              </div>
            </div>
          </div>
        </div>
      );
    } else {
      price = (
        <div className={styles.flexItem}>$150/mo</div>
      );
      subInfo = (
        <div className={styles.flexItem}>
          <div>
            20 Private Repositories <br/>
            10 Docker Engines <br/>
            Email Support
          </div>
        </div>
      );
    }
    return (
      <FullSection title='Package Subscriptions'>
        <div className={'columns large-12 ' + styles.flexbox}>
          <div className={styles.flexRow}>
            {cloudStarter}
            {price}
            {subInfo}
            {cloudActionButton}
            {curlBody}
          </div>
        </div>
      </FullSection>
    );
  }
});

module.exports = EnterpriseSubscriptions;
