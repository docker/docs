'use strict';

import React, { Component, PropTypes } from 'react';
import RepositoryList from 'components/common/repositoryList';
import { createStructuredSelector } from 'reselect';
import { connect } from 'react-redux';
import { getReposForUsername } from 'selectors/repositories';
import { Map } from 'immutable';
import { listRepositories } from 'actions/repositories';
import autoaction from 'autoaction';

import consts from 'consts';
import Spinner from 'components/common/spinner';

const mapState = createStructuredSelector({
  repos: getReposForUsername
});

@connect(mapState)
@autoaction({
  listRepositories: (props) => {
    return {
      namespace: props.params.username
    };
  }
}, {
  listRepositories
})
export default class UserRepos extends Component {

  static propTypes = {
    repos: PropTypes.instanceOf(Map)
  }

  render () {

    const {
      repos
    } = this.props;

    const status = [
      [consts.repositories.LIST_REPOSITORIES]
    ];

    return (
      <Spinner loadingStatus={ status }>
        <RepositoryList
          context='user'
          repositories={ repos.toArray() } />
      </Spinner>
    );
  }
}
