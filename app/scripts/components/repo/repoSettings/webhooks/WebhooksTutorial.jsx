'use strict';
import React, { Component, PropTypes } from 'react';
const { func } = PropTypes;

// Some issues getting InlineSVG loaded
const InlineSVG = require('svg-inline-react');

import styles from './WebhooksTutorial.css';

export default class WebhooksTutorial extends Component {
  static propTypes = {
    addWebhook: func.isRequired
  }
  render() {
    return (
      <div className={styles.tutorial}>
        <h1 className={styles.headline}>Create a Webhook</h1>
        <button
          className={styles.addWebhook}
          title="Add a webhook"
          onClick={this.props.addWebhook}
        >
          <div className={styles.webhookIcon}>
            <InlineSVG src={require('./webhook.svg')}/>
          </div>
          <div className={styles.plusIcon}>
            <InlineSVG src={require('./plus.svg')}/>
          </div>
        </button>
        <p className={styles.details}>
          A webhook is an HTTP call-back triggered by a specific event.
          You can create a single webhook to start and connect multiple
          webhooks to further build out your workflow.
        </p>
      </div>
    );
  }
}
