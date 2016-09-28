'use strict';

import React, { Component, PropTypes } from 'react';
import styles from './radioGroup.css';
import SVG from 'components/common/svg';
import css from 'react-css-modules';

@css(styles)
export default class Radio extends Component {
  static propTypes = {
    active: PropTypes.bool.isRequired,
    styles: PropTypes.object
  }

  render() {

    /*
     Safari needs the following properties set as HTML attributes and cannot derive
     them from CSS:
     - cx
     - cy
     - r

     Colors and other aesthetic information are in the associated CSS file.
     */
    const sharedProps = {
      cx: '11px',
      cy: '11px'
    };

    return (

      <SVG
        className={ this.props.styles.svg }
        viewBox='0 0 22 22'>
        { this.props.active ?
          <g>
            <circle className={ styles.active } {...sharedProps} r='9'></circle>
            <circle className={ styles.inner } {...sharedProps} r='4'></circle>
          </g> :
          <circle className={ styles.circle } {...sharedProps} r='9'></circle> }
      </SVG>
    );
  }
}
