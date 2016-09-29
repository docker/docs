'use strict';

import React, { Component } from 'react';
import styles from './Row.css';

export default class Row extends Component {

  render() {
    return (
      <div className={styles.flexRow}>
        {this.props.children}
      </div>
    );
  }
}
