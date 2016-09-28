import React, { PropTypes } from 'react';
import { Select } from 'common';
import css from './styles.css';
const { any, func, array } = PropTypes;
const noOp = () => {};

const stopPropagation = (e) => {
  e.stopPropagation();
};

const AccountSelect = ({ options, fields, onSelectChange }) => {
  let accountSelects = null;
  if (options.length > 1) {
    accountSelects = (
      <div className={css.account} onClick={stopPropagation}>
        <div className={css.sectionTitle}>
          Attach this subscription to an account
        </div>
        <Select
          {...fields.account}
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

AccountSelect.propTypes = {
  fields: any.isRequired,
  onSelectChange: func.isRequired,
  options: array.isRequired,
};

export default AccountSelect;
