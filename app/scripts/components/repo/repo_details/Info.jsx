'use strict';

import React, { PropTypes, Component } from 'react';
import RepoShortDescription from './info/RepoShortDescription.jsx';
import RepoFullDescription from './info/RepoFullDescription.jsx';
import SourceRepositoryCard from 'common/SourceRepositoryCard';
import PullCommand from './info/PullCommand.jsx';
import Owner from './info/Owner.jsx';
import Comments from './info/Comments.jsx';
const { string, shape, object, bool, number } = PropTypes;

const debug = require('debug')('RepositoryDetailsInfo');


export default class RepositoryDetailsInfo extends Component {
  static propTypes = {
    description: string.isRequired,
    fullDescription: string.isRequired,
    isPrivate: bool.isRequired,
    isAutomated: bool.isRequired,
    name: string.isRequired,
    namespace: string.isRequired,
    status: number.isRequired,
    user: object,
    JWT: string,
    canEdit: bool,
    autoBuildStore: shape({
      provider: string.isRequired,
      repo_web_url: string.isRequired
    })
  }

  render() {
    const {
            namespace,
            name,
            isAutomated,
            JWT,
            user
          } = this.props;
    const { provider, repo_web_url } = this.props.autoBuildStore;

    const pullCommandCard = <PullCommand name={name} namespace={namespace} />;
    const ownerCard = <Owner namespace={namespace} />;
    let sourceCard = null;
    if (isAutomated) {
      sourceCard = <SourceRepositoryCard provider={provider} url={repo_web_url} />;
    }

    return (
      <div>
        <div className='row'>
          <div className='large-8 columns'>
            <RepoShortDescription {...this.props} />
            <RepoFullDescription {...this.props} />
          </div>
          <div className='large-4 columns'>
            {pullCommandCard}
            {ownerCard}
            {sourceCard}
          </div>
        </div>
        <div className='row'>
          <div className='large-8 columns'>
            <Comments JWT={JWT}
                      user={user}
                      name={name}
                      namespace={namespace} />
          </div>
        </div>
      </div>
    );
  }

}
