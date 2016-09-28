import React, { Component } from 'react';
import RCSlider from 'rc-slider';
import './slider.css';
import './styles.css';

export default class Slider extends Component {
  render() {
    return (
      <RCSlider {...this.props} />
    );
  }
}
