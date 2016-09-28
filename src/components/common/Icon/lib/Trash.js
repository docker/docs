import React, { Component } from 'react';
import svg from './_Svg';

@svg
class TrashIcon extends Component {
  render() {
    return <path d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zm9.5-15l-1-1h-5l-1 1H5v2h14V4h-3.5z" />;
  }
}

export default TrashIcon;
