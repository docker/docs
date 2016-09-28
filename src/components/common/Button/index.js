import React, { Component, PropTypes } from 'react';
import css from './styles.css';
import { values } from 'lodash';
import { variants } from 'lib/constants';

function getStylesFromProps({
  icon = false,
  text = false,
  variant = 'primary',
  outlined = false,
  inverted = false,
  fullWidth = false,
}) {
  const styleNames = ['button', variant];
  const end = () => styleNames.map(name => css[name]).join(' ');

  if (inverted) styleNames.push('inverted');

  if (icon === true) {
    styleNames.push('icon');
    return end();
  } else if (icon === 'left') {
    styleNames.push('iconLeft');
  } else if (icon === 'right') {
    styleNames.push('iconRight');
  }

  if (text) {
    styleNames.push('text');
    return end();
  }

  if (outlined) styleNames.push('outlined');

  if (fullWidth) styleNames.push('fullWidth');

  return end();
}

const validationError = [
  'Button components can have only one of',
  '"outlined", "text" and "icon" props.',
].join(' ');

const iconError = '"icon" can only be {true}, "left" or "right".';

export default class Button extends Component {
  static propTypes = {
    variant: PropTypes.oneOf(values(variants)),
    children: PropTypes.any,
    disabled: PropTypes.bool,
    outlined: function validateOutlined({ text, icon, outlined }) {
      return outlined === true && (text === true || icon === true) ?
        new Error(validationError) : null;
    },

    text: function validateLink({ outlined, icon, text }) {
      return text === true && (outlined === true || icon === true) ?
        new Error(validationError) : null;
    },

    icon: function validateIcon({ outlined, text, icon }) {
      if (icon === 'left' || icon === 'right') return null;
      else if (icon && icon !== true) return new Error(iconError);
      return icon === true && (outlined === true || text === true) ?
        new Error(validationError) : null;
    },

    inverted: PropTypes.bool,
    className: PropTypes.string,
    element: PropTypes.oneOfType([PropTypes.element, PropTypes.string]),
    split: PropTypes.bool,
    fullWidth: PropTypes.bool,
  }

  render() {
    const { className = '', element, split, children } = this.props;
    const propStyles = getStylesFromProps(this.props);
    const props = {
      ...this.props,
      className: `dbutton ${className} ${propStyles}`,
    };

    let button;

    if (!element) {
      if (split) {
        button = (
          <button {...props}>{children}</button>
        );
      } else {
        button = (
          <button {...props}>{children}</button>
        );
      }
      return button;
    }

    return typeof element === 'string' ?
      React.createElement(element, props, children) :
      React.cloneElement(element, props, children);
  }
}
