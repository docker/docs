import React, { Component } from 'react';
import svg from './_Svg';

@svg
class BurgerIcon extends Component {
  render() {
    return (
      <g>
        <path d="M0 0h24v24H0z" fill="none" />
        <path d="M3 18h18v-2H3v2zm0-5h18v-2H3v2zm0-7v2h18V6H3z" />
      </g>
    );
  }
}

export default BurgerIcon;
