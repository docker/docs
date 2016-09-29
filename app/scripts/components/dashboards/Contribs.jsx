'use strict';

import React, { PropTypes, Component, createClass } from 'react';
const { array } = PropTypes;
import { Link } from 'react-router';
import connectToStores from 'fluxible-addons-react/connectToStores';

import DashboardContribsStore from 'stores/DashboardContribsStore';
import RepositoriesList from 'common/RepositoriesList';
import { BlankSlates, BlankSlate } from 'common/BlankSlate';
import Pagination from 'common/Pagination';

var debug = require('debug')('DashboardContribs');

var ReposContrib = createClass({
  displayName: 'ReposContrib',
  propTypes: {
    contribs: array
  },
  _onChangePage(pageNumber) {
    debug('_onChangePage', pageNumber);
    this.props.history.pushState(null, '/contributed/', {page: pageNumber});
  },
  render() {

    var currentPageNumber = parseInt(this.props.location.query.page, 10);
    var maybePagination;
    if (this.props.contribs && this.props.contribs.length > 0) {
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
            <RepositoriesList repos={this.props.contribs} />
            {maybePagination}
          </div>
        </div>
      );

    } else {

      return (
        <BlankSlates title="No contributed repositories yet"
                     subtext="Repositories that you are a collaborator of, will show up here.">
          <BlankSlate link='/' icon='fa-home' title='Dashboard' />
        </BlankSlates>
      );

    }
  }
});

export default connectToStores(ReposContrib,
                               [
                                 DashboardContribsStore
                               ],
                               function({ getStore }, props) {
                                 return getStore(DashboardContribsStore)
                                   .getState();
                               });
