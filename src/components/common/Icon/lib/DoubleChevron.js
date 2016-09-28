import React, { Component } from 'react';
import svg from './_Svg';

@svg
class DoubleChevronIcon extends Component {
  render() {
    return (
      <g fillRule="evenodd">
        <path d="M12 7.41L13.41 6l6 6-6 6L12 16.59 16.58 12 12 7.41z" />
        <path d="M5 7.41L6.41 6l6 6-6 6L5 16.59 9.58 12 5 7.41z" />
      </g>
    );
  }
}

export default DoubleChevronIcon;
