import React, { Component, PropTypes } from 'react';
import { RadioButtonGroup as RadioButtonGroupMui }
  from 'material-ui/RadioButton';

const { any, string } = PropTypes;

export default class RadioButtonGroup extends Component {

  static propTypes = {
    name: string.isRequired,
    defaultSelected: string,
    children: any,
  }

  render() {
    // eslint-disable-next-line no-use-before-define
    const { children, ...props } = this.props;
    return (
      <RadioButtonGroupMui {...props}>
        {children}
      </RadioButtonGroupMui>
    );
  }
}
