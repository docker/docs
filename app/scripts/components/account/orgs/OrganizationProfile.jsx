'use strict';

import React from 'react';
import connectToStores from 'fluxible-addons-react/connectToStores';
import _ from 'lodash';

import saveOrgProfileAction from '../../../actions/saveOrgProfile';
import toggleVisibility from '../../../actions/toggleVisibility.js';
import OrganizationStore from '../../../stores/OrganizationStore';
import PrivateRepoUsageStore from '../../../stores/PrivateRepoUsageStore.js';
import SimpleInput from 'common/SimpleInput.jsx';
import { SplitSection } from '../../common/Sections.jsx';
import Route404 from '../../common/RouteNotFound404Page.jsx';
import Button from '@dux/element-button';
import styles from './OrganizationProfile.css';

var OrganizationProfile = React.createClass({
  contextTypes: {
    executeAction: React.PropTypes.func.isRequired
  },
  propTypes: {
    currentOrg: React.PropTypes.shape({
      orgname: React.PropTypes.string,
      full_name: React.PropTypes.string,
      location: React.PropTypes.string,
      company: React.PropTypes.string,
      profile_url: React.PropTypes.string
    }),
    isOwner: React.PropTypes.bool.isRequired,
    JWT: React.PropTypes.string,
    error: React.PropTypes.string,
    success: React.PropTypes.string,
    defaultVisibility: React.PropTypes.oneOf(['public', 'private'])
  },
  getInitialState: function() {
    return {
      orgname: this.props.currentOrg.orgname,
      fullName: this.props.currentOrg.full_name,
      location: this.props.currentOrg.location,
      company: this.props.currentOrg.company,
      profileUrl: this.props.currentOrg.profile_url,
      gravatarEmail: this.props.currentOrg.gravatar_email,
      defaultVisibility: this.props.defaultVisibility
    };
  },
  onSubmit: function(e) {
    e.preventDefault();
    var updatedOrg = {
      full_name: this.state.fullName,
      location: this.state.location,
      company: this.state.company,
      profile_url: this.state.profileUrl,
      gravatar_email: this.state.gravatarEmail
    };
    this.context.executeAction(saveOrgProfileAction, {
      jwt: this.props.JWT,
      orgname: this.state.orgname,
      organization: updatedOrg
    });
    this.context.executeAction(toggleVisibility, {
      JWT: this.props.JWT,
      username: this.state.orgname,
      visibility: this.state.defaultVisibility
    });
  },
  orgFullNameChange: function(e) {
    this.setState({fullName: e.target.value});
  },
  orgCompanyChange: function(e) {
    this.setState({company: e.target.value});
  },
  locationChange: function(e) {
    this.setState({location: e.target.value});
  },
  profileUrlChange: function(e) {
    this.setState({profileUrl: e.target.value});
  },
  gravatarEmailChange: function(e) {
    this.setState({gravatarEmail: e.target.value});
  },
  toggleClick: function({ visibility }) {
    return (e) => {
      this.setState({defaultVisibility: visibility});
    };
  },
  render: function() {
    var maybeError = <span />;
    var maybeSuccess = <span />;
    if (this.props.error) {
      maybeError = <span className='alert-box alert radius'>{this.props.error}</span>;
    } else if(this.props.success) {
      //TODO: this could be an alert box with a close icon that can be closed when the user wants to dismiss
      //Or time out after some time (ideally i would like to see this as a notification and that's it)
      maybeSuccess = <div><br /><span className='alert-box success radius'>{this.props.success}</span></div>;
    }
    if (this.props.isOwner) {
      return (
        <div>
          <br />
          <SplitSection title={'Organization: ' + this.state.orgname}
                        subtitle={<p>This information is private to users with access to this organization.</p>}>
            <form onSubmit={this.onSubmit}>
              <br />
              <div className="row">
                <div className='columns large-3'>
                  Default Visibility:
                </div>
                <ul className='inline-list columns large-9'>
                  <li>
                    <label className={styles.visibility} onClick={this.toggleClick({visibility: 'public'})}>
                      <input type="radio"
                             name="defaultVisibility"
                             value="public"
                             checked={this.state.defaultVisibility === 'public'}
                             readOnly/> public</label>
                  </li>
                  <li>
                    <label className={styles.visibility} onClick={this.toggleClick({visibility: 'private'})}>
                      <input type="radio"
                             name="defaultVisibility"
                             value="private"
                             checked={this.state.defaultVisibility === 'private'}
                             readOnly/> private</label>
                  </li>
                </ul>
              </div>

              <label className={styles.label}>Organization Full Name</label>
              <SimpleInput type="text"
                           onChange={this.orgFullNameChange}
                           value={this.state.fullName}/>

              <label className={styles.label}>Company</label>
              <SimpleInput type="text"
                           onChange={this.orgCompanyChange}
                           value={this.state.company}/>

              <label className={styles.label}>Location</label>
              <SimpleInput type="text"
                           onChange={this.locationChange}
                           value={this.state.location}/>

              <label className={styles.label}>Website</label>
              <SimpleInput type="text"
                           onChange={this.profileUrlChange}
                           value={this.state.profileUrl}/>

              <label className={styles.label}>Gravatar Email</label>
              <SimpleInput type="text"
                           onChange={this.gravatarEmailChange}
                           value={this.state.gravatarEmail}/>

              <Button type="submit">Save</Button>
              {maybeError}
              {maybeSuccess}
            </form>
          </SplitSection>
        </div>
      );
    } else {
      return (
        <Route404 />
      );
    }
  }
});

export default connectToStores(OrganizationProfile,
  [
    OrganizationStore,
    PrivateRepoUsageStore
  ],
  function({ getStore }, props) {
    return _.merge({}, getStore(OrganizationStore).getState(), {defaultVisibility: getStore(PrivateRepoUsageStore).getState().defaultRepoVisibility});
  });
