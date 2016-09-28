'use strict';

import React, { PropTypes, Component } from 'react';
import { Link } from 'react-router';
import FA from 'common/FontAwesome';
import Card, { Block } from '@dux/element-card';
import styles from './DeployKeys.css';
var debug = require('debug')('DeployKeys');

class DeployKeys extends Component {
  static propTypes = {
    autoBuildStore: PropTypes.object.isRequired
  }
  render() {
    const {
      autoBuildStore
    } = this.props;
    const pk = autoBuildStore.deploykey.public_key;
    let deployKeyHelp;
    if (autoBuildStore.provider.toLowerCase() === 'github') {
      deployKeyHelp = (
        <a href='https://help.github.com/articles/managing-deploy-keys'>
          <FA icon='fa-github' /> Deploy key help
        </a>
      );
    } else if (autoBuildStore.provider.toLowerCase() === 'bitbucket') {
      deployKeyHelp = (
        <a href='https://confluence.atlassian.com/display/BITBUCKET/Use+deployment+keys'>
          <FA icon='fa-bitbucket' /> Deploy key help
        </a>
      );
    }
    const subTitle = (
      <div>
        <p>
          If your automated build is linked to private repository on github or bitbucket,
          we need a way to have access to the repository. We do this with deploy keys.
          We try to do this step automatically for you, but sometimes we don't have access to do it.
          When this happens you need to add it yourself.
        </p>
        {deployKeyHelp}
      </div>
    );
    if (pk) {
      return (
        <Card heading='Deploy Key'>
          <Block>
            <h5>{subTitle}</h5>
            <textarea readonly className={styles.publicKeyText}>{pk}</textarea>
          </Block>
        </Card>);
    } else {
      return (<span />);
    }
  }
}

module.exports = DeployKeys;
