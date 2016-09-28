import React, { Component } from 'react';
import { Uptime } from 'common';
import asExample from '../../asExample';
import mdHeader from './header.md';
import mdApi from './api.md';

const d1 = new Date();
const d2 = new Date();
const d3 = new Date();
const d4 = new Date();
const d5 = new Date();

d2.setMinutes(d2.getMinutes() - 5);
d3.setHours(d3.getHours() - 2);
d4.setMonth(d4.getMonth() - 3);
d5.setFullYear(d5.getFullYear() - 1);

const propsP = {
  style: {
    display: 'inline-block',
    minWidth: '120px',
    marginRight: '18px',
  },
};

@asExample(mdHeader, mdApi)
export default class UptimeDoc extends Component {
  render() {
    return (
      <div>
        <p {...propsP} ><Uptime since={+d1} interval={1000} /></p>
        <p {...propsP} ><Uptime since={+d2} /></p>
        <p {...propsP} ><Uptime since={d3.toString()} /></p>
        <p {...propsP} ><Uptime since={d4} /></p>
        <p {...propsP} ><Uptime since={d5} /></p>
      </div>
    );
  }
}
