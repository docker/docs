'use strict';

import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import connectToStores from 'fluxible-addons-react/connectToStores';
import _ from 'lodash';

import styles from './AccountSettings.css';
import EmailForm from './forms/EmailForm';
import AccountInfoForm from './forms/AccountInfoForm';
import ChangePassForm from './forms/ChangePassForm';
import convertToOrgAction from '../../actions/convertToOrganization';
import toggleVisibility from '../../actions/toggleVisibility.js';
import PrivateRepoUsageStore from '../../stores/PrivateRepoUsageStore.js';
import { PageHeader, Button } from 'dux';
import classnames from 'classnames';
import { SplitSection } from '../common/Sections.jsx';
import DocumentTitle from 'react-document-title';

var debug = require('debug')('AccountSettings');

var Settings = React.createClass({
  contextTypes: {
    executeAction: React.PropTypes.func.isRequired
  },
  getDefaultProps: function() {
    return {
      user: {},
      JWT: ''
    };
  },
  render: function() {

    return (
      <DocumentTitle title='Account Settings - Docker Hub'>
        <div>
          <PageHeader title='Account Settings' />
          <div className={'row ' + styles.body}>
            <div className="columns large-12">
              <DefaultVisibility defaultRepoVisibility={this.props.defaultRepoVisibility}
                                 JWT={this.props.JWT}
                                 username={this.props.user.username}/>
              <EmailForm JWT={this.props.JWT} user={this.props.user}/>
              <ChangePassForm username={this.props.user.username} JWT={this.props.JWT}/>
              <AccountInfoForm JWT={this.props.JWT} user={this.props.user}/>
              <ToOrgForm userType={this.props.user.userType} history={this.props.history} username={this.props.user.username}/>
            </div>
          </div>
        </div>
      </DocumentTitle>
    );
  }
});

var DefaultVisibility = React.createClass({
  displayName: 'DefaultVisibility',
  PropTypes: {
    defaultRepoVisibility: PropTypes.string.isRequired
  },
  contextTypes: {
    executeAction: React.PropTypes.func.isRequired
  },
  toggleClick(visibility) {
    const _this = this;
    return (e) => {
      e.preventDefault();
      this.context.executeAction(toggleVisibility, _.merge(visibility, {JWT: _this.props.JWT, username: _this.props.username}));
    };
  },
  render: function() {

    return (
      <div>
        <SplitSection title="Default Repository Visibility"
                      subtitle={<p>Update the default visibility for your repositories.</p>}>
          <div>
            <form className={styles.margin}>
              <ul className='inline-list'>
                <li>
                  <label className={styles.visibility} onClick={this.toggleClick({visibility: 'public'})}>
                    <input type="radio"
                           name="defaultVisibility"
                           value="public"
                           checked={this.props.defaultRepoVisibility === 'public'}
                           readOnly /> public</label>
                </li>
                <li>
                  <label className={styles.visibility} onClick={this.toggleClick({visibility: 'private'})}>
                    <input type="radio"
                           name="defaultVisibility"
                           value="private"
                           checked={this.props.defaultRepoVisibility === 'private'}
                           readOnly /> private</label>
                </li>
              </ul>
            </form>
          </div>
        </SplitSection>
      </div>
    );
  }
});

var ToOrgForm = React.createClass({
  _toOrgClick: function(e) {
    e.preventDefault();
    this.props.history.pushState(null, '/account/convert-to-org/');
  },
  render: function() {

    if (this.props.userType === 'Organization') {
      return (
        <div className="to-org">
          <div className="form-section-header columns large-5">
            <div className="form-section-title">This is an Organization account for:</div>
            <div className="form-section-subtitle">{this.props.username}</div>
          </div>
          <br />
        </div>
      );
    } else {
      return (
        <SplitSection title="Convert your Account to an Organization"
                      module={false}
                      subtitle={<p>To use organization features you must convert your account from a "User" to an "Organization"</p>}>
          <div className={styles.margin}>
            <Button intent='warning' onClick={this._toOrgClick}>Convert to Organization</Button>
          </div>
        </SplitSection>
      );
    }
  }
});

export default connectToStores(Settings,
                               [
                                 PrivateRepoUsageStore
                               ],
                               function({ getStore }, props) {
                                 return getStore(PrivateRepoUsageStore).getState();
                               });
