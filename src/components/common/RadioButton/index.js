import React, { Component } from 'react';
import RadioButtonMui from 'material-ui/RadioButton';

export default class RadioButton extends Component {
  render() {
    return (
      <RadioButtonMui {...this.props} />
    );
  }
}
