import React, { Component } from 'react';
import { Slider } from 'common';
import asExample from '../../asExample';
import mdHeader from './header.md';
import mdApi from './api.md';

@asExample(mdHeader, mdApi)
export default class SliderDoc extends Component {
  render() {
    return (
      <div className="clearfix">
        <h3>Slider</h3>
        <br />
        <br />
        <Slider min={0} max={1000} step={16} />
        <br />
      </div>
    );
  }
}
