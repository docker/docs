import React, { Component, PropTypes } from 'react';
import { Checkbox } from 'common';
import css from './styles.css';
import forEach from 'lodash/forEach';
const { arrayOf, bool, func, shape, string } = PropTypes;

export default class FilterList extends Component {
  static propTypes = {
    className: string,
    clearAllText: string,
    filters: arrayOf(shape({
      disabled: bool,
      isChecked: bool,
      label: string,
      value: string,
    })).isRequired,
    onChange: func.isRequired,
    title: string,
  }

  onChange = (val) => (e, isInputChecked) => {
    // Get all of the values of the checked filters and return them as a
    // comma separated string
    const checkedFilters = [];
    forEach(this.props.filters, ({ isChecked, value }) => {
      // This checkbox was previously checked or it was the one just checked
      if (val !== value && isChecked || val === value && isInputChecked) {
        checkedFilters.push(value);
      }
    });
    this.props.onChange(checkedFilters.join(','));
  }

  clearAll = () => {
    this.props.onChange('');
  }

  mkFilterLine = ({ disabled, isChecked, label, value }) => {
    const labelStyle = { color: disabled ? '#c0c9ce' : '#445d6e' };
    return (
      <Checkbox
        className={css.input}
        checked={isChecked}
        disabled={disabled}
        key={value}
        label={label}
        labelStyle={labelStyle}
        onCheck={this.onChange(value)}
      />
    );
  }

  render() {
    const { className = '', clearAllText, filters, title } = this.props;
    let clearAll;
    if (clearAllText) {
      clearAll = (
        <div className={css.clearAll} onClick={this.clearAll}>
          { clearAllText }
        </div>
      );
    }
    return (
      <div className={className}>
        <div className={css.title}>{ title }</div>
        { clearAll }
        { filters.map(this.mkFilterLine) }
      </div>
    );
  }
}
