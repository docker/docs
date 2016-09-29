'use strict';

import React, { PropTypes, Component } from 'react';
import connectToStores from 'fluxible-addons-react/connectToStores';
import Card, { Block } from '@dux/element-card';
import SourceRepositoryCard from 'common/SourceRepositoryCard';
import RepoDetailsDockerfileStore from '../../../stores/RepoDetailsDockerfileStore';
var debug = require('debug')('Dockerfile');
import styles from './Dockerfile.css';
import Code from 'common/Code';

class Dockerfile extends Component {

  static propTypes = {
    dockerfile: PropTypes.string.isRequired
  }

  render() {
    const { provider, repo_web_url } = this.props.autoBuildStore;
    return (
      <div className='row'>
        <div className='large-8 columns'>
          <Card heading='Dockerfile'>
            <span className={styles.dockerfile}><Code>{this.props.dockerfile}</Code></span>
          </Card>
        </div>
        <div className='large-4 columns'>
          <SourceRepositoryCard provider={provider} url={repo_web_url} />
        </div>
      </div>
    );
  }
}

export default connectToStores(Dockerfile,
                               [
                                   RepoDetailsDockerfileStore
                               ],
                               function({ getStore }, props) {
                                   return getStore(RepoDetailsDockerfileStore).getState();
                               });
