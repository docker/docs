'use strict';

import React, { PropTypes, Component } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import Immutable from 'immutable';
import * as RepositoriesActions from 'actions/repositories';
import { resetStatus } from 'actions/status';
import Spinner from 'components/common/spinner';
import { Tabs, Tab } from 'components/common/tabs';
import consts from 'consts';
import { mapActions } from 'utils';
import autoaction from 'autoaction';
import styles from './repository.css';
import css from 'react-css-modules';
import FontAwesome from 'components/common/fontAwesome';
import {
  getRawRepoState,
  getAccessLevel,
  getRepoForName
} from 'selectors/repositories';

const mapRepositoryState = (state, props) => ({
  repositories: getRawRepoState,
  status: state.status,
  accessLevel: getAccessLevel(state, props),
  repo: getRepoForName(state, props)
});

@connect(
  mapRepositoryState,
  mapActions({
    ...RepositoriesActions,
    resetStatus
  })
)
@autoaction({
  resetStatus: {
    key: (p) => `${p.params.namespace}/${p.params.name}`
  },
  // On this component we use the 'getUserRepoAccess' endpoint, which returns
  // both the repository object *and* the user access level in one call.
  getUserRepoAccess: (props) => {
    return {
      user: window.user.name,
      namespace: props.params.namespace,
      repo: props.params.name
    };
  }
}, {
  ...RepositoriesActions,
  resetStatus
})
@css(styles, {allowMultiple: true})
export default class Repository extends Component {

  static propTypes = {
    actions: PropTypes.object,
    children: PropTypes.any,
    params: PropTypes.object,
    repositories: PropTypes.object,
    status: PropTypes.instanceOf(Immutable.Map),
    accessLevel: PropTypes.string,
    repo: PropTypes.object
  };

  static contextTypes = {
    router: PropTypes.object.isRequired
  };

  componentWillReceiveProps(next) {
    const {
      params: {
        name,
        namespace
      }
    } = next;
    const deleteRepoStatus = next.status.getIn([consts.repositories.DELETE_REPOSITORY, namespace, name, 'status']);
    if (deleteRepoStatus === consts.loading.SUCCESS) {
      // Repo deleted
      // we use replace rather than push so the deleted history item is inaccessible
      this.context.router.replace('/repositories');
    }
  }

  render() {
    const {
      namespace,
      name
      } = this.props.params;
    const getStatus = [consts.repositories.GET_REPOSITORY_WITH_USER_PERMISSIONS, namespace, name, window.user.name];
    const getOwnershipStatus = [consts.teams.GET_TEAM_MEMBER, namespace, 'owners', window.user.name];

    const {
      accessLevel,
      repo
      } = this.props;

    const {
      visibility,
      namespaceType,
      shortDescription
      } = repo || {};

    return (
      <div styleName='wrapper'>
        <Link to={ '/repositories' } styleName='backToRepos'>
          <FontAwesome icon='fa-long-arrow-left'/>
          Back to repositories
        </Link>

        <div styleName='repositoryDetail'>

          <span styleName='book'><FontAwesome icon='fa-book '/></span>

          <div styleName='repoInfo'>
            { namespaceType === 'user' ?
              <Link to={ `/users/${namespace}` }>{ namespace }</Link>
              :
              <Link to={ `/orgs/${namespace}` }>{ namespace }</Link>
            }
            <span styleName='divider'>/</span>
            <span styleName='name'>{ name }</span>
            { visibility === 'private' ?
              <span styleName='visibility'>private</span> :
              null
            }
            <p styleName='shortDescription'>
              { shortDescription ?
                shortDescription :
                (<i>Description is empty for this repository.</i>)
              }
            </p>
          </div>
        </div>

        <Tabs header id='repository-tabs'>
          <Tab><Link to={ `/repositories/${namespace}/${name}/details` }>Info</Link></Tab>
          { accessLevel === 'admin' && namespaceType === 'organization' &&
          <Tab><Link to={ `/repositories/${namespace}/${name}/teams` }>Permissions</Link></Tab>
          }
          <Tab><Link to={ `/repositories/${namespace}/${name}/tags` }>Tags</Link></Tab>
          { accessLevel === 'admin' &&
          <Tab id='repo-settings-tab'><Link to={ `/repositories/${namespace}/${name}/settings` }>Settings</Link></Tab>
          }
        </Tabs>

        <Spinner loadingStatus={ [getStatus, getOwnershipStatus] }>
          { this.props.children }
        </Spinner>

      </div>
    );
  }
}

export { RepositoryDetailsTab } from './details';
export { RepositoryTeamsTab } from './teams';
export { RepositoryUsersTab } from './users';
export { RepositorySettingsTab } from './settings';
export { RepositoryTagsTab } from './tags';
export RepositoryActivityTab from './activity';
