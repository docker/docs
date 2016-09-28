import React, { Component } from 'react';
import { LabelField, Input, Select } from 'common';
import asExample from '../../asExample';
import mdHeader from './header.md';
import mdApi from './api.md';

@asExample(mdHeader, mdApi)
export default class LabelFieldDoc extends Component {
  render() {
    return (
      <div className="clearfix">
        <h3>Label field with text input</h3>
        <LabelField label="NODE CLUSTER NAME">
          <Input
            id="label-field-doc-input"
            hintText="Node cluster name"
            fullWidth
          />
        </LabelField>
        <br />
        <hr />
        <h3>Label field with select</h3>
        <LabelField label="NODE CLUSTER NAME">
          <Select />
        </LabelField>
        <br />
        <hr />
        <h3>Label field with unstyled text area</h3>
        <LabelField label="NODE CLUSTER NAME">
          <textarea />
        </LabelField>
        <br />
        <hr />
      </div>
    );
  }
}
