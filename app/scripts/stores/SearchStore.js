'use strict';
const createStore = require('fluxible/addons/createStore');
const _ = require('lodash');
const debug = require('debug')('STORE:SearchStore');

//Store to keep track of searches
//TODO: Autocomplete support?
//Query API on general search and return results in a Search Results Component
var SearchStore = createStore({
    storeName: 'SearchStore',
    handlers: {
      SUBMIT_SEARCH_QUERY: '_submitSearchQuery',
      UPDATE_SEARCH_FILTER: '_updateSearchFilter',
      UPDATE_SEARCH_SORT: '_updateSearchSort',
      UPDATE_SEARCH_PAGE: '_updateSearchPage',
      UPDATE_SEARCH_OTHERFILTERS: '_updateSearchFilters',
      PROCESS_SEARCH_RESULTS: '_processSearchResults',
      SEARCH_ERROR: '_handleSearchError'
    },
    initialize: function() {
      this.results = '';
      this.queryResult = {};
      this.page = 1;
      this.count = 0;
      this.next = false;
      this.prev = false;
    },
    getQueryParams: function() {
      //transition to will always have `q` appended as query param at the very least
      //Other query params like: `s` -> sort by | `t=User` -> user | `t=Organization` -> Org | `f=official`
      // `f=automated_builds` | `s=date_created`, `s=last_updated`, `s=alphabetical`, `s=stars`, `s=downloads`
      // `s=pushes`
      var queryParams = {
        q: this.query || '',
        page: this.page || 1,
        isOfficial: this.isOfficial || 0,
        isAutomated: this.isAutomated || 0,
        pullCount: this.pullCount || 0,
        starCount: this.starCount || 0
      };
      return queryParams;
    },
    _submitSearchQuery: function(payload: string) {
      this.query = payload;
      this.results = null;
      this.emitChange();
    },
    _processSearchResults: function(searchResult) {
      this.queryResult = searchResult;
      this.count = searchResult.count;
      this.results = searchResult.results;
      this.next = searchResult.next;
      this.prev = searchResult.previous;
      this.emitChange();
    },
    _handleSearchError: function(searchError) {
      //TODO: Some form of common error handling across all components
      debug(searchError);
    },
    _updateSearchPage: function(page) {
      this.page = page;
      this.emitChange();
    },
    _updateSearchFilters: function(params) {
      this.isAutomated = params.isAutomated;
      this.isOfficial = params.isOfficial;
      this.pullCount = params.pullCount;
      this.starCount = params.starCount;
      this.emitChange();
    },
    getState: function() {
      return {
        query: this.query,
        page: this.page,
        queryResult: this.queryResult,
        results: this.results,
        isOfficial: this.isOfficial,
        isAutomated: this.isAutomated,
        pullCount: this.pullCount,
        starCount: this.starCount,
        count: this.count,
        next: this.next,
        prev: this.prev
      };
    },
    dehydrate: function() {
      return this.getState();
    },
    rehydrate: function(state) {
      this.query = state.query;
      this.page = state.page;
      this.queryResult = state.queryResult;
      this.results = state.results;
      this.isOfficial = state.isOfficial;
      this.isAutomated = state.isAutomated;
      this.pullCount = state.pullCount;
      this.starCount = state.starCount;
      this.count = state.count;
      this.next = state.next;
      this.prev = state.prev;
    }
});

module.exports = SearchStore;
