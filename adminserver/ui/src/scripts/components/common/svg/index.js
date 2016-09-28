'use strict';

import React, { Component, PropTypes } from 'react';
const { string, node, func, number, oneOfType } = PropTypes;

export default class SVG extends Component {

  static propTypes = {
    width: oneOfType([number, string]),
    height: oneOfType([number, string]),
    viewBox: string,
    styleName: string,
    children: node,
    className: string,
    onClick: func
  }

  render () {
    const {
      width,
      height,
      viewBox,
      children,
      className,
      onClick
      } = this.props;
    return (
      <svg
        onClick={ onClick ? onClick : undefined }
        className={ className ? className : undefined }
        viewBox={ viewBox ? viewBox : undefined }
        width={ width ? width : undefined }
        height={ height ? height : undefined }
        xmlns='http://www.w3.org/2000/svg'
        version='1.1'>
        { children }
      </svg>
    );
  }

}
