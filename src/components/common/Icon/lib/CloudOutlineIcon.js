import React, { Component, PropTypes } from 'react';
import svg from './_Svg';

const { bool } = PropTypes;

@svg(80)
class CloudOutlineIcon extends Component {
  static propTypes = {
    filled: bool,
  }
  render() {
    const { filled } = this.props;
    let fill1 = '#f7f8f9'; // $color-porcelain
    let fill2 = '#e0e4e7'; // $color-geyser

    if (filled) {
      fill1 = '#d1eefd';
      fill2 = '#1aaaf8';
    }
    return (
      <g>
        <path fill={fill1} d="M21.8 64.5c-10.1 0-18.2-8.2-18.2-18.2 0-9.3 7-17.1 16.3-18.1h.3l.1-.2c4-7.6 11.8-12.4 20.3-12.4 10.9 0 20.4 7.8 22.5 18.5l.1.4h.4c7.9.5 14.1 7.1 14.1 15 0 8.3-6.8 15.1-15.1 15.1l-40.8-.1z" />
        <path fill={fill2} d="M40.5 16c10.7 0 19.9 7.6 22 18.1l.1.8.8.1C71 35.4 77 41.8 77 49.4c0 8-6.6 14.6-14.6 14.6H21.8C12 64 4 56 4 46.2c0-9.1 6.8-16.7 15.8-17.6l.5-.1.2-.5c4-7.4 11.6-12 20-12m0-1c-9 0-16.9 5.1-20.8 12.6-9.4 1-16.7 9-16.7 18.6C3 56.6 11.4 65 21.8 65h40.6C71 65 78 58 78 49.4c0-8.2-6.4-14.9-14.5-15.5C61.3 23.1 51.9 15 40.5 15z" />
      </g>
    );
  }
}

export default CloudOutlineIcon;
