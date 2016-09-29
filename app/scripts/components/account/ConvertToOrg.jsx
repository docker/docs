'use strict';

import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import _ from 'lodash';
import connectToStores from 'fluxible-addons-react/connectToStores';
import classnames from 'classnames';

import DUXInput from '../common/DUXInput.jsx';
import convertToOrgAction from '../../actions/convertToOrganization';
import updateToOrgOwner from '../../actions/updateToOrgOwner.js';
import ConvertToOrgStore from '../../stores/ConvertToOrgStore.js';
import { PageHeader, Module } from 'dux';
import Button from '@dux/element-button';

import { FullSection } from '../common/Sections.jsx';

import styles from './ConvertToOrg.css';

var debug = require('debug')('ConvertToOrg');

var _mkErrors = function(err, key) {
  const errorClass = classnames({
    [ styles.center ]: true,
    [ styles.warning ]: true
  });
  return (
    <div key={ key } className={errorClass}>
      { err }
    </div>
  );
};

var ConvertToOrg = React.createClass({
  contextTypes: {
    executeAction: PropTypes.func.isRequired
  },
  propTypes: {
    user: PropTypes.shape({
      username: PropTypes.string
    }),
    JWT: PropTypes.string,
    convertError: PropTypes.bool.isRequired,
    error: PropTypes.object
  },
  getInitialState: function() {
    return {
      newOwner: ''
    };
  },
  onChangeOwner: function(e) {
    e.preventDefault();
    var newOwner = e.target.value;
    this.context.executeAction(updateToOrgOwner, { newOwner: newOwner });
  },
  onCancelClick: function(e) {
    e.preventDefault();
    this.props.history.pushState(null, '/account/settings/');
  },
  submitChangeOrg: function(e) {
    e.preventDefault();
    debug('Change user to org');
    this.context.executeAction(convertToOrgAction,
      {jwt: this.props.JWT, username: this.props.user.username, newOwner: this.props.newOwner});
  },
  render: function() {
    var disabled = !this.props.newOwner;
    var intent = this.props.convertError ? 'alert' : null;
    let error = null;
    if (this.props.convertError) {
      error = (
        <div>
          { _.map(this.props.error, _mkErrors) }
        </div>
      );
    }
    return (
      <div>
        <PageHeader title='Convert your Account to an Organization' />
        <div className={'row ' + styles.body}>
          <div className="columns large-8 large-offset-2 end">
            <FullSection title='How it works'>
              <div className={'columns large-12 ' + styles.toOrgContent}>Your user account will be transformed into an organization account where all administrative duties are left to another user or group of users. You will no longer be able to login to this account.</div>
            </FullSection>
            <FullSection title='Email Addresses'>
              <div className={'columns large-12 ' + styles.toOrgContent}>Email addresses for this account will be removed, freeing them up to be used for any other accounts.</div>
            </FullSection>
            <FullSection title='Linked Accounts'>
              <div className={'columns large-12 ' + styles.toOrgContent}>Converting your account removes any associations to other services like GitHub or Atlassian Bitbucket. You will be able to link your external accounts to another Docker Hub user.</div>
            </FullSection>
            <FullSection title='Billing'>
              <div className={'columns large-12 ' + styles.toOrgContent}>Billing details and Private Repository plans will remain attached to this account after it is converted to an organization.</div>
            </FullSection>
            <FullSection title='Repositories'>
              <div className={'columns large-12 ' + styles.toOrgContent}>Repository namespaces and names remain unchanged. Any user collaborators that you have configured for these repositories will be removed and must be reconfigured using group collaborators.</div>
            </FullSection>
            <FullSection title='Automated Builds'>
              <div className={'columns large-12 ' + styles.toOrgContent}>Automated Builds for this account will be updated to appear as if they were originallly configured by the initial organization owner. Any user in a group with 'admin' level access to a repository will be able to edit Automated Build Configurations.</div>
            </FullSection>
            <div className="row">
              <div className={'columns large-12 ' + styles.center}>
                <h5 className={'toOrg-title ' + styles.warning}>WARNING</h5>
                <div className={'toOrg-body ' + styles.warning}>This account conversion operation can not be undone.</div>
              </div>
            </div>
            <div className="row">
              <div className="columns large-12">
                <p>In order to complete the conversion of your account to an organization you will need to enter the Docker ID of an **existing** Docker Hub user account.
                  The user account you specify will become a member of the Owners group and will have full administrative privileges to manage the organization.</p>
                <Module intent={intent}>
                  { error }
                  <form className={styles.form} onSubmit={this.submitChangeOrg}>
                    <div className="row">
                      <div className="columns large-7">
                        <DUXInput label='Existing Docker ID (username required)'
                                  onChange={this.onChangeOwner}
                                  value={this.props.newOwner}/>
                      </div>
                      <div className="columns large-5">
                        <div className="row">
                          <div className='columns large-5'>
                            <Button variant='alert' size='small' ghost onClick={this.onCancelClick}>Cancel</Button>
                          </div>
                          <div className='columns large-7'>
                            <Button type='submit' size='small' disabled={disabled}>Save and Continue</Button>
                          </div>
                        </div>
                      </div>
                    </div>
                  </form>
                </Module>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
});

export default connectToStores(ConvertToOrg,
  [ConvertToOrgStore],
  function({ getStore }, props) {
    return getStore(ConvertToOrgStore).getState();
  });
