import React, { PropTypes } from 'react';
import { Select } from 'components/common';

import css from './styles.css';

const generateMonths = () => {
  const months = [];
  for (let i = 1; i < 13; i++) {
    const num = `0${i}`.slice(-2);
    months.push({ label: num, value: num });
  }
  return months;
};

const generateYears = () => {
  const years = [];
  const thisYear = new Date().getFullYear();
  for (let i = thisYear; i < thisYear + 20; i++) {
    years.push({ label: i, value: i });
  }
  return years;
};

const CardExpiration = props => {
  const { fields, onSelectChange, initialized } = props;
  const { expMonth, expYear } = fields;
  const monthClass = (expMonth.touched && expMonth.error) ?
    css.expirationErr : '';
  const yearClass = (expYear.touched && expYear.error) ?
    css.expirationErr : '';
  const error = !!monthClass || !!yearClass ?
    (<div className={css.dateErr}>Invalid Expiration</div>) : null;
  const accountSelects = (
    <div className={css.expiration}>
      <div className={css.title}>Expiration Date</div>
      <div className={css.dates}>
        <Select
          {...fields.expMonth}
          disabled={!!initialized.expMonth}
          onBlur={() => {}}
          placeholder="Month"
          onChange={onSelectChange('expMonth')}
          className={monthClass}
          options={generateMonths()}
          ignoreCase
          clearable={false}
        />
        <Select
          {...fields.expYear}
          disabled={!!initialized.expYear}
          onBlur={() => {}}
          placeholder="Year"
          onChange={onSelectChange('expYear')}
          className={yearClass}
          options={generateYears()}
          ignoreCase
          clearable={false}
        />
      </div>
      {error}
    </div>
  );
  return accountSelects;
};

CardExpiration.propTypes = {
  fields: PropTypes.any.isRequired,
  onSelectChange: PropTypes.func.isRequired,
};

export default CardExpiration;
