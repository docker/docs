import React, { Component } from 'react';
import { map } from 'lodash';
import * as components from './components';

export default class Styleguide extends Component {
  render() {
    return (
      <div id="module-styleguide">
        { map(components, (Element, key) => <Element key={key} />) }
      </div>
    );
  }
}
