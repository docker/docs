'use strict';

import React, { Component } from 'react';
import styles from './Header.css';

export default class Header extends Component {

  render() {
    return (
      <div className={styles.flexHeader}>
        {this.props.children}
      </div>
    );
  }
}
