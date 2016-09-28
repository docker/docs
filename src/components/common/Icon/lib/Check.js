import React, { Component } from 'react';
import svg from './_Svg';

@svg
class CheckIcon extends Component {
  render() {
    return (
      <g>
        <path d="M0 0h24v24H0z" fill="none" />
        <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z" />
      </g>
    );
  }
}

export default CheckIcon;
