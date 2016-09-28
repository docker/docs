import React, { PropTypes } from 'react';
import { Select } from '../lib/common';
import { noOp } from '../lib/helpers';
import css from './styles.css';
const { string, func, array } = PropTypes;

const Accounts = ({ options, onSelectChange, selectedNamespace }) => {
  let accountSelects = null;
  if (options.length > 1) {
    accountSelects = (
      <div className={css.account}>
        <div className={css.sectionTitle}>
          Account
        </div>
        <Select
          value={selectedNamespace}
          className={css.accountSelect}
          clearable={false}
          ignoreCase
          onBlur={noOp}
          onChange={onSelectChange}
          options={options}
          placeholder="Account"
        />
      </div>
    );
  }
  return accountSelects;
};

Accounts.propTypes = {
  selectedNamespace: string.isRequired,
  onSelectChange: func.isRequired,
  options: array.isRequired,
};

export default Accounts;
