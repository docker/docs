'use strict';

import React, { Component, PropTypes } from 'react';
import { mkAvatarForNamespace } from 'utils/avatar';
const { func, object, string } = PropTypes;

import styles from './Gravatar.css';

export default class GravatarOption extends Component {
  static propTypes = {
    className: string,
    mouseDown: func,
    mouseEnter: func,
    mouseLeave: func,
    option: object.isRequired
  }

  handleMouseDown = (event) => {
    event.preventDefault();
    event.stopPropagation();
    this.props.onSelect(this.props.option, event);
  }

  handleMouseEnter = (event) =>{
    this.props.onFocus(this.props.option, event);
  }

  handleMouseMove = (event) => {
    if (this.props.focused) {
      return;
    }
    this.props.onFocus(this.props.option, event);
  }

  render () {
    const { label, value } = this.props.option;
    return (
      <div className={this.props.className}
           onMouseEnter={this.handleMouseEnter}
           onMouseMove={this.handleMouseMove}
           onMouseDown={this.handleMouseDown}>
        <img className={styles.gravatar} src={mkAvatarForNamespace(value)}/>
        {label}
      </div>
    );
  }

}
