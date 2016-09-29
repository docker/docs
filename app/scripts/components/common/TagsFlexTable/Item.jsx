'use strict';

import React, { Component, PropTypes } from 'react';
import styles from './Item.css';
const { bool, number } = PropTypes;
import classnames from 'classnames';

export default class Item extends Component {

  static propTypes = {
    grow: number,
    header: bool
  }

  static defaultProps = {
    header: false
  }

  render() {
    const { header, grow } = this.props;

    const itemClass = classnames({
      [styles.flexItem]: true,
      [styles.flexHeaderItem]: header,
      [styles.flexItemGrow1]: !grow || grow === 1,
      [styles.flexItemGrow2]: grow === 2,
      [styles.flexItemGrow3]: grow === 3
    });

    return (
      <div className={itemClass}>
        {this.props.children}
      </div>
    );
  }
}
