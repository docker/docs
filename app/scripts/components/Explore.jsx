'use strict';

import React, { createClass, PropTypes } from 'react';
import connectToStores from 'fluxible-addons-react/connectToStores';
// TODO: This is dirty; replace with ExploreStore?
import DashboardReposStore from '../stores/DashboardReposStore';
import RepositoryListItem from './common/RepositoryListItem';
import Pagination from './common/Pagination.jsx';
import { PageHeader } from 'dux';
import DocumentTitle from 'react-document-title';
const debug = require('debug')('Explore');

function mkRepoListItem(repo) {
  /**
   * TODO: after snakecase to camelcase in hub-js-sdk is done remove
   * explicit props
   */
  return (
    <RepositoryListItem {...repo}
                        isPrivate={repo.is_private}
                        isTrusted={repo.is_automated}
                        starCount={repo.star_count}
                        pullCount={repo.pull_count}
                        key={repo.name} />
    );
}

var ReposList = createClass({
  displayName: 'ExploreReposList',
  _onChangePage(pageNumber) {
    debug('_onChangePage', pageNumber);
    this.props.history.pushState(null, '/explore/', {page: pageNumber});
  },
  render() {
    var content = (<div></div>);
    var currentPageNumber = parseInt(this.props.location.query.page, 10);
    if(this.props.repos && this.props.repos.length > 0) {
      content = (
        <ul className='large-12 columns no-bullet'>
          {this.props.repos.map(mkRepoListItem)}
        </ul>
      );
    }

    return (
      <DocumentTitle title='Explore - Docker Hub'>
        <div>
          <PageHeader title='Explore Official Repositories' />
          <div className='row explore-repo-list'>
            {content}
          </div>
          <div className='row'>
            <div className='large-12 columns'>
              <Pagination next={this.props.next} prev={this.props.prev}
                onChangePage={this._onChangePage}
                currentPage={currentPageNumber || 1}
                pageSize={10} />
            </div>
          </div>
        </div>
      </DocumentTitle>
    );
  }
});

export default connectToStores(ReposList,
                               [
                                 DashboardReposStore
                               ],
                               function({ getStore }, props) {
                                 return getStore(DashboardReposStore)
                                              .getState();
                               });
