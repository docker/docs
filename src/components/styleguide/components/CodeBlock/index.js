import React, { Component } from 'react';
import { CodeBlock } from 'common';
import asExample from '../../asExample';
import mdHeader from './header.md';
import mdApi from './api.md';

@asExample(mdHeader, mdApi)
export default class CodeBlockDoc extends Component {
  render() {
    return (
      <div className="clearfix">
        <CodeBlock>
          $ code here
          <br />
          $ code here
        </CodeBlock>
      </div>
    );
  }
}
