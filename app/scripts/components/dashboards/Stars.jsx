'use strict';

import React, { PropTypes, Component, createClass } from 'react';
const { array, func, object } = PropTypes;
import { Link } from 'react-router';
import connectToStores from 'fluxible-addons-react/connectToStores';

import DashboardStarsStore from 'stores/DashboardStarsStore';
import RepositoriesList from 'common/RepositoriesList';
import { BlankSlates, BlankSlate } from 'common/BlankSlate';
import Pagination from 'common/Pagination';

var debug = require('debug')('DashboardStars');

var ReposStarred = createClass({
  displayName: 'ReposStarred',
  propTypes: {
    starred: array,
    history: object.isRequired
  },
  _onChangePage(pageNumber) {
    debug('_onChangePage', pageNumber);
    this.props.history.pushState(null, '/stars/', {page: pageNumber});
  },
  render() {

    var currentPageNumber = parseInt(this.props.location.query.page, 10);
    var maybePagination;
    if (this.props.starred && this.props.starred.length > 0) {
      maybePagination = (
        <div className='row'>
          <div className='large-12 columns'>
            <Pagination next={this.props.next} prev={this.props.prev}
                        onChangePage={this._onChangePage}
                        currentPage={currentPageNumber || 1}
                        pageSize={10} />
          </div>
        </div>
      );
    }

    if(this.props.starred && this.props.starred.length > 0) {

      return (
        <div className='row'>
          <div className='large-12 columns'>
            <RepositoriesList repos={this.props.starred} />
            {maybePagination}
          </div>
        </div>
      );

    } else {

      return (
        <BlankSlates title="You haven't starred anything yet"
                     subtext="Why don't you explore a bit?">
          <BlankSlate link='/explore/' icon='fa-compass' title='Go Exploring' />
        </BlankSlates>
      );

    }
  }
});

export default connectToStores(ReposStarred,
                               [
                                 DashboardStarsStore
                               ],
                               function({ getStore }, props) {
                                 return getStore(DashboardStarsStore)
                                              .getState();
                               });
