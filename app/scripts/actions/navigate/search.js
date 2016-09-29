'use strict';
var debug = require('debug')('navigate::search');

import {
  Search
} from 'hub-js-sdk';

export default function search({actionContext, payload, done, maybeData}){
  debug('Hit /search:: Query = ' + JSON.stringify(payload.location.query));
  debug('Searching for: ' + payload.location.query.q);
  //the filter param values are `0` because the python api accepts either `0` or `False`, we decided to use `0`
  var searchQueryParams = {
    query: payload.location.query.q || '',
    page: payload.location.query.page || '',
    isAutomated: payload.location.query.isAutomated || 0,
    isOfficial: payload.location.query.isOfficial || 0,
    starCount: payload.location.query.starCount || 0,
    pullCount: payload.location.query.pullCount || 0
  };

  //TODO: Maybe we should have a single dispatch here?
  actionContext.dispatch('SUBMIT_SEARCH_QUERY', searchQueryParams.query);
  actionContext.dispatch('UPDATE_SEARCH_PAGE', searchQueryParams.page);
  actionContext.dispatch('UPDATE_SEARCH_OTHERFILTERS', searchQueryParams);

  //This is to search repositories
  Search.searchRepos(maybeData.token, searchQueryParams, function(err, res) {
    if (err) {
      debug(err);
      actionContext.dispatch('SEARCH_ERROR', err);
      done();
    } else {
      debug(res);
      var queryResult = res.body;
      //Query Result details (for paging)
      if (queryResult) {
        actionContext.dispatch('PROCESS_SEARCH_RESULTS', queryResult);
        done();
      } else {
        done();
      }
    }
  });
}
