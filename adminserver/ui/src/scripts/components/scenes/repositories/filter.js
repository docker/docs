'use strict';

import React, { PropTypes, Component } from 'react';
import styles from './filter.css';
import css from 'react-css-modules';

import Typeahead from 'react-typeahead-component';

// This represents a namespace shown within the typeahead component as a search
// result.
//
// Typeahead passes data, index, userInputValue, inputValue and isSelected here.
@css(styles)
class FilterItem extends Component {
  static propTypes = {
    data: PropTypes.any
  }

  renderHeader(optionData) {
    if (optionData.index === 0 && optionData.type) {
      return (
        <div styleName='type'>{ optionData.type }</div>
      );
    }

    return null;
  }

  render() {
    const { data: optionData } = this.props;

    return (
      <div>
        { this.renderHeader(optionData) }
        <div styleName='result'>{ optionData.value }</div>
      </div>
    );
  }
}

const Filter = ({ onChange, onSelect, orgs = [], username, searchTerm = '' }) => {
  const options = [
    {
      index: 0,
      value: 'All accounts'
    }, {
      index: 0,
      type: 'User',
      value: username
    }
  ].concat(orgs.map((org, index) => ({
    index,
    type: 'Organizations',
    value: org
  })));

  return (
    <div>
      Filter by
      <div styleName='typeaheadSelectContainer'>
        <Typeahead
          options={ options }
          optionTemplate={ FilterItem }
          onChange={ onChange }
          onOptionClick={ onSelect }
          inputValue={ searchTerm }
          placeholder='All accounts' />
      </div>
    </div>
  );
};

export default css(Filter, styles);
