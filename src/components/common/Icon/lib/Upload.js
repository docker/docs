import React, { Component } from 'react';
import svg from './_Svg';

@svg
class UploadIcon extends Component {
  render() {
    return (
      <g fill-rule="evenodd">
        <path d="M0 0h24v24H0V0z" fill="none" />
        <path d="M18.125 9.908C17.558 7.104 15.033 5 12 5 9.592 5 7.5 6.332 6.458 8.283 3.95 8.543 2 10.613 2 13.125 2 15.815 4.242 18 7 18h10.833c2.3 0 4.167-1.82 4.167-4.063 0-2.144-1.708-3.883-3.875-4.03zM14 12v4h-4v-4H8l4-4 4 4h-2z" />
      </g>
    );
  }
}

export default UploadIcon;
