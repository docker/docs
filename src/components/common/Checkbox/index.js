import React, { Component, PropTypes } from 'react';
import CheckboxMui from 'material-ui/Checkbox';

const { string, bool } = PropTypes;

export default class Checkbox extends Component {
  static propTypes = {
    className: string,
    checked: bool,
  }

  render() {
    const { className = '', checked } = this.props;
    let fill = '#8f9ea8';
    if (checked) {
      fill = '#1aaaf8';
    }
    // TODO: Remove hardcoded color from JS
    // We need to find a way to make CSS variables
    // accesible from CSS and JS
    const iconStyle = { fill };

    return (
      <CheckboxMui
        iconStyle={iconStyle}
        {...this.props}
        className={className}
      />
    );
  }
}
