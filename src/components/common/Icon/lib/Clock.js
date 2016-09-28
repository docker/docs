import React, { Component } from 'react';
import svg from './_Svg';

@svg
class ClockIcon extends Component {
  render() {
    return <path d="M12.5 7H11v6l5.25 3.15.75-1.23-4.5-2.67V7zm-.51-5C6.47 2 2 6.48 2 12s4.47 10 9.99 10C17.52 22 22 17.52 22 12S17.52 2 11.99 2zM12 20c-4.42 0-8-3.58-8-8s3.58-8 8-8 8 3.58 8 8-3.58 8-8 8z" />;
  }
}

export default ClockIcon;
