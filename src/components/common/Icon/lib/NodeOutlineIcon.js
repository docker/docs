import React, { Component, PropTypes } from 'react';
import svg from './_Svg';

const { bool } = PropTypes;

@svg(80)
class NodeOutlineIcon extends Component {
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
        <path fill={fill1} d="M32 75.5c-4.1 0-7.5-3.4-7.5-7.5V12c0-4.1 3.4-7.5 7.5-7.5h16c4.1 0 7.5 3.4 7.5 7.5v56c0 4.1-3.4 7.5-7.5 7.5H32zm4-48c-2.5 0-4.5 2-4.5 4.5s2 4.5 4.5 4.5 4.5-2 4.5-4.5-2-4.5-4.5-4.5zm0-16c-2.5 0-4.5 2-4.5 4.5s2 4.5 4.5 4.5 4.5-2 4.5-4.5-2-4.5-4.5-4.5z" />
        <path fill={fill2} d="M36 27c-2.8 0-5 2.2-5 5s2.2 5 5 5 5-2.2 5-5-2.2-5-5-5zm0 9c-2.2 0-4-1.8-4-4s1.8-4 4-4 4 1.8 4 4-1.8 4-4 4zM48 4H32c-4.4 0-8 3.6-8 8v56c0 4.4 3.6 8 8 8h16c4.4 0 8-3.6 8-8V12c0-4.4-3.6-8-8-8zm7 64c0 3.8-3.1 7-7 7H32c-3.9 0-7-3.1-7-7V12c0-3.8 3.1-7 7-7h16c3.9 0 7 3.1 7 7v56zM36 11c-2.8 0-5 2.2-5 5s2.2 5 5 5 5-2.2 5-5-2.2-5-5-5zm0 9c-2.2 0-4-1.8-4-4s1.8-4 4-4 4 1.8 4 4-1.8 4-4 4z" />
      </g>
    );
  }
}

export default NodeOutlineIcon;
