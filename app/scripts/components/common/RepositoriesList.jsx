'use strict';

import React, { PropTypes, createClass } from 'react';
const { array, bool, element, func, oneOfType } = PropTypes;
import PendingDeleteRepositoryItem from './PendingDeleteRepositoryItem';
import RepositoryListItem from './RepositoryListItem';
import { PENDING_DELETE } from 'common/enums/RepoStatus';
import has from 'lodash/object/has';

var debug = require('debug')('RepositoriesList');

function mkRepoListItem(repo) {

  const {
    is_automated,
    is_official,
    is_private,
    name,
    namespace,
    pull_count,
    repo_name,
    star_count,
    status
  } = repo;
  /**
   * TODO: after snakecase to camelcase in hub-js-sdk is done remove
   * explicit props
   */
  if(has(repo, 'is_private')) {
    repo.isPrivate = is_private;
  }
  if(has(repo, 'is_official')) {
    repo.isOfficial = is_official;
  }
  if(has(repo, 'is_offical')) {
    // This is a legit typo in the ES results
    repo.isOfficial = repo.is_offical;
  }
  if(has(repo, 'repo_name')) {
    repo.repoName = repo_name;
  }

  let key = repo_name || (namespace + '/' + name);

  if (status === PENDING_DELETE) {
    return (
      <PendingDeleteRepositoryItem {...repo} key={key}/>
    );
  } else {
    return (
      <RepositoryListItem {...repo}
                          key={key}
                          isAutomated={is_automated}
                          starCount={star_count}
                          pullCount={pull_count}/>
    );
  }
}

var BlankSlate = createClass({
  displayName: 'BlankSlate',
  render() {
    return (
      <div className="blankslate-alt">
            <div className="row">
              <div className='large-12 columns'>
                <h1>No Repositories Yet.</h1>
                <a href="/add/repository/"
                   className="button primary">Create Repository</a>
            </div>
          </div>
      </div>
    );
  }
});

export default createClass({
  displayName: 'RepositoriesList',
  propTypes: {
    blankSlate: element,
    repos: array.isRequired
  },
  getDefaultProps() {
    return {
      blankSlate: (<BlankSlate />),
      repos: []
    };
  },
  render() {
    const {
      blankSlate,
      repos
    } = this.props;

    let content = blankSlate;

    if(repos && repos.length > 0) {
      content = (
        <div className='profile-repos'>
          <ul className='large-12 columns no-bullet'>
            {repos.map(mkRepoListItem, this)}
          </ul>
        </div>
      );
    }

    return (
      <div className='row'>
        {content}
      </div>
    );
  }
});
