import React, { Component } from 'react';
import svg from './_Svg';

@svg
export default class Analytics extends Component {
  render() {
    return (
      <g>
        <path d="M15 9v12h2V9h-2zM7 9v12h2V9H7zm6 12V3h-2v18h2zm8 0V3h-2v18h2zM3 21h2V3H3v18z" />
        <path fill="none" d="M0 24V0h24v24z" />
      </g>
    );
  }
}
