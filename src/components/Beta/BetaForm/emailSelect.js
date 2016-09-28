import React, { PropTypes } from 'react';
import {
  Select,
  Input,
} from 'common';

import css from './styles.css';

const Emails = ({
  accountEmails,
  fields,
  onSelectChange,
  initialized = {},
}) => {
  let emails;
  if (accountEmails.length <= 1) {
    emails = (
      <Input
        readOnly
        className={css.input}
        value={accountEmails[0]}
        id={'email'}
        inputStyle={{ color: 'white', width: '100%' }}
        underlineFocusStyle={{ borderColor: 'white' }}
        placeholder="Email"
        style={{ marginBottom: '14px', width: '' }}
      />
    );
  } else {
    emails = (
      <Select
        {...fields.email}
        onBlur={() => {}}
        className={css.select}
        placeholder="Email"
        disabled={!!initialized.email}
        style={{ marginBottom: '10px', width: '' }}
        options={accountEmails.map((e) => {
          return { label: e, value: e };
        })}
        onChange={onSelectChange}
        ignoreCase
        clearable={false}
      />
    );
  }
  return emails;
};

Emails.propTypes = {
  accountEmails: PropTypes.array.isRequired,
  fields: PropTypes.any.isRequired,
  initialized: PropTypes.object,
  onSelectChange: PropTypes.func.isRequired,
};

export default Emails;
