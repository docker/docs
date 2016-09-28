'use strict';

import AutobuildStore from '../../stores/AutobuildStore';
import Route404 from '../common/RouteNotFound404Page.jsx';
import AutobuildBlankSlate from './AutobuildBlankSlate.jsx';
import React, { PropTypes, cloneElement } from 'react';
import { Link } from 'react-router';
import connectToStores from 'fluxible-addons-react/connectToStores';
import { SecondaryNav } from 'dux';
import LiLink from '../common/LiLink';
const debug = require('debug')('COMPONENT:Autobuild');

let AutobuildSourcesNav = React.createClass({
  displayName: 'AutobuildSourcesNav',
  propTypes: {
    namespace: PropTypes.string,
    githubAccount: PropTypes.object,
    bitbucketAccount: PropTypes.object
  },
  _checkAccountsAndGetLinks: function() {
    var liLinks = [];
    if (this.props.githubAccount) {
      liLinks.push(
      <LiLink to={`/add/automated-build/${this.props.namespace}/github/orgs/`} key="github">
        <i className='fa fa-github'/> GitHub ({this.props.githubAccount.login})
      </LiLink>);
    }

    if (this.props.bitbucketAccount) {
      liLinks.push(
        <LiLink to={`/add/automated-build/${this.props.namespace}/bitbucket/orgs/`} key="bitbucket">
          <i className='fa fa-bitbucket'/> Bitbucket ({this.props.bitbucketAccount.login})
        </LiLink>);
    }

    if (!this.props.githubAccount || !this.props.bitbucketAccount) {
      liLinks.push(<LiLink to="/account/authorized-services/" key="link"><i className='fa fa-gear' /> Link Accounts</LiLink>);
    }

    return liLinks;
  },
  render: function() {
    var liLinks = this._checkAccountsAndGetLinks();
    return (
      <div className="secondary-nav">
        <SecondaryNav>
          <ul>
            {liLinks}
          </ul>
        </SecondaryNav>
      </div>
    );
  }
});

var Autobuild = React.createClass({
  propTypes: {
    user: React.PropTypes.object,
    JWT: React.PropTypes.string
  },
  contextTypes: {
    executeAction: React.PropTypes.func.isRequired
  },
  render: function() {
    if (!this.props.JWT) {
      return (
        <Route404 />
      );
    } else if (!this.props.githubAccount && !this.props.bitbucketAccount) {
      return (
        <div className="settings-wrapper">
          <AutobuildSourcesNav namespace={this.props.params.userNamespace}
                               githubAccount={this.props.githubAccount}
                               bitbucketAccount={this.props.bitbucketAccount}/>
          <AutobuildBlankSlate/>
        </div>
      );
    } else {
      return (
        <div className="settings-wrapper">
          <AutobuildSourcesNav namespace={this.props.params.userNamespace}
                               githubAccount={this.props.githubAccount}
                               bitbucketAccount={this.props.bitbucketAccount}/>
          {this.props.children && cloneElement(this.props.children, {
            user: this.props.user,
            JWT: this.props.JWT,
            githubAccount: this.props.githubAccount,
            bitbucketAccount: this.props.bitbucketAccount
          })}
        </div>
      );
    }
  }
});

export default connectToStores(Autobuild,
  [
    AutobuildStore
  ],
  function({ getStore }, props) {
    return getStore(AutobuildStore).getState();
  });
