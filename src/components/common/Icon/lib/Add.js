import React, { Component } from 'react';
import svg from './_Svg';

@svg
class AddIcon extends Component {
  render() {
    return <path d="M13 11V5h-2v6H5v2h6v6h2v-6h6v-2" />;
  }
}

export default AddIcon;
