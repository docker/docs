import React, { Component } from 'react';
import svg from './_Svg';

@svg(80)
class AngularIcon extends Component {
  render() {
    return (
      <g>
        <path fill="#e23237" d="M15 22.3l24.7-8.8L65 22.1l-4.1 32.6-21.2 11.8-20.9-11.6L15 22.3z" />
        <path fill="#b52e31" d="M65 22.1l-25.3-8.6v53l21.2-11.7L65 22.1z" />
        <path fill="#fff" d="M39.7 19.7L24.3 53.9l5.7-.1 3.1-7.7h13.8l3.3 7.8 5.5.1c.1-.1-16-34.3-16-34.3zm.1 11L45 41.5h-9.8l4.6-10.8z" />
      </g>
    );
  }
}

export default AngularIcon;
