import React, { Component } from 'react';
import svg from './_Svg';

@svg
class WarningIcon extends Component {
  render() {
    return (
      <g>
        <path fill="none" d="M0 0h24v24H0V0z" />
        <path d="M1 21h22L12 2 1 21zm12-3h-2v-2h2v2zm0-4h-2v-4h2v4z" />
      </g>
    );
  }
}

export default WarningIcon;
