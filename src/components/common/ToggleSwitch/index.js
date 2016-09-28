import React, { Component } from 'react';
import SwitchMui from 'material-ui/Toggle';

// Ref: http://www.material-ui.com/#/components/toggle
export default class Switch extends Component {
  render() {
    return (
      <SwitchMui {...this.props} />
    );
  }
}
