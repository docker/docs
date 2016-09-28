'use strict';

import React, { PropTypes, Component } from 'react';
import { Link } from 'react-router';
import connectToStores from 'fluxible-addons-react/connectToStores';
import UserProfileReposStore from '../../stores/UserProfileReposStore';
import RepositoriesList from '../common/RepositoriesList';
import Pagination from '../common/Pagination';
import moment from 'moment';
import { Module } from 'dux';

var debug = require('debug')('UserProfileRepos');

class Repos extends Component {

  static propTypes = {
    user: PropTypes.object.isRequired,
    repos: PropTypes.array,
    next: PropTypes.string,
    prev: PropTypes.string
  }

  _onChangePage = (pageNumber) => {
    const username = this.props.params.user;
    this.props.history.pushState(null, `/u/${username}/`, {page: pageNumber});
  }

  render() {

    debug(this.props);
    if(this.props.repos && this.props.repos.length > 0) {

      return (
        <div className='row'>
          <div className='large-12 columns'>
            <RepositoriesList repos={this.props.repos} />
            {this.renderPagination()}
          </div>
        </div>
      );

    } else {
      debug(this.props);
      const { username, orgname } = this.props.user;
      const namespace = username || orgname;
      return (
        <Module>
        <Link to={`/u/${namespace}/`}>This user has not created any repos yet</Link>
        </Module>
      );
    }
  }
  renderPagination = (e) => {
    if(this.props.repos && this.props.repos.length > 0) {
      const currentPageNumber = parseInt(this.props.location.query.page, 10);
      return (
        <div className='row'>
          <div className='large-12 columns'>
            <Pagination next={this.props.next} prev={this.props.prev}
                        onChangePage={this._onChangePage}
                        currentPage={currentPageNumber || 1}
                        pageSize={10} />
          </div>
        </div>
      );
    } else {
      return null;
    }
  }
}

export default connectToStores(Repos,
                               [
                                 UserProfileReposStore
                               ],
                               function({ getStore }, props) {
                                 return getStore(UserProfileReposStore)
                                              .getState();
                               });
