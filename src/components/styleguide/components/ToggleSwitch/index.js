import React, { Component } from 'react';
import { ToggleSwitch as Toggle } from 'common';
import asExample from '../../asExample';
import mdHeader from './header.md';
import mdApi from './api.md';

@asExample(mdHeader, mdApi)
export default class ToggleExampleSimple extends Component {
  render() {
    const styles = {
      block: {
        maxWidth: 250,
      },
      toggle: {
        marginBottom: 16,
      },
    };

    return (
      <div style={styles.block}>
        <Toggle
          label="Simple"
          style={styles.toggle}
        />
        <Toggle
          label="Toggled by default"
          defaultToggled
          style={styles.toggle}
        />
        <Toggle
          label="Disabled"
          disabled
          style={styles.toggle}
        />
        <Toggle
          label="Label on the right"
          labelPosition="right"
          style={styles.toggle}
        />
      </div>
    );
  }
}
