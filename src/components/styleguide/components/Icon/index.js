import React, { Component } from 'react';
import { map } from 'lodash';
import * as icons from 'common/Icon';
import asExample from '../../asExample';
import { sizes, variants } from 'lib/constants';

import mdHeader from './header.md';
import mdApi from './api.md';

@asExample(mdHeader, mdApi)
export default class IconDoc extends Component {
  render() {
    const { DockerIcon, TutumIcon } = icons;

    return (
      <div>
        <h4>Available Icons</h4>
        { map(icons, (Element, key) => <Element key={key} />) }
        <h4>Sizes</h4>
        { map(sizes, (size, key) => {
          return (
            <div key={key}>
              <h5>{`${size}`}</h5>
              <DockerIcon size={size} />
            </div>
          );
        })}
        <h4>Variants</h4>
        { map(variants, (variant, key) => {
          return (
            <div key={key}>
              <h5>{`${variant}`}</h5>
              <TutumIcon variant={variant} />
            </div>
          );
        })}
      </div>
    );
  }
}
