import React, { Component } from 'react';
import svg from './_Svg';

@svg
class ImgIcon extends Component {
  render() {
    return (
      <g>
        <circle fill="none" cx="16" cy="8" r="2" />
        <path fill="none" d="M11 16.51L8.5 13.5 5 18h14l-4.5-6" /><path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-3 3c1.104 0 2 .896 2 2s-.896 2-2 2-2-.896-2-2 .896-2 2-2zM5 18l3.5-4.5 2.5 3.01L14.5 12l4.5 6H5z" />
      </g>
    );
  }
}

export default ImgIcon;
