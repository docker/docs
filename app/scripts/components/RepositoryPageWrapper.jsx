'use strict';

import styles from './RepositoryPageWrapper.css';
import React, { PropTypes, Component, cloneElement } from 'react';
const { string, number, bool, shape, array, object, func } = PropTypes;
import { Link } from 'react-router';
import _ from 'lodash';
import includes from 'lodash/collection/includes';
import merge from 'lodash/object/merge';
import omit from 'lodash/object/omit';
import words from 'lodash/string/words';
import moment from 'moment';
import FA from 'common/FontAwesome';
const REPOSTATUS = require('../stores/repostore/Constants').STATUS;

import RouteNotFound404Page from 'common/RouteNotFound404Page';
import toggleStarred from '../actions/toggleStarred';
import RepositoryPageStore from '../stores/RepositoryPageStore';
import DashboardNamespacesStore from '../stores/DashboardNamespacesStore';
import connectToStores from 'fluxible-addons-react/connectToStores';
import DocumentTitle from 'react-document-title';

var debug = require('debug')('RepositoryPageWrapper');

class Privacy extends Component {
  static propTypes = {
    isOfficial: bool.isRequired,
    isPrivate: bool.isRequired,
    isAutomated: bool.isRequired
  }
  render() {
    const publicOrPrivate = this.props.isPrivate ? 'Private' : 'Public';
    if (this.props.isOfficial) {
      return (
        <div className={styles.repoLabel}>
          <div className={styles.privacy}>Official Repository</div>
        </div>
      );
    } else if (this.props.isAutomated) {
      return (
        <div className={styles.repoLabel}>
          <div className={styles.privacy}>{publicOrPrivate} | Automated Build</div>
        </div>
      );
    } else {
      return (
        <div className={styles.repoLabel}>
          <div className={styles.privacy}>{publicOrPrivate} Repository</div>
        </div>
      );
    }
  }
}

class StarRepo extends Component {
  static propTypes = {
    hasStarred: bool.isRequired,
    toggleStar: func.isRequired
  }
  render() {
    const { toggleStar, hasStarred } = this.props;

    let toggle = toggleStar(true);
    let icon = 'fa-star-o';

    if (hasStarred) {
      toggle = toggleStar(false);
      icon = 'fa-star';
    }

    return (
      <span onClick={toggle}
            className={styles.repoStar}>
        <FA icon={icon} />
      </span>
    );
  }
}

class RepositoryPageWrapper extends Component {
  static propTypes = {
    description: string.isRequired,
    fullDescription: string.isRequired,
    hasStarred: bool.isRequired,
    isPrivate: bool.isRequired,
    isAutomated: bool.isRequired,
    lastUpdated: string,
    name: string.isRequired,
    namespace: string.isRequired,
    status: number.isRequired,
    canEdit: bool.isRequired,
    comments: shape({
      count: number.isRequired,
      results: array.isRequired
    }),
    user: object,
    JWT: string,
    history: object.isRequired,
    ownedNamespaces: array.isRequired,
    namespaces: array,
    currentUserContext: string
  }

  static contextTypes = {
    executeAction: func.isRequired
  }

  _toggleStar = status => e => {
    e.preventDefault();
    var repoShortName = this.props.namespace + '/' + this.props.name;
    this.context.executeAction(toggleStarred, {
      jwt: this.props.JWT,
      repoShortName: repoShortName,
      status: status
    });
  };

  render() {
    if (this.props.STATUS === REPOSTATUS.REPO_NOT_FOUND) {
      return <RouteNotFound404Page />;
    }

    let privacy;

    //Redirect to user profile page if public user, if org (and owner) redirect to dashboard, if public org, org profile page
    let namespaceLink;
    var {namespace, ownedNamespaces, user} = this.props;
    if (includes(ownedNamespaces, this.props.namespace) && namespace !== user.username) {
      namespaceLink = <Link to={`/u/${this.props.namespace}/dashboard/`}>{this.props.namespace}</Link>;
    } else {
      namespaceLink = <Link to={`/u/${this.props.namespace}/`}>{this.props.namespace}</Link>;
    }
    const starRepo = <StarRepo hasStarred={this.props.hasStarred}
                               toggleStar={this._toggleStar} />;
    let repoShortName = (
      <div>
        {namespaceLink}
        <span>/</span>
        <Link to={`/r/${this.props.namespace}/${this.props.name}/`}>{this.props.name}</Link>
        {starRepo}
      </div>
    );
    if (this.props.isOfficial || this.props.namespace === 'library') {
      repoShortName = (
        <div>
          <Link to={`/r/_/${this.props.name}/`}>{this.props.name}</Link>
          {starRepo}
        </div>
      );
    }

    let computedLastUpdated = 'never';
    if (this.props.lastUpdated) {
      computedLastUpdated = moment(this.props.lastUpdated).fromNow();
    }

    return (
      <DocumentTitle title={`${this.props.namespace}/${this.props.name} - Docker Hub`}>
        <div className='full-width repository-page'>
          <div className={styles.repoHeader}>
            <div className='row'>
              <div className='large-12 columns'>
                <Privacy isOfficial={this.props.isOfficial || this.props.namespace === 'library'}
                         isPrivate={this.props.isPrivate}
                         isAutomated={this.props.isAutomated}/>
                <h2 className={styles.repoTitle}>{repoShortName}</h2>
                <span className={styles.repoSubtitle}>Last pushed: {computedLastUpdated}</span>
              </div>
            </div>
          </div>
          {this.props.children && cloneElement(this.props.children, omit(this.props, 'children'))}
        </div>
      </DocumentTitle>
    );
  }
}

export default connectToStores(RepositoryPageWrapper,
                               [
                                 RepositoryPageStore,
                                 DashboardNamespacesStore
                               ],
                               ({ getStore }, props) => {
                                 return merge({},
                                              getStore(DashboardNamespacesStore).getState(),
                                              getStore(RepositoryPageStore).getState());
                               });
