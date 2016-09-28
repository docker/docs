'use strict';

import React from 'react';
import css from 'react-css-modules';
import styles from './noUsers.css';

const UnstyledNoUsers = ({ noun = '' }) => (
  <p styleName='no-users'>
    { noun === ''
      ? 'There are no users to show'
      : `This ${noun} has no members`
    }
  </p>
);

const NoUsers = css(UnstyledNoUsers, styles);

export default NoUsers;
export {
  UnstyledNoUsers
};
