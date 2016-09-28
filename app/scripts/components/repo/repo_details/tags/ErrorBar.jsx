'use strict';

import React, { PropTypes, Component } from 'react';
import styles from './ErrorBar.css';

export default class ErrorBar extends Component {
  render() {
    return <div className={styles.bar}>&nbsp;</div>;
  }
}
