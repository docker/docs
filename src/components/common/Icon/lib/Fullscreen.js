import React, { Component } from 'react';
import svg from './_Svg';

@svg
class Fullscreen extends Component {
  render() {
    return <path d="M5 10h2V7h3V5H5v5zm2 4H5v5h5v-2H7v-3zm10 3h-3v2h5v-5h-2v3zM14 5v2h3v3h2V5h-5z" />;
  }
}

export default Fullscreen;
