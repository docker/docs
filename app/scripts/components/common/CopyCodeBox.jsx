'use strict';

import React, { PropTypes, Component } from 'react';
import { findDOMNode } from 'react-dom';
import styles from './CopyCodeBox.css';
import Button from '@dux/element-button';
const debug = require('debug')('CopyCodeBox');
import repeat from 'lodash/string/repeat';
const { bool, number, string } = PropTypes;

export default class CopyCodeBox extends Component {
  static propTypes = {
    content: string.isRequired,
    dollar: bool,
    lines: number
  }

  static defaultProps = {
    dollar: false,
    lines: 1
  }

  selectContents = (e) => {
    const content = findDOMNode(this.refs.content);
    let range;
    if ( document.selection ) {
        range = document.body.createTextRange();
        range.moveToElementText(content);
        range.select();
    } else if ( window.getSelection ) {
        range = document.createRange();
        range.selectNodeContents(content);
        window.getSelection().removeAllRanges();
        window.getSelection().addRange(range);
    }
  }

  copyContentsToClipboard = (e) => {
    /* Works in Chrome and Firefox 41.x
     * TODO: Does not work in Safari
     */
    this.selectContents(e);
    try {
      const success = document.execCommand('copy');
      debug(`Copy worked: ${success}`);
    } catch (err) {
      debug('Cannot copy.');
    }
  }
  render() {
    const { content, dollar, lines } = this.props;
    let dollarDiv;
    const button = (
      <Button onClick={this.copyContentsToClipboard}
              ghost>
        Copy
      </Button>
    );
    if (dollar) {
      const inner = repeat('$ \n', lines);
      dollarDiv = (
        <div className={styles.dollarBox}>{inner}</div>
      );
    }
    return (
      <div className={styles.wrapper}>
        <div className={styles.copyBox}>
          {dollarDiv}
          <div className={styles.contentBox} ref='content'>{content}</div>
        </div>
        {button}
      </div>
    );
  }
}
