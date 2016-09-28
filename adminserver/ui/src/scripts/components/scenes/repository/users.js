'use strict';

import React, { Component } from 'react';
import { connect } from 'react-redux';
import * as RepositoriesActions from 'actions/repositories';
import Spinner from 'components/common/spinner';
import MemberList from 'components/scenes/organization/memberList';
import consts from 'consts';
import { mapActions } from 'utils';
import autoaction from 'autoaction';

const mapRepositoryState = (state) => ({
  currentPage: state.repositories.getIn(['ui', 'repoUsersPage']),
  repositories: state.repositories
});

@connect(mapRepositoryState, mapActions(RepositoriesActions))
@autoaction({
  listRepoUserAccess: (_, state) => { return { namespace: state.router.params.namespace, repo: state.router.params.name }; }
}, RepositoriesActions)
export class RepositoryUsersTab extends Component {

  static propTypes = {
    actions: React.PropTypes.object,
    currentPage: React.PropTypes.number,
    params: React.PropTypes.object.isRequired,
    repositories: React.PropTypes.object
  }

  render() {
    const {namespace, name} = this.props.params;
    const listRepoUserAccessStatus = [consts.repositories.LIST_REPO_USER_ACCESS, namespace, name];

    let users = [{name: namespace}]; // TODO use a full user object
    (this.props.repositories.getIn(['repositories', namespace, name, 'access', 'users']) || [])
      .forEach((x) => {
        if (x.name === namespace) {
          users[0] = x.user;
        } else {
          users.push(x.user);
        }
      });

    return (
      <div>
        <Spinner loadingStatus={ [listRepoUserAccessStatus] }>
          <MemberList
            members={ users }
            canDelete={ false } />
        </Spinner>
      </div>
    );
  }
}
