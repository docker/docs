import React, { Component } from 'react';
import svg from './_Svg';

@svg
export default class BaseImages extends Component {
  render() {
    return (
      <g>
        <path fill="none" d="M0 0h24v24H0z" />
        <path d="M14 6H8v2h6V6zm4.006-4H5.994C4.894 2 4 2.89 4 3.99v16.02C4 21.102 4.894 22 5.994 22h12.012c1.1 0 1.994-.89 1.994-1.99V3.99C20 2.898 19.105 2 18.006 2zM18 20H6V4h12v16zm-5-10H8v2h5v-2z" />
      </g>
    );
  }
}
