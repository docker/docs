import React, { Component } from 'react';
import {
  Menu,
  Button,
  EllipsisIcon,
  ChevronIcon,
} from 'common';
import asExample from '../../asExample';
import mdHeader from './header.md';
import mdApi from './api.md';

@asExample(mdHeader, mdApi)
export default class IconMenuDoc extends Component {
  render() {
    const trigger1 = <Button icon variant="dull"><EllipsisIcon /></Button>;
    const items1 = [
      { label: 'First', value: 'first_value' },
      { label: 'Second', value: 'second_value' },
      { label: 'Third', value: 'third_value' },
    ];

    const trigger2 = (
      <Button icon="right" outlined><ChevronIcon />Offsetted</Button>
    );
    const items2 = [
      { label: 'Redeploy', value: 'redeploy_value' },
      { label: 'Scale', value: 'scale_value', disabled: true },
      { label: 'Terminate', value: 'terminate_value' },
    ];
    // eslint-disable-next-line no-alert
    const cb = (item) => alert(`Selected ${item}`);
    return (
      <div>
        <Menu
          style={{ marginRight: 20 }}
          trigger={trigger1}
          onSelect={cb}
          items={items1}
        />
        <Menu
          offset={[0, -5]}
          style={{ marginRight: 20 }}
          trigger={trigger2}
          onSelect={cb}
          items={items2}
        />
      </div>
    );
  }
}
