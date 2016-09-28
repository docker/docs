import React, { Component } from 'react';
import { Avatar } from 'common';
import asExample from '../../asExample';
import mdHeader from './header.md';
import mdApi from './api.md';

@asExample(mdHeader, mdApi)
export default class AvatarDoc extends Component {
  render() {
    return (
      <div className="clearfix">
        <h3>With image</h3>
        <Avatar src="https://www.docker.com/sites/default/files/Engine.png" />
        <br />
        <h3>With image | bigger size</h3>
        <Avatar
          src="https://www.docker.com/sites/default/files/Engine.png"
          size={128}
        />
      </div>
    );
  }
}
