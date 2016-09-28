'use strict';

import React, { Component, PropTypes } from 'react';
const { oneOf, string, bool } = PropTypes;
import classnames from 'classnames';
import _ from 'lodash';

const sizes = ['lg', '2x', '3x', '4x', '5x'];
const animations = ['spin', 'pulse'];
const flips = ['horizontal', 'vertical'];
const rotations = [90, 180, 270];

export default class FontAwesome extends Component {
  static propTypes = {
    animate: oneOf(animations),
    className: string,
    fixedWidth: bool,
    flip: oneOf(flips),
    icon: string.isRequired,
    invert: bool,
    rotate: oneOf(rotations),
    size: oneOf(sizes),
    stack: bool
  }

  render() {

    const {
      animate,
      className,
      fixedWidth,
      flip,
      icon,
      invert,
      rotate,
      size,
      stack,
      ...props
      } = this.props;

    const classes = classnames({
      [className]: className,
      'fa': true,
      'fa-fw': fixedWidth,
      'fa-inverse': invert,
      [icon]: true,
      [`fa-${animate}`]: _.includes(animations, animate),
      [`fa-${size}`]: _.includes(sizes, size) && !stack,
      [`fa-flip-${flip}`]: _.includes(flips, flip),
      [`fa-rotate-${rotate}`]: _.includes(rotations, rotate),
      [`fa-stack-${size}`]: _.includes(sizes, size) && stack
    });

    return (<i { ...props } className={ classes } />);
  }
}

