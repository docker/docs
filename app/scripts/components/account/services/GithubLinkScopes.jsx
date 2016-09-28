'use strict';

import React from 'react';
import GithubLinkStore from '../../../stores/GithubLinkStore';
import connectToStores from 'fluxible-addons-react/connectToStores';
const debug = require('debug')('COMPONENT:GithubLinkScopes');
import githubOathAction from '../../../actions/githubOauth';

import { SplitSection } from '../../common/Sections';
import { Button } from 'dux';

var GithubLinkScopes = React.createClass({
  displayName: 'GithubLinkScopes',
  contextTypes: {
    executeAction: React.PropTypes.func.isRequired
  },
  _generateGHStateString: function() {
    //Generate random string
    return Math.random().toString(36).substring(8);
  },
  _limitedScope: function(evt) {
    evt.preventDefault();
    var ss = this._generateGHStateString();
    window.open(
      'https://github.com/login/oauth/authorize?client_id=' +
      this.props.githubClientID +
      '&state=' +
      ss, '_blank');
    this.context.executeAction(githubOathAction, {stateString: ss});
  },
  _recommendedScope: function(evt) {
    evt.preventDefault();
    var ss = this._generateGHStateString();
    window.open(
      'https://github.com/login/oauth/authorize?client_id=' +
      this.props.githubClientID +
      '&scope=repo' +
      '&state=' +
      ss, '_blank');
    this.context.executeAction(githubOathAction, {stateString: ss});
  },
  render: function() {
    return (
      <div>
        <br />
        <SplitSection title='Connect to GitHub'
                      subtitle='We let you choose how much access we have to your GitHub account.'>
          <div className='row'>
            <div className='large-12 columns'>
              <h5>Public and Private (Recommended)</h5>
              <ul>
                <li>Read and Write access to public and private repositories.
                  (We only use write access to add service hooks and add deploy keys)</li>
                <li>Required if you want to setup an Automated Build from a private GitHub repository.</li>
                <li>Required if you want to use a private GitHub organization.</li>
                <li>We will automatically configure the service hooks and deploy keys for you.</li>
              </ul>
            </div>
            <div className='large-2 columns large-centered'>
              <Button size='small' onClick={this._recommendedScope}>Select</Button>
            </div>
            <div className='large-12 columns'>
              <h5>Limited Access</h5>
              <ul>
                <li>Public read only access.</li>
                <li>Only works with public repositories and organizations.</li>
                <li>You will need to manually make changes to your repositories in order to use Automated Build.</li>
              </ul>
            </div>
            <div className='large-2 columns large-centered'>
              <Button size='small' onClick={this._limitedScope}>Select</Button>
            </div>
          </div>
        </SplitSection>
      </div>
    );
  }
});

export default connectToStores(GithubLinkScopes,
                                [
                                  GithubLinkStore
                                ],
                                function({ getStore }, props) {
                                  return getStore(GithubLinkStore).getState();
                                });
