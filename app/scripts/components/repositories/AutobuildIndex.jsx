'use strict';

import React, { PropTypes, Component } from 'react';
import { Link } from 'react-router';
import FA from 'common/FontAwesome';
import classnames from 'classnames';
import styles from './AutobuildIndex.css';

export default class AutobuildIndex extends Component {
  static propTypes = {
    githubAccount: PropTypes.object,
    bitbucketAccount: PropTypes.object
  }

  render() {
    let githubText;
    let githubLink;
    let bitbucketText;
    let bitbucketLink;
    const linkClasses = classnames({
      'button': true,
      [styles.link]: true
    });

    if (this.props.githubAccount) {
      githubLink = `/add/automated-build/${this.props.params.userNamespace}/github/orgs/`;
      githubText = (
        <h2>Create Auto-build</h2>
      );
    } else {
      githubLink = '/account/authorized-services/';
      githubText = (
        <h2>Link Account</h2>
      );
    }
    if (this.props.bitbucketAccount) {
      bitbucketLink = `/add/automated-build/${this.props.params.userNamespace}/bitbucket/orgs/`;
      bitbucketText = (
        <h2>Create Auto-build</h2>
      );
    } else {
      bitbucketLink = '/account/authorized-services/';
      bitbucketText = (
        <h2>Link Account</h2>
      );
    }

    return (
      <div className={`row ${styles.row}`}>
        <div className='large-4 large-offset-2 columns'>
          <Link className={linkClasses} to={githubLink}>
            {githubText}
            <FA icon='fa-github'/>
            <h2>Github</h2>
          </Link>
        </div>
        <div className='large-4 columns end'>
          <Link className={linkClasses} to={bitbucketLink}>
            {bitbucketText}
            <FA icon='fa-bitbucket'/>
            <h2>Bitbucket</h2>
          </Link>
        </div>
      </div>
    );
  }
}
