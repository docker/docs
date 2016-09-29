'use strict';

import React, { PropTypes, Component } from 'react';
const { array, string, func } = PropTypes;
import connectToStores from 'fluxible-addons-react/connectToStores';

import UserProfileStarsStore from 'stores/UserProfileStarsStore';
import RepositoriesList from 'common/RepositoriesList';
import Pagination from 'common/Pagination';
import { Module } from 'dux';

var debug = require('debug')('UserProfileStars');

var UserStars = React.createClass({
  displayName: 'UserStars',
  propTypes: {
    starred: array,
    next: string,
    prev: string
  },
  _onChangePage(pageNumber) {
    this.props.history.pushState(null, `/u/${this.props.user.username}/starred/`, {page: pageNumber});
  },
  render() {

    var currentPageNumber = parseInt(this.props.location.query.page, 10);
    var maybePagination;

    if(this.props.starred && this.props.starred.length > 0) {
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
        <Module>This user does not have any starred repos.</Module>
      );

    }
  }
});

export default connectToStores(UserStars,
  [
    UserProfileStarsStore
  ],
  function({ getStore }, props) {
    return getStore(UserProfileStarsStore)
      .getState();
  });
