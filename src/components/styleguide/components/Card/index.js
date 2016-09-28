import React, { Component } from 'react';
import { Card } from 'common';
import asExample from '../../asExample';
import mdHeader from './header.md';
import mdApi from './api.md';

@asExample(mdHeader, mdApi)
export default class CardDoc extends Component {
  render() {
    return (
      <div className="clearfix">
        <h3>With title</h3>
        <Card title="Muse">
          <p>Paranoia is in bloom</p>
          <p>The PR, transmission will resume</p>
          <p>They'll try to, push drugs that keep us all dumbed down</p>
          <p>And hope that, we will never see the truth around</p>
          <p>(So come on!!)</p>
        </Card>

        <h3>"Shy" card</h3>
        <Card shy title="Muse again">
          <p>They will not force us</p>
          <p>They will stop degraining us</p>
          <p>They will not control us</p>
          <p>Will be be victorious</p>
          <p>(So come on!!)</p>
        </Card>
        <div style={{ width: '49%', float: 'left' }}>
          <h3>No title</h3>
          <Card>
            <p>M-m-m-m-m-m-m-mad-mad-mad</p>
            <p>M-m-m-m-m-m-m-mad-mad-mad</p>
            <p>M-m-m-m-m-m-m-mad-mad-mad</p>
            <p>I can't get this memories out of my mind</p>
            <p>It's some kind of madness, is started to evolve</p>
          </Card>
        </div>
        <div style={{ width: '49%', float: 'right' }}>
          <h3>No title</h3>
          <Card>
            <p>M-m-m-m-m-m-m-mad-mad-mad</p>
            <p>M-m-m-m-m-m-m-mad-mad-mad</p>
            <p>M-m-m-m-m-m-m-mad-mad-mad</p>
            <p>I can't get this memories out of my mind</p>
            <p>It's some kind of madness, is started to evolve</p>
          </Card>

        </div>
      </div>
    );
  }
}
