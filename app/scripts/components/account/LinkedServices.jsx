'use strict';

import React, { PropTypes, Component } from 'react';
import connectToStores from 'fluxible-addons-react/connectToStores';
import { PageHeader } from 'dux';

import AutobuildStore from '../../stores/AutobuildStore';
import BitbucketLinkStore from '../../stores/BitbucketLinkStore';
import GithubLinkStore from '../../stores/GithubLinkStore';
import GithubLinkAction from '../../actions/linkGithub';
import BitbucketLinkAction from '../../actions/linkBitbucket';
import GithubUnlinkAction from '../../actions/unlinkGithub';
import BitbucketUnlinkAction from '../../actions/unlinkBitbucket';
import { Button, Module } from 'dux';
import { SplitSection } from '../common/Sections.jsx';
import FA from '../common/FontAwesome';
import classnames from 'classnames';
import styles from './LinkedServices.css';
import merge from 'lodash/object/merge';

const { string, array, object, func } = PropTypes;

class LinkedServices extends Component {

  static contextTypes = {
    executeAction: func.isRequired
  }

  static propTypes = {
    JWT: string.isRequired,
    githubAccount: object,
    githubRepos: array,
    bitbucketAccount: object,
    bitbucketRepos: array,
    gitlabAccount: object,
    gitlabRepos: array,
    bitbucketError: string,
    bbAuthUrl: string,
    githubError: string
  }

  authServiceClick = (e) => {
    e.preventDefault();
  }

  render() {
    var github = {service: 'Github', account: this.props.githubAccount};
    var bitbucket = {service: 'Bitbucket', account: this.props.bitbucketAccount};
    var gitlab = {service: 'Gitlab', account: this.props.gitlabAccount};
    let maybeError;
    const errorMsg = this.props.githubError || this.props.bitbucketError;
    if (errorMsg) {
      maybeError = <span className='alert-box alert radius'>{errorMsg}</span>;
    }
    return (
      <div>
        <PageHeader title='Linked Accounts & Services' />
        <div className={'row ' + styles.body}>
          <div className="columns large-12">

            <SplitSection title='Linked Accounts'
                          subtitle={<p>These account links are currently used for Automated Builds,
                          so that we can access your project lists and help you configure your Automated Builds.
                          <strong>&nbsp; Please note: A github/bitbucket account can be connected to only one docker hub account at a time.</strong></p>}>
                <div className="row">
                  <LinkedAccount data={github}
                                 JWT={this.props.JWT}
                                 history={this.props.history}/>
                  <LinkedAccount data={bitbucket}
                                 bbAuthUrl={this.props.bbAuthUrl}
                                 JWT={this.props.JWT}
                                 history={this.props.history}/>
                </div>
                <br />
                {maybeError}
            </SplitSection>
          </div>
        </div>
      </div>
    );
  }
}

class LinkedAccount extends Component {
  static propTypes = {
    data: object,
    bbAuthUrl: string,
    JWT: string,
    history: object.isRequired
  }

  static contextTypes = {
    executeAction: PropTypes.func.isRequired
  }

  linkAction = (provider, e) => {
    e.preventDefault();
    if (provider.service.toLowerCase() === 'github') {
      this.context.executeAction(GithubLinkAction, this.props.JWT);
      this.props.history.pushState(null, '/account/authorized-services/github-permissions/');
    } else if (provider.service.toLowerCase() === 'bitbucket') {
      const bbWin = window.open();
      bbWin.location = this.props.bbAuthUrl;
    }
  }

  unlinkAction = (provider, e) => {
    e.preventDefault();
    if (provider.service.toLowerCase() === 'github') {
      this.context.executeAction(GithubUnlinkAction, this.props.JWT);
    } else if (provider.service.toLowerCase() === 'bitbucket') {
      this.context.executeAction(BitbucketUnlinkAction, this.props.JWT);
    }
  }

  render() {
    var service = this.props.data.service;
    var icon;
    if (service.toLowerCase() === 'github') {
      icon = 'fa-github';
    } else if (service.toLowerCase() === 'bitbucket') {
      icon = 'fa-bitbucket';
    }
    var account = this.props.data.account;
    let linkClass = classnames({
      [styles.service]: true,
      [styles.unlink]: account,
      [styles.link]: !account
    });
    if (account) {
      return (
        <div className='columns large-6'>
          <div className={linkClass} onClick={this.unlinkAction.bind(null, this.props.data)}>
            <div className={'row ' + styles.title}>
              <div className="columns large-5">
                <img src={account.avatar_url} className={styles.icon} />
              </div>
              <div className={'columns large-7 ' + styles.access}>
                <span>{account.login}:</span><br/>
                <span>read/write access</span>
              </div>
            </div><br/>
            <span className={styles.name}>Unlink {service}</span>
          </div>
        </div>
      );
    } else {
      return (
      <div className='columns large-6'>
        <div className={linkClass} onClick={this.linkAction.bind(null, this.props.data)}>
          <FA icon={icon}/><br/>
          <span className={styles.name}>Link {service}</span>
        </div>
      </div>
      );
    }
  }
}

export default connectToStores(LinkedServices,
  [
    AutobuildStore,
    BitbucketLinkStore,
    GithubLinkStore
  ],
  function({ getStore }, props) {
    return merge(
      {},
      getStore(AutobuildStore).getState(),
      { bitbucketError: getStore(BitbucketLinkStore).getState().error,
        bbAuthUrl: getStore(BitbucketLinkStore).getState().authURL,
        githubError: getStore(GithubLinkStore).getState().error }
    );
  });
