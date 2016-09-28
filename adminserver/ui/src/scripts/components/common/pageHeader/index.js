'use strict';

import React from 'react';
import styles from './pageHeader.css';

const PageHeader = ({ title }) => (
  <div className={ styles.pageHeader }>
    <h1>{ title }</h1>
  </div>
);

export default PageHeader;
