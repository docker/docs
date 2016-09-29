'use strict';
import React, { PropTypes, Component } from 'react';
import FA from 'common/FontAwesome';
import styles from './BuildStatus.css';
const { number } = PropTypes;

/**
 * possible build status codes
 * -2: 'exception',
 * -1: 'error',
 * 0: 'pending',
 * 1: 'claimed',
 * 2: 'started',
 * 3: 'cloned',
 * 4: 'readme',
 * 5: 'dockerfile',
 * 6: 'built',
 * 7: 'bundled',
 * 8: 'uploaded',
 * 9: 'pushed',
 * 10: 'done',
 * 11: 'queued'
 */

export default class BuildStatus extends Component {
  static propTypes = {
    status: number.isRequired
  };

  render() {
    const { status } = this.props;

    switch (status) {
      case -4:
        return (
          <span className={styles.canceled}>
            <FA fixedWidth={true} icon='fa-ban' /> Canceled
          </span>
        );
      case -2:
      case -1:
        return (
          <span className={styles.warning}>
            <FA fixedWidth={true} icon='fa-exclamation' /> Error
          </span>
        );
      case 0:
      case 1:
      case 11:
        return (
          <span className={styles.primary}>
            <FA fixedWidth={true} icon='fa-clock-o' /> Queued
          </span>
        );
      case 2:
      case 3:
      case 4:
      case 5:
        return (
          <span className={styles.primary}>
            <FA fixedWidth={true} icon='fa-refresh' /> Building
          </span>
        );
      case 6:
      case 7:
      case 8:
      case 9:
        return (
          <span className={styles.primary}>
            <FA fixedWidth={true} icon='fa-refresh' /> Building
          </span>
        );
      case 10:
        return (
          <span className={styles.success}>
            <FA fixedWidth={true} icon='fa-check' /> Success
          </span>
        );
      default:
        return (
          <span className={styles.primary}>
            <FA fixedWidth={true} icon='fa-question' /> Unknown
          </span>
        );
    }
  }
}
