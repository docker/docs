'use strict';

import React from 'react';

import styles from './style.css';

export default class DescriptParagraph extends React.Component {

  static propTypes = {
    children: React.PropTypes.node
  }

  render() {
    return (
      <div className={ styles.container }>
        <p>{ this.props.children }</p>
      </div>
    );
  }
}
