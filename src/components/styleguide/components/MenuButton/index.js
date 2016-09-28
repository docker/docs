import React, { Component } from 'react';
import {
  MenuButton,
} from 'common';
import asExample from '../../asExample';
import mdHeader from './header.md';
import mdApi from './api.md';

@asExample(mdHeader, mdApi)
export default class MenuButtonDoc extends Component {
  render() {
    const items1 = [
      { label: 'First', value: 'first_value' },
      { label: 'Second', value: 'second_value' },
      { label: 'Third', value: 'third_value' },
    ];
    // eslint-disable-next-line no-alert
    const cb = (item) => alert(`Selected ${item}`);
    return (
      <div>
        <MenuButton
          onClick={() => console.log('MenuButton Clicked')}
          onSelect={cb}
          items={items1}
        >Menu Button</MenuButton>
      </div>
    );
  }
}
