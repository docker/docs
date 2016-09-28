import React, { Component } from 'react';
import svg from './_Svg';

@svg
export default class Infrastructure extends Component {
  render() {
    return (
      <g>
        <path fill="none" d="M0 0h24v24H0z" />
        <path d="M15.998 3C14.888 3 14 3.893 14 4.995L10 5c0-1.107-.895-2-1.998-2H3.998C2.888 3 2 3.893 2 4.995v14.01C2 20.107 2.895 21 3.998 21h4.004C9.112 21 10 20.107 10 19.005V7h4v12.005c0 1.102.895 1.995 1.998 1.995h4.004c1.11 0 1.998-.893 1.998-1.995V4.995A1.997 1.997 0 0 0 20.002 3h-4.004zM5 11a1 1 0 1 1 0-2 1 1 0 0 1 0 2zm0-4a1 1 0 1 1 0-2 1 1 0 0 1 0 2zm12 4a1 1 0 1 1 0-2 1 1 0 0 1 0 2zm0-4a1 1 0 1 1 0-2 1 1 0 0 1 0 2z" />
      </g>
    );
  }
}
