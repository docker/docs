'use strict';

import React from 'react';
import css from 'react-css-modules';
import styles from './button.css';

const UnstyledButton = ({ onClick, children, variant = 'primary', disabled = false, type, id }) => (
  <button
    id={ id }
    type={ type }
    disabled={ disabled }
    onClick={ onClick }
    styleName={ `button ${variant}` }>
      { children }
  </button>
);

const Button = css(UnstyledButton, styles, { allowMultiple: true });

export default Button;
export {
  UnstyledButton
};
