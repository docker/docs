'use strict';

import React, { PropTypes, Component } from 'react';
import hljs from 'highlight.js';
const debug = require('debug')('Code');
import styles from './Code.css';

function renderCode(str) {
  const { value } = hljs.highlightAuto(str);
  return value;
}

export default class Code extends Component {
  render() {
    if(!this.props.children) {
      return null;
    }

    const html = {
      __html: renderCode(this.props.children)
    };

    return (
      <div className={styles.code}>
      <div className='hljs'
           dangerouslySetInnerHTML={html}></div>
      </div>
    );
  }
}
