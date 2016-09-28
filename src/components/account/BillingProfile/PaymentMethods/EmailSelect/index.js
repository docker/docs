import React, { PropTypes } from 'react';
import {
  Select,
  Input,
} from 'common';

const Emails = ({ accountEmails, fields, onSelectChange }) => {
  let emails;
  if (accountEmails.length <= 1) {
    emails = (
      <Input
        readOnly
        value={accountEmails[0]}
        id={'email'}
        placeholder="Email"
        style={{ marginBottom: '14px', width: '' }}
      />
    );
  } else {
    emails = (
      <Select
        {...fields.email}
        onBlur={() => {}}
        placeholder="Email"
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
  onSelectChange: PropTypes.func.isRequired,
};

export default Emails;
