import React, { Component } from 'react';
import {
  DockerFlatIcon,
  PrivateIcon,
  PublicIcon,
  Tab,
  Tabs,
} from 'common';
import asExample from '../../asExample';
import mdHeader from './header.md';
import mdApi from './api.md';

@asExample(mdHeader, mdApi)
export default class TabsTabDoc extends Component {

  state = {
    selected1: 0,
    selected2: 'val2',
    selected3: [true, false, true, false, false],
    selected4: 'v1',
  }

  handleSelect1 = (ev, index) => {
    this.setState({ selected1: index });
  }

  handleSelect2 = (ev, value) => {
    this.setState({ selected2: value });
  }

  handleSelect3 = (index) => {
    const selected3 = [].concat(this.state.selected3);
    selected3[index] = !selected3[index];
    this.setState({ selected3 });
  }

  handleSelect4 = (ev, value) => {
    this.setState({ selected4: value });
  }

  render() {
    return (
      <div className="clearfix">
        <h3>Default</h3>
        <Tabs
          selected={this.state.selected1}
          onSelect={this.handleSelect1}
        >
          <Tab>Sao Paulo</Tab>
          <Tab>Santiago</Tab>
          <Tab>Rio</Tab>
          <Tab>Lima</Tab>
          <Tab>Buenos Aires</Tab>
        </Tabs>
        <br />
        <h3>By Value (no index) </h3>
        <Tabs
          selected={this.state.selected2}
          onSelect={this.handleSelect2}
        >
          <Tab value="v1">Madrid</Tab>
          <Tab value="v2">London</Tab>
          <Tab value="v3">Paris</Tab>
          <Tab value="v4">Copenhagen</Tab>
          <Tab value="v5">Dublin</Tab>
        </Tabs>
        <br />
        <h3>Without Tabs Controller</h3>
        <div>
          <Tab
            selected={this.state.selected3[0]}
            onClick={() => this.handleSelect3(0)}
          >
            San Francisco
          </Tab>
          <Tab
            selected={this.state.selected3[1]}
            onClick={() => this.handleSelect3(1)}
          >
            New York
          </Tab>
          <Tab
            selected={this.state.selected3[2]}
            onClick={() => this.handleSelect3(2)}
          >
            Portland
          </Tab>
          <Tab
            selected={this.state.selected3[3]}
            onClick={() => this.handleSelect3(3)}
          >
            Chicago
          </Tab>
          <Tab
            selected={this.state.selected3[4]}
            onClick={() => this.handleSelect3(4)}
          >
            Boston
          </Tab>
        </div>
        <br />
        <h3>Icon Tabs</h3>
        <Tabs
          selected={this.state.selected4}
          onSelect={this.handleSelect4}
          icons
        >
          <Tab value="v1">
            <DockerFlatIcon />
          </Tab>
          <Tab value="v2">
            <PrivateIcon />
          </Tab>
          <Tab value="v3">
            <PublicIcon />
          </Tab>
        </Tabs>
      </div>
    );
  }
}
