'use strict';
import React, { Component, PropTypes } from 'react';
import { VelocityTransitionGroup } from 'velocity-react';
// Some issues getting InlineSVG loaded
const InlineSVG = require('svg-inline-react');

import styles from './Options.css';

export default class PipelineOptions extends Component {
  static propTypes = {
    isActive: PropTypes.bool.isRequired,
    toggleHistory: PropTypes.func.isRequired,
    delete: PropTypes.func.isRequired
  }

  renderMenu() {
    return (
      <div className={styles.options_showMenu} key='menu'>
        <div className={styles.menu}>
          <div
            className={styles.menuButton}
            role='button'
            title='View history'
            onClick={this.props.toggleHistory}
          >
            View history
          </div>
          <div
            className={[
              styles.menuButton,
              styles.menuButtonDanger
            ].join(' ')}
            role='button'
            title='Delete'
            onClick={this.props.delete}
          >
            Delete
          </div>
        </div>
      </div>
    );
  }

  renderHistory() {
    return (
      <div
        className={styles.history}
        key='hideHistory'
        role='button'
        title='Hide history'
        onClick={this.props.toggleHistory}
      >
        <div className={styles.historyClose}>
          <InlineSVG src={require('../plus.svg')}/>
        </div>
        Hide history
      </div>
    );
  }

  render() {
    return (
      <div className={styles.options}>
        { this.props.isActive ? this.renderHistory() : this.renderMenu() }
      </div>
    );
  }
}
