import React, { Component, PropTypes } from 'react';
import svg from './_Svg';

const { bool } = PropTypes;

@svg(80)
class ServiceOutlineIcon extends Component {
  static propTypes = {
    filled: bool,
  }
  render() {
    const { filled } = this.props;
    let fill1 = '#f7f8f9'; // $color-porcelain
    let fill2 = '#e0e4e7'; // $color-geyser

    if (filled) {
      fill1 = '#d1eefd';
      fill2 = '#1aaaf8';
    }
    return (
      <g>
        <path fill={fill1} d="M36.5 27.5h7v-7h-7v7zm35.5 9H59.5V8c0-1.9-1.6-3.5-3.5-3.5H24c-1.9 0-3.5 1.6-3.5 3.5v28.5H8c-1.9 0-3.5 1.6-3.5 3.5v32c0 1.9 1.6 3.5 3.5 3.5h64c1.9 0 3.5-1.6 3.5-3.5V40c0-1.9-1.6-3.5-3.5-3.5zm-44.5-25h25v25h-25v-25zm9 57h-25v-25h25v25zm32 0h-25v-25h25v25zm-9-16h-7v7h7v-7zm-32 0h-7v7h7v-7z" />
        <path fill={fill2} d="M56 5c1.6 0 3 1.4 3 3v29h13c1.6 0 3 1.4 3 3v32c0 1.6-1.3 3-3 3H8c-1.6 0-3-1.4-3-3V40c0-1.6 1.4-3 3-3h13V8c0-1.6 1.4-3 3-3h32M27 37h26V11H27v26m16 32h26V43H43v26m-32 0h26V43H11v26m32-48v6h-6v-6h6m16 32v6h-6v-6h6m-32 0v6h-6v-6h6M56 4H24c-2.2 0-4 1.8-4 4v28H8c-2.2 0-4 1.8-4 4v32c0 2.2 1.8 4 4 4h64c2.2 0 4-1.8 4-4V40c0-2.2-1.8-4-4-4H60V8c0-2.2-1.8-4-4-4zM28 36V12h24v24H28zm16 32V44h24v24H44zm-32 0V44h24v24H12zm32-48h-8v8h8v-8zm16 32h-8v8h8v-8zm-32 0h-8v8h8v-8z" />
      </g>
    );
  }
}

export default ServiceOutlineIcon;
