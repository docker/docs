import React, { Component } from 'react';
import svg from './_Svg';

@svg
class Return extends Component {
  render() {
    return (
      <g fill="none" fillRule="evenodd">
        <path d="M0 0h24v24H0V0z" />
        <path d="M19 8v4H5.83l3.58-3.59L8 7l-6 6 6 6 1.41-1.41L5.83 14H21V8h-2z" fill="#000" />
      </g>
    );
  }
}

export default Return;
