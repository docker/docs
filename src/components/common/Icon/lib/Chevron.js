import React, { Component } from 'react';
import svg from './_Svg';

@svg
class ChevronIcon extends Component {
  render() {
    return (
      <g>
        <path d="M16.59 8.59L12 13.17 7.41 8.59 6 10l6 6 6-6z" />
        <path d="M0 0h24v24H0z" fill="none" />
      </g>
    );
  }
}

export default ChevronIcon;
