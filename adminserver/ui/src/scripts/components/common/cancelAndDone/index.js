'use strict';

import React from 'react';
import Button from 'components/common/button';
// css
import css from 'react-css-modules';
import styles from './cancelAndDone.css';

export const UnstyledCancelAndDone = ({ onCancel, onDone }) => (
  <div styleName='buttons'>
    <div>
      <Button
        type='button'
        variant='primary outline'
        onClick={ (evt) => { evt.preventDefault(); onCancel(evt); } }>Cancel</Button>
      <Button
        variant='primary'
        type='submit'
        onClick={ onDone }>Save</Button>
    </div>
  </div>
);

export default css(UnstyledCancelAndDone, styles);
