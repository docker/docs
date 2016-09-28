'use strict';

import React, { Component, PropTypes } from 'react';
import FontAwesome from 'components/common/fontAwesome';
import styles from './copyToClipboard.css';
import css from 'react-css-modules';

@css(styles)
export default class CopyToClipboard extends Component {

  static propTypes = {
    children: PropTypes.node
  };

  copyToClipboard = () => {

    if (document.queryCommandSupported('copy')) {
      // to write to the clipboard you have to create a selection range
      const range = document.createRange();
      // then tell it to select the contents of a node
      range.selectNodeContents(this.refs.copyNode);

      // get the selection and clear all ranges already there
      const s = window.getSelection();
      s.removeAllRanges();
      // add our selected text
      s.addRange(range);

      // copy the text to the clipboard
      const success = document.execCommand('copy');

      // report success
      // what should we do on failure?
      if (success) {
        alert('Copied to clipboard');
      }
    }

  };

  render () {
    return (
      <div styleName='copyToClipboard'>
        <p ref='copyNode'>{ this.props.children }</p>
        <span styleName='clipboard' onClick={ ::this.copyToClipboard }>
          <FontAwesome icon='fa-clipboard'/>
        </span>
      </div>
    );
  }
}
