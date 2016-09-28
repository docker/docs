import React, { PropTypes } from 'react';
import { Select, Input } from 'components/common';

const EmailSelect = ({
  accountEmails,
  fields,
  onSelectChange,
  initialized,
}) => {
  let emails;
  if (accountEmails.length <= 1) {
    // If there is only one email, render a readOnly input with that email
    // selected
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
        disabled={!!(initialized && initialized.email)}
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

EmailSelect.propTypes = {
  accountEmails: PropTypes.array.isRequired,
  fields: PropTypes.any.isRequired,
  initialized: PropTypes.object,
  onSelectChange: PropTypes.func.isRequired,
};

export default EmailSelect;
