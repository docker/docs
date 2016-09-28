'use strict';

import React, { Component, PropTypes } from 'react';
const { func, bool } = PropTypes;
import styles from './toggle.css';
import SVG from 'components/common/svg';

export ToggleWithLabel from './toggleWithLabel.js';

export default class ToggleSwitch extends Component {

  static propTypes = {
    on: bool.isRequired,
    onClick: func
  }

  render() {
    const { on } = this.props;

    return (

      <SVG
        onClick={ this.props.onClick }
        className={ on ? styles.ToggleSwitchOn : styles.ToggleSwitch }>
        <rect
          x='7'
          y='7'
          rx='8'
          ry='8'
          width='34'
          height='14'/>
        <text
          x={ on ? 12 : 24 }
          y='17'>{ on ? 'on' : 'off' }</text>
        <circle
          cx={ on ? 34 : 14 }
          cy='14'
          r='10'></circle>
      </SVG>

    );
  }
}
