import React, { Component } from 'react';
import { RadioButtonGroup, RadioButton } from 'common';
import asExample from '../../asExample';
import mdHeader from './header.md';
import mdApi from './api.md';

@asExample(mdHeader, mdApi)
export default class RadioButtonDoc extends Component {
  render() {
    const styles = {
      radioButton: {
        marginBottom: 16,
      },
    };

    return (
      <div className="clearfix">
        <h3>Radio Buttons</h3>
        <br />
        <RadioButtonGroup name="dockerWhales" defaultSelected="mobyDock">
            <RadioButton
              value="mollyDock"
              label="Molly Dock"
              style={styles.radioButton}
            />
            <RadioButton
              value="mobyDock"
              label="Moby Dock"
              style={styles.radioButton}
            />
        </RadioButtonGroup>
        <br />
      </div>
    );
  }
}
