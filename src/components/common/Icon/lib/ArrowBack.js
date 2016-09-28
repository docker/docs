import React, { Component } from 'react';
import svg from './_Svg';

@svg
class ArrowBack extends Component {
  render() {
    return (
      <g>
        <path fill="none" d="M-1 0h24v24H-1V0z" />
        <path d="M19 11H6.83l5.59-5.59L11 4l-8 8 8 8 1.41-1.41L6.83 13H19v-2z" />
      </g>
    );
  }
}

export default ArrowBack;
