'use strict';

import React, { Component, PropTypes } from 'react';
const { array, object, func } = PropTypes;
import Typeahead from 'react-typeahead-component';
import FA from 'components/common/fontAwesome';
import ui from 'redux-ui';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import * as searchActions from 'actions/search';
import throttle from 'lodash/throttle';
import { mapActions } from 'utils';
import { createStructuredSelector } from 'reselect';
import { limitedSearchResults } from './selectors';
import SearchResultItem from './searchResultItem.js';
import styles from './search.css';

const mapState = createStructuredSelector({
  searchResults: limitedSearchResults
});

@ui({
  key: 'globalSearch',
  state: {
    term: '',
    // The max number of orgs/repos etc we can list
    maxResultsPerResource: 4,
    // The max search results to show across all resources
    maxResults: 12
  },
  persist: true
})
@connect(mapState, mapActions({ ...searchActions, push }))
export default class Search extends Component {

  static propTypes = {
    ui: object.isRequired,
    updateUI: func.isRequired,
    actions: object.isRequired,

    searchResults: array
  }

  static defaultProps = {
    searchResults: []
  }

  componentWillMount() {
    this.search = throttle(::this.props.actions.searchAll, 250);
  }

  updateSearch(evt) {
    const { value } = evt.target;
    this.props.updateUI('term', value);
    if (value === '') {
      this.props.actions.clearSearch();
    } else {
      this.search(value);
    }
  }

  // Called when the enter key is pressed on a suggestion
  onSelectResult(evt, item) {
    if (evt.which === 13 && item.type) {
      evt.target.blur();
      this.navigateTo(item);
    }
  }

  onClickResult(_, item) {
    this.navigateTo(item);
  }

  navigateTo(item) {
    const { data } = item;
    let url;

    // TODO: Navigate to item
    switch(item.type) {
    case 'repo':
      url = `/repositories/${data.namespace}/${data.name}`;
      break;
    case 'org':
      url = `/orgs/${data.name}`;
      break;
    case 'user':
      url = `/users/${data.name}`;
      break;
    }
    this.props.actions.push(url);
  }

  render() {
    return (
      <div className={ styles.searchInputContainer }>
        <Typeahead
          options={ this.props.searchResults }
          filterOption='name'
          displayOption={ data => data }
          optionTemplate={ SearchResultItem }
          onChange={ ::this.updateSearch }
          inputValue={ this.props.ui.term }
          onKeyDown={ ::this.onSelectResult }
          onOptionClick={ ::this.onClickResult }
          className={ styles.typeaheadContainer }
          placeholder='Search Trusted Registry' />
        <FA icon='fa-search' className={ styles.icon } />
      </div>
    );
  }
}
