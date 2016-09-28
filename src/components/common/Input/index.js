import React, { Component, PropTypes } from 'react';
import TextField from 'material-ui/TextField';
import css from './styles.css';
import classnames from 'classnames';
const { string, bool, object, oneOf } = PropTypes;
import { TINY, SMALL, REGULAR, LARGE } from 'lib/constants/sizes';
// widths hardcoded from variables css $input-width-xsmall through -large
const widths = {
  [TINY]: '82px',
  [SMALL]: '144px',
  [REGULAR]: '176px',
  [LARGE]: '238px',
};
const sizes = Object.keys(widths);

export default class Input extends Component {
  static propTypes = {
    className: string,
    id: string.isRequired,
    size: oneOf(sizes),
    style: object,
    readOnly: bool,
  }

  focus() {
    this.refs.input_text_field.focus();
  }

  render() {
    const { className = '', size, style, readOnly } = this.props;

    // TODO: Remove hardcoded color from JS
    // We need to find a way to make CSS variables
    // accesible from CSS and JS
    const underlineStyles = {
      borderColor: '#1AAAF8',
    };

    const inputStyles = { marginBottom: '14px', width: '', ...style };
    if (size) {
      inputStyles.width = widths[size];
    }

    const classes = classnames({
      [css.main]: true,
      [className]: !!className,
    });

    return (
      <TextField
        ref="input_text_field"
        name="input_text_field"
        underlineFocusStyle={underlineStyles}
        disabled={readOnly}
        {...this.props}
        style={inputStyles}
        className={classes}
      />
    );
  }
}
