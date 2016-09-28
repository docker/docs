import React, { Component } from 'react';
import { CircularLoading } from 'common';
import asExample from '../../asExample';
import mdHeader from './header.md';
import mdApi from './api.md';

@asExample(mdHeader, mdApi)
export default class CodeBlockDoc extends Component {
  render() {
    return (
      <div>
        tiny: <CircularLoading size="tiny" />
        small: <CircularLoading size="small" />
        regular: <CircularLoading size="regular" />
        large: <CircularLoading size="large" />
        xlarge: <CircularLoading size="xlarge" />
      </div>
    );
  }
}
