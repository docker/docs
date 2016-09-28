'use strict';

import React, { Component, PropTypes } from 'react';
import cn from 'classnames';
import styles from './tabs.css';

export default class Tabs extends Component {
  static propTypes = {
    children: PropTypes.node.isRequired,
    header: PropTypes.bool,
    sidebar: PropTypes.bool
  }

  static defaultProps = {
    header: false,
    sidebar: false
  }

  render() {
    const { children, header, sidebar } = this.props;

    return (
      <div className={ cn({ [styles.header]: header, [styles.sidebar]: sidebar }) }>
        <ul className={ styles.container }>
          { children }
        </ul>
      </div>
    );
  }
}
