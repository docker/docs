import React, { Component, PropTypes } from 'react';
import css from './styles.css';
import { SearchIcon } from '../Icon';
import { SMALL } from 'lib/constants/sizes';
const { func, string } = PropTypes;

export default class GlobalSearchBar extends Component {
  static propTypes = {
    placeholder: string,
    onChange: func.isRequired,
    onSubmit: func.isRequired,
    value: string,
  }

  onChange = (e) => {
    e.preventDefault();
    this.props.onChange(e.target.value);
  }

  onSubmit = (e) => {
    e.preventDefault();
    this.props.onSubmit(this.props.value);
  }

  focus = () => {
    this.refs.global_search_bar.focus();
  }

  render() {
    const { placeholder, value } = this.props;

    return (
      <div className={`dGlobalSearchBar ${css.wrapper}`}>
        <form onSubmit={this.onSubmit}>
          <SearchIcon className={css.icon} size={SMALL} />
          <input
            ref="global_search_bar"
            className={css.bar}
            onChange={this.onChange}
            placeholder={placeholder}
            value={value}
          />
        </form>
      </div>
    );
  }
}
