'use strict';

import React, { PropTypes, createClass, cloneElement } from 'react';
import connectToStores from 'fluxible-addons-react/connectToStores';
import UserProfileStore from '../stores/UserProfileStore';
import RepositoriesList from './common/RepositoriesList';
import RouteNotFound404Page from './common/RouteNotFound404Page';
import moment from 'moment';
import { SecondaryNav } from 'dux';
import FA from 'common/FontAwesome';
import { Link } from 'react-router';
import _ from 'lodash';
import styles from './UserProfileWrapper.css';
import DocumentTitle from 'react-document-title';

var debug = require('debug')('UserProfileWrapper');

var UserShape = {
  id: PropTypes.string.isRequired,
  username: PropTypes.string,
  orgname: PropTypes.string,
  full_name: PropTypes.string.isRequired,
  location: PropTypes.string.isRequired,
  company: PropTypes.string.isRequired,
  profile_url: PropTypes.string.isRequired,
  date_joined: PropTypes.string.isRequired,
  gravatar_url: PropTypes.string.isRequired
};

var ProfileCard = createClass({
  displayName: 'ProfileCard',
  propTypes: UserShape,
  render() {
    var gravatar = this.props.gravatar_url;
    /**
     * If the gravatar has a size equal to 80 in the url,
     * change it to 512. This is a HACK.
     */
    if(this.props.gravatar_url && this.props.gravatar_url.match(/s=80/)) {
      gravatar = gravatar.replace('s=80', 's=512');
    }

    var maybeLocation = null;
    var maybeCompany = null;
    var maybeProfileUrl = null;

    if (this.props.location) {
      maybeLocation = <li className={styles.item}><FA icon='fa-map-marker'/> {this.props.location}</li>;
    }
    if (this.props.maybeCompany) {
      maybeCompany = <li className={styles.item}><FA icon='fa fa-building'/> {this.props.company}</li>;
    }
    if (this.props.profile_url) {
      maybeProfileUrl = <li className={styles.item}><FA icon='fa fa-home'/> <a href={this.props.profile_url}>{this.props.profile_url}</a></li>;
    }

    return (
      <div className='row'>
        <div className='large-12 columns'>
          <div className={styles.gravatar}>
            <img src={gravatar} />
          </div>
        </div>
        <div className='large-12'>
          <h1 className={styles.heading}>{this.props.username || this.props.orgname}</h1>
          <h2 className={styles.heading}>{this.props.full_name}</h2>
          <ul className={styles.userinfo}>
            {maybeLocation}
            {maybeCompany}
            {maybeProfileUrl}
            <li className={styles.item}><FA icon='fa fa-clock-o'/> Joined {moment(this.props.date_joined).format('MMMM YYYY')}</li>
          </ul>
        </div>
      </div>
    );
  }
});

var UserProfile = createClass({
  displayName: 'UserProfile',
  propTypes: {
    user: PropTypes.shape(UserShape)
  },
  render() {
    if(this.props.STATUS === '404' || _.isEmpty(this.props.user)) {
      return (<RouteNotFound404Page />);
    } else {
      let namespace;
      var maybeStars = null;
      if(this.props.user.orgname) {
        namespace = this.props.user.orgname;
      } else {
        namespace = this.props.user.username;
        maybeStars = (<li><Link to={`/u/${namespace}/starred/`}>Stars</Link></li>);
      }
      return (
        <DocumentTitle title={`${namespace} - Docker Hub`}>
          <div>
            <SecondaryNav>
              <ul>
                <li><Link to={`/u/${namespace}/`}>Repos</Link></li>
                {maybeStars}
              </ul>
            </SecondaryNav>
            <div className='row'>
              <div className='large-3 columns'>
                <ProfileCard {...this.props.user} />
              </div>
              <div className='large-9 columns'>
                {this.props.children && cloneElement(this.props.children, {user: this.props.user})}
              </div>
            </div>
          </div>
        </DocumentTitle>
      );
    }
  }
});

export default connectToStores(UserProfile,
                               [
                                 UserProfileStore
                               ],
                               function({ getStore }, props) {
                                 return getStore(UserProfileStore)
                                               .getState();
                               });
