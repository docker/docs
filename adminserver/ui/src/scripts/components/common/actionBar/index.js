'use strict';

import React, { Component, PropTypes } from 'react';
import styles from './actionBar.css';

export default class ActionBar extends Component {
  static propTypes = {
    className: PropTypes.string,
    children: PropTypes.node.isRequired
  }

  render() {
    return <div className={ `${styles.bar} ${this.props.className}` }>{ this.props.children }</div>;
  }
}
