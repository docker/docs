'use strict';
import React, { PropTypes } from 'react';
import FluxibleMixin from 'fluxible-addons-react/FluxibleMixin';

import Pagination from '../common/Pagination.jsx';
import ResultsNotFound from './ResultsNotFound.jsx';
import SearchStore from './../../stores/SearchStore';
import SearchBar from './SearchBar.jsx';
import Spinner from '../Spinner.jsx';
import FA from '../common/FontAwesome';
import RepositoriesList from '../common/RepositoriesList';
import { PageHeader } from 'dux';
import _ from 'lodash';
import findKey from 'lodash/object/findKey';

var debug = require('debug')('COMPONENT:Search');

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
    pullCount: state.pullCount || 0,
    starCount: state.starCount || 0
  };
  return queryParams;
};

var Search = React.createClass({
  mixins: [FluxibleMixin],
  statics: {
    storeListeners: [SearchStore]
  },
  contextTypes: {
    getStore: React.PropTypes.func.isRequired
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
  _getCurrentFilter: function() {
    const query = this.props.location.query;
    return findKey(query, (val, key) => {
      return val === '1' && key !== 'q' && key !== 'page';
    });
  },
  _retransitionToSearch: function() {
    this.props.history.pushState(null, '/search/', _getQueryParams(this.state));
  },
  _onFilterChange: function(event) {
    event.preventDefault();
    if (event.target.value === 'isAutomated') {
      this.setState({isAutomated: 1, isOfficial: 0, starCount: 0, pullCount: 0}, this._retransitionToSearch);
    } else if (event.target.value === 'isOfficial') {
      this.setState({isOfficial: 1, isAutomated: 0, starCount: 0, pullCount: 0}, this._retransitionToSearch);
    } if (event.target.value === 'starCount') {
      this.setState({starCount: 1, isOfficial: 0, isAutomated: 0, pullCount: 0}, this._retransitionToSearch);
    } else if (event.target.value === 'pullCount') {
      this.setState({pullCount: 1, isOfficial: 0, starCount: 0, isAutomated: 0}, this._retransitionToSearch);
    } else if (event.target.value === 'all') {
      this.setState({pullCount: 0, isOfficial: 0, starCount: 0, isAutomated: 0}, this._retransitionToSearch);
    }
  },
  _onChangePage(pageNumber) {
    pageNumber = parseInt(pageNumber, 10);
    this.setState({
      page: pageNumber
    }, function() {
      this.props.history.pushState(null, '/search/', _getQueryParams(this.state));
    });
  },
  _renderMessage() {
    var message;
    if (this.state.count === 0) {
      if (this.state.isAutomated === '0' && this.state.isOfficial === '0') {
        message = _.isEmpty(this.state.query) ? `Your search is empty!` : `Your search of '${this.state.query}' did not match any repository names or descriptions.`;
      } else {
        let filterType = this.state.isAutomated === '1' ? 'automated builds' : 'official repositories';
        message = `There are no ${filterType} matching '${this.state.query}'. Try removing this filter to see more results.`;
      }
      return (
        <div className="row">
          <div className="large-8 large-centered columns">
            <ResultsNotFound heading={ <FA icon='fa-search' /> } message={ message } />
          </div>
        </div>
      );
    } else {
      return <span></span>;
    }
  },
  _renderFilterBar() {
    if (this.state.count === 0 && this.state.isAutomated === '0' && this.state.isOfficial === '0') {
      return <span></span>;
    } else {
      return (
        <div className="row">
          <div className="large-4 columns large-offset-8">
            <select defaultValue="all" value={this._getCurrentFilter()} onChange={this._onFilterChange}>
              <option value="all">All</option>
              <option value="isOfficial">Official</option>
              <option value="isAutomated">Automated</option>
              <option value="pullCount">Downloads</option>
              <option value="starCount">Stars</option>
            </select>
          </div>
        </div>
      );
    }
  },
  render: function() {
    var maybePagination;
    if (this.state.results && this.state.results.length > 0 && this.state.count > 10) {
      maybePagination = (
        <div className='row'>
          <div className='large-12 columns'>
            <Pagination next={this.state.next} prev={this.state.prev}
                        onChangePage={this._onChangePage}
                        currentPage={parseInt(this.state.page, 10) || 1}
                        pageSize={10}/>
          </div>
        </div>
      );
    }
    return (
      <div className="search-page">
        <PageHeader title={`Repositories (${this.state.count || 0 })`} />
        <br />
        <div className="search-results-container">
          { this._renderFilterBar() }
          <RepositoriesList repos={this.state.results || []} blankSlate={ this._renderMessage() }/>
          {maybePagination}
        </div>
      </div>
    );
  }
});

module.exports = Search;
