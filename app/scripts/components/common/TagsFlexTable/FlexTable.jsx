'use strict';

import React, { Component, PropTypes } from 'react';
const { bool } = PropTypes;
import styles from './FlexTable.css';
import classnames from 'classnames';

export default class FlexTable extends Component {
  static propTypes = {
    success: bool,
    error: bool,
    isWrappedInACard: bool
  };

  static defaultProps = {
    isWrappedInACard: false
  };

  render() {
    const { error, isWrappedInACard, success } = this.props;
    const tableClass = classnames({
      [styles.flexTable]: true,
      [styles.defaultBorder]: !isWrappedInACard,
      [styles.inACardBorder]: isWrappedInACard,
      [styles.success]: success,
      [styles.error]: error
    });

    return (
      <div className={tableClass}>
        {this.props.children}
      </div>
    );
  }
}
