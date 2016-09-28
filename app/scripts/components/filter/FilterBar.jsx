'use strict';
import React, { createClass, PropTypes } from 'react';
import FA from 'common/FontAwesome';
import styles from './FilterBar.css';
const { bool, func, string } = PropTypes;
const debug = require('debug')('COMPONENT:FilterBar');

const FilterBar = createClass({
  //Optional placeholder string
  //Filter function to invoke when the query is submitted
  propTypes: {
    placeholder: string,
    onFilter: func.isRequired,
    onClick: func
  },
  getDefaultProps() {
    return {
      placeholder: 'Type to filter'
    };
  },
  getInitialState() {
    return {
      query: ''
    };
  },
  _clearQuery: function(event) {
    this.setState({
      query: ''
    });
    this.props.onFilter('');
  },
  _handleQueryChange: function(event) {
    event.preventDefault();
    this.props.onFilter(event.target.value);
    this.setState({query: event.target.value});
  },
  _onSubmit: function(event) {
    event.preventDefault();
  },
  render: function() {
    const { onClick, placeholder } = this.props;
    const { query } = this.state;
    const maybeCancel = query ? '╳' : '';

    return (
      <div className="row">
        <form ref='filterForm' className={'large-12 columns filterbar-container ' + styles.form} onSubmit={this._onSubmit}>
          <div className={styles.filterInput + ' row collapse'}>
            <div className="small-11 columns">
              <input type="text"
                     placeholder={placeholder}
                     onClick={onClick}
                     onChange={this._handleQueryChange}
                     value={query} />
            </div>
            <div className="small-1 columns">
              <span className={styles.iconStyle + ' postfix'} onClick={this._clearQuery}>{maybeCancel}</span>
            </div>
          </div>
        </form>
      </div>
    );
  }
});

module.exports = FilterBar;
