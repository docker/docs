'use strict';
/**
TODO: UNUSED COMPONENT. SHOULD REMOVE
*/
import React from 'react';
import ResultItem from './SearchResultItem.jsx';
var debug = require('debug')('COMPONENT:SearchResults');

var SearchResults = React.createClass({
  contextTypes: {
    getStore: React.PropTypes.func.isRequired
  },
  getResultItem: function(resultItem, idx) {
    if (resultItem.short_description) {
      if (resultItem.short_description.length > 100) {
        resultItem.short_description = resultItem.short_description.substring(0, 100) + '...';
      } else if (resultItem.short_description.length === 0) {
        resultItem.short_description = 'No description set';
      }
    }
    return <ResultItem key={idx} resultItem={resultItem} onClick={this._handleSearchResultClick.bind(null, resultItem.repo_name)}/>;
  },
  _handleSearchResultClick: function(repoName, e) {
    e.preventDefault();
    //TODO: Handle official images/repos differently
    this.props.history.pushState(null, `/r/${repoName}/`);
  },
  render: function() {
    var results = this.props.results;
    var resultItems = [];
    if (results) {
      resultItems = results.map(this.getResultItem);
    }
    return (
      <div>
        <ul className="search-results-list">{resultItems}</ul>
      </div>
    );
  }
});

module.exports = SearchResults;
