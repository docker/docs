'use strict';

import React from 'react';
import css from 'react-css-modules';
import styles from './dropdown.css';

const UnstyledDropdown = ({ children, disabled = false, values }) => (
  <div
    styleName={ `wrapper ${ disabled ? 'disabled' : '' }` }
  >
    <select
      { ...values }
      styleName='default'>
        { children }
    </select>
  </div>
);

const Dropdown = css(UnstyledDropdown, styles, { allowMultiple: true });

export default Dropdown;
export {
  UnstyledDropdown
};
