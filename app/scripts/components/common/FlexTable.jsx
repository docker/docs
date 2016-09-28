'use strict';

import React, { Component, PropTypes } from 'react';
const { number, bool, func, string } = PropTypes;
import includes from 'lodash/collection/includes';
import styles from './FlexTable.css';
import classnames from 'classnames';

export class FlexTable extends Component {
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

export class FlexRow extends Component {
  static propTypes = {
    onMouseOver: func,
    onMouseLeave: func,
    onClick: func,
    selectable: bool
  }

  render() {
    const {
      onClick,
      onMouseOver,
      onMouseLeave,
      selectable
    } = this.props;
    return (
      <div className={selectable ? styles.selectableFlexRow : styles.flexRow}
           onClick={onClick}
           onMouseOver={onMouseOver}
           onMouseLeave={onMouseLeave}>
        {this.props.children}
      </div>
    );
  }
}

export class FlexHeader extends Component {

  render() {
    return (
      <div className={styles.flexHeader}>
        {this.props.children}
      </div>
    );
  }
}

export class FlexItem extends Component {

  static propTypes = {
    grow: number,
    end: bool,
    noPadding: bool
  }

  static defaultProps = {
    grow: 1,
    noPadding: false
  }

  render() {
    //Optional itemClass usually used for setting widths
    //noPadding removes top and bottom padding
    const { grow, end, noPadding } = this.props;

    const itemClass = classnames({
      [styles.flexItem]: true,
      [styles.flexEnd]: end,
      [styles.flexItemPadding]: !noPadding,
      [styles[`flexItemGrow${grow}`]]: includes([1, 2, 3, 4, 5, 6], grow)
    });

    return (
      <div className={itemClass}>
        {this.props.children}
      </div>
    );
  }
}
