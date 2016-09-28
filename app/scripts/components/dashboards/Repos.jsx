'use strict';

import React, { createClass, PropTypes } from 'react';
import { Link } from 'react-router';
import connectToStores from 'fluxible-addons-react/connectToStores';
import DashboardReposStore from 'stores/DashboardReposStore';
import DashboardNamespacesStore from 'stores/DashboardNamespacesStore';
import RepositoryListItem from 'common/RepositoryListItem';
import Pagination from 'common/Pagination';
import FilterBar from '../filter/FilterBar';
import getAllReposForFiltering from '../../actions/getAllReposForFiltering';
import RepositoriesList from 'common/RepositoriesList';
import FA from 'common/FontAwesome';
import { PageHeader, ColModule } from 'dux';
import { BlankSlates, BlankSlate } from 'common/BlankSlate';
import merge from 'lodash/object/merge';
import filter from 'lodash/collection/filter';
import includes from 'lodash/collection/includes';
import isEmpty from 'lodash/lang/isEmpty';
import styles from './Repos.css';

const { array, bool, func } = PropTypes;

import { STATUS as COMMONSTATUS } from 'stores/common/Constants';
const { ATTEMPTING, SUCCESSFUL } = COMMONSTATUS;

const debug = require('debug')('DashboardRepos');

const ReposList = createClass({
  displayName: 'ReposList',
  contextTypes: {
    executeAction: func.isRequired
  },
  propTypes: {
    isOwner: bool,
    repos: array.isRequired
  },
  getInitialState() {
    return {
      reposList: this.props.repos,
      currentQuery: '',
      paginationMode: true,
      orgName: ''
    };
  },
  _onChangePage(pageNumber) {
    debug('_onChangePage', pageNumber);
    const { location, history, params } = this.props;
    if (location.pathname.indexOf('/dashboard') === -1) {
      history.pushState(null, '/', {page: pageNumber});
    } else {
      const { user } = params;
      history.pushState(null, `/u/${user}/dashboard/`, {page: pageNumber});
    }
  },
  _onFilterRepos(query) {
    if (query === '') {
      this.setState({
        reposList: this.props.repos
      });
    } else {
      this.setState({
        reposList: filter(this.props.repos, (repoItem) => {
          if (includes(repoItem.name, query)) {
            return repoItem;
          }
        }),
        currentQuery: query
      });
    }
  },
  _onFilterClick(e) {
    e.preventDefault();
    if (this.state.paginationMode) {
      const { user } = this.props.params;
      this.context.executeAction(getAllReposForFiltering, {
        jwt: this.props.JWT,
        user: user || this.props.user.username
      });
      this.setState({
        paginationMode: false,
        orgName: user || '',
        loadingFilterResults: true
      });
    }
  },
  _renderMessage() {
    if (!isEmpty(this.state.currentQuery)) {
      return (
        <div className='columns large-6'>
          <h6>No matching repositories for '{this.state.currentQuery}'.</h6>
        </div>
      );
    } else {
      return <span></span>;
    }
  },
  renderReposList() {
    const { paginationMode, reposList } = this.state;
    const { isOwner, repos } = this.props;

    return <RepositoriesList repos={paginationMode ? repos : reposList}
                             showPendingDelete={isOwner || true}
                             blankSlate={this._renderMessage()}/>;
  },
  componentWillReceiveProps(nextProps) {
    //We assume that loading filter results will be true only on click of the filter
    //And, the next time, we receive new props we set it to false
    //Also, we don't call the API after we get into filter mode, the user is stuck with filter mode after he clicks on filter
    //We also, would go back to pagination on click out and back
    if (nextProps.STATUS === SUCCESSFUL) {
      this.setState({
        reposList: nextProps.repos,
        currentQuery: ''
      });
    }

    //If the context changes and the dashboard goes into an org context
    //This route will have a param only when there is an Org dashboard loading -> /u/<orgname>/
    //We will set the reposList to the nextProps in this scenario and also clear the query
    const { params } = this.props;
    const orgnameParam = params.user;
    if (orgnameParam && orgnameParam !== this.state.orgName) {
      this.setState({
        orgName: orgnameParam,
        reposList: nextProps.repos,
        paginationMode: true,
        currentQuery: ''
      });
    }
  },
  render() {

    const {
      next,
      prev,
      repos,
      STATUS,
      params,
      location,
      currentUserContext
    } = this.props;

    let namespace = params.user;
    var content = (
      <BlankSlates title='Welcome to Docker Hub' subtext='Here are a few things to get you started.'>
        <BlankSlate link='/add/repository/' icon='fa-book' title='Create Repository' query={{ namespace }} />
        <BlankSlate link='/organizations/add/' icon='fa-users' title='Create Organization' />
        <BlankSlate link='/explore/' icon='fa-compass' title='Explore Repositories' />
      </BlankSlates>
    );
    var currentPageNumber = parseInt(location.query ? location.query.page : 1, 10);
    var maybePagination;
    if (repos && repos.length > 0 && this.state.paginationMode) {
      maybePagination = (
        <div className='row'>
          <div className='large-12 columns'>
            <Pagination next={next} prev={prev}
                        onChangePage={this._onChangePage}
                        currentPage={currentPageNumber || 1}
                        pageSize={10} />
          </div>
        </div>
      );
    }

    //A filter bar on top of the repositories list, when we are doing client side filtering
    let maybeMessage;
    var showTotal = 0;
    var showCount = 0;
    if (repos) {
      showTotal = repos.length;
    }
    if (this.state.reposList) {
      showCount = this.state.reposList.length;
    }
    if (STATUS !== ATTEMPTING && !this.state.paginationMode) {
      maybeMessage = <span className={styles.repoCount}>{`Showing ${showCount} of ${showTotal}`}</span>;
    } else if (STATUS === ATTEMPTING) {
      maybeMessage = (
        <span className={styles.repoCount}>
          <FA icon='fa-spin fa-spinner' /> Loading ...
        </span>
      );
    }

    let maybeFilter = (
      <div className='row'>
        <div className='large-6 columns'>
          <FilterBar onFilter={this._onFilterRepos}
                     onClick={this._onFilterClick}
                     placeholder='Type to filter repositories by name' />
        </div>
        <div className='large-4 columns end'>
          {maybeMessage}
        </div>
      </div>
    );

    if(repos && repos.length > 0) {
      const tryNautilusLink = this.props.user.orgname ?
        `/u/${this.props.user.orgname}/dashboard/billing/` :
        '/account/billing-plans/';
      let createRepoLink = '/add/repository/';
      if (currentUserContext) {
        createRepoLink = '/add/repository/?namespace=' + currentUserContext;
      }
      content = (
        <div>
          <PageHeader title='Repositories'>
            <Link to={createRepoLink} className='button'>Create Repository <FA icon='fa-plus' /></Link>
          </PageHeader>
          <div className='row'>
            <div className='large-9 columns'>
              <br />
              {maybeFilter}
              {this.renderReposList()}
              {maybePagination}
            </div>
            <div className='large-3 columns text-center'>
              <ColModule>
                <div className={styles.nautilusAdShieldIcon}>
                  <svg width="30px" height="30px" viewBox="1080 274 30 30" version="1.1">
                      <g id="ic-security-black-24-px" stroke="none" strokeWidth="1" fill="none" fill-rule="evenodd" transform="translate(1080.000000, 274.000000)">
                          <path d="M14.5,1.20833333 L3.625,6.04166667 L3.625,13.2916667 C3.625,19.9979167 8.265,26.2691667 14.5,27.7916667 C20.735,26.2691667 25.375,19.9979167 25.375,13.2916667 L25.375,6.04166667 L14.5,1.20833333 L14.5,1.20833333 Z M14.5,14.4879167 L22.9583333,14.4879167 C22.3179167,19.46625 18.995,23.9008333 14.5,25.2904167 L14.5,14.5 L6.04166667,14.5 L6.04166667,7.6125 L14.5,3.85458333 L14.5,14.4879167 L14.5,14.4879167 Z" id="Shape" fill="#8F9EA8"></path>
                          <polygon id="Shape" points="0 0 29 0 29 29 0 29"></polygon>
                      </g>
                  </svg>
                </div>
                <div className={styles.nautilusAdContent}>
                  <h6>Docker Security Scanning</h6>
                  <div className={styles.nautilusAdDesc}>
                    Protect your repositories from vulnerabilities.
                    <br />
                    <Link to={tryNautilusLink}>Try it free</Link>
                  </div>
                </div>
              </ColModule>
            </div>
          </div>
        </div>
      );
    }

    return content;
  }
});

export default connectToStores(ReposList,
  [
    DashboardReposStore,
    DashboardNamespacesStore
  ],
  function({ getStore }, props) {
    return merge(
      {},
      getStore(DashboardReposStore).getState(),
      {currentUserContext: getStore(DashboardNamespacesStore).getState().currentUserContext}
    );
  });
