'use strict';
import React from 'react';
import FluxibleMixin from 'fluxible-addons-react/FluxibleMixin';
import SearchStore from '../../stores/SearchStore';
import styles from './SearchBar.css';
import FA from '../common/FontAwesome';

var debug = require('debug')('COMPONENT:SearchBar');

var _getQueryParams = function(state) {
  //transition to will always have `q` appended as query param at the very least
  //Other query params like: `s` -> sort by | `t=User` -> user | `t=Organization` -> Org | `f=official`
  // `f=automated_builds` | `s=date_created`, `s=last_updated`, `s=alphabetical`, `s=stars`, `s=downloads`
  // `s=pushes`
  var queryParams = {
    q: state.query || '',
    page: state.page || 1,
    isAutomated: state.isAutomated || 0,
    isOfficial: state.isOfficial || 0,
    starCount: state.starCount || 0,
    pullCount: state.pullCount || 0
  };
  return queryParams;
};

var SearchBar = React.createClass({
  mixins: [FluxibleMixin],
  statics: {
    storeListeners: [SearchStore]
  },
  contextTypes: {
    getStore: React.PropTypes.func.isRequired
  },
  getDefaultProps() {
    return {
      placeholder: 'Search'
    };
  },
  getInitialState: function() {
    return this.context.getStore(SearchStore).getState();
  },
  //on Search Store Change
  onChange: function() {
    //When a search query has been submitted
    var state = this.context.getStore(SearchStore).getState();
    this.setState(state);
  },
  _handleQueryChange: function(event) {
    event.preventDefault();
    //Change page to number 1 when the query is changed
    this.setState({
      page: 1
    });
    this.setState({query: event.target.value});
  },
  _handleQuerySubmit: function(event) {
    event.preventDefault();
      //second parameter will be empty object always since we don't have /search/{?}/
    //third param will be the query /search/?q=whatever&s=blah&f=bleh
    this.props.history.pushState(null, '/search/', _getQueryParams(this.state));
  },
  render: function() {
    var searchQuery = this.state.query;
    var inputPlaceholder = this.props.placeholder;
    return (
      <div className="row">
        <form className="large-12 columns" onSubmit={this._handleQuerySubmit}>
          <div className="searchbar">
            <input type="text"
                   placeholder={inputPlaceholder}
                   className={styles.searchInput}
                   onChange={this._handleQueryChange}
                   value={searchQuery} />
              <div className={styles.fa}>
                <FA icon='fa-search' />
              </div>
          </div>
        </form>
      </div>
    );
  }
});

module.exports = SearchBar;
