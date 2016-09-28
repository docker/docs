'use strict';

import React from 'react';
import css from 'react-css-modules';
import styles from './box.css';

export const UnstyledBox = ({ children }) =>
  <div styleName='box'>{ children }</div>;

const Box = css(UnstyledBox, styles);

export default Box;
