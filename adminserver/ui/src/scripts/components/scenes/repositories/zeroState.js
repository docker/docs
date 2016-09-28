'use strict';

import React from 'react';
import styles from './repositories.css';

const ZeroState = () => (
  <div className={ styles.zero }>
    <img src='/public/img/whale.png' alt='No Repositories' />
    <p>
      You have <b>no repositories</b>...<br />
      make one now!
    </p>
    <div className={ styles.arrow }></div>
  </div>
);

export default ZeroState;
