'use strict';
import React, { Component, PropTypes } from 'react';
import FA from 'common/FontAwesome';

import styles from './Status.css';

export default class Status extends Component {
  static propTypes = {
    status: PropTypes.oneOf(['init', 'pending', 'success', 'error', 'failure']),
    onClick: PropTypes.func
  }

  render() {
    const klasses = [styles.status];

    if (this.props.onClick) {
      klasses.push(styles.clickable);
    }

    let icon, status;

    switch(this.props.status) {
      case 'init':
      case 'pending':
        icon = <FA icon='fa-clock-o'/>;
        status = 'pending';
        break;
      case 'success':
        icon = <FA icon='fa-check'/>;
        status = 'success';
        break;
      case 'error':
      case 'failure':
        icon = <FA icon='fa-exclamation'/>;
        status = 'error';
        break;
      default:
        icon = <FA icon='fa-question'/>;
        status = 'unknown';
        break;
    }

    klasses.push(styles[`status_${status}`]);

    return (
      <div
        title={`Status: ${status}`}
        className={klasses.join(' ')}
        onClick={this.props.onClick}
      >
        {icon}
        {status}
      </div>
    );
  }
}
