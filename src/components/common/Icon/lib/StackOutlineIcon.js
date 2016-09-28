import React, { Component, PropTypes } from 'react';
import svg from './_Svg';

const { bool } = PropTypes;

@svg(80)
class StackOutlineIcon extends Component {
  static propTypes = {
    filled: bool,
  }
  render() {
    const { filled } = this.props;
    let fill1 = '#f7f8f9'; // $color-porcelain
    let fill2 = '#e0e4e7'; // $color-geyser

    if (filled) {
      fill1 = '#ccf4f4';
      fill2 = '#00cbca';
    }
    return (
      <g>
        <path fill={fill1} d="M20.5 59.5v-7h55v7h-55zm-8-16v-7h55v7h-55zm-8-16v-7h55v7h-55z" />
        <path fill={fill2} d="M59 21v6H5v-6h54m8 16v6H13v-6h54m8 16v6H21v-6h54M60 20H4v8h56v-8zm8 16H12v8h56v-8zm8 16H20v8h56v-8z" />
      </g>
    );
  }
}

export default StackOutlineIcon;
