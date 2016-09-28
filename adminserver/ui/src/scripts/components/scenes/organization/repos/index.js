'use strict';

import React, { Component, PropTypes } from 'react';
const { bool, instanceOf, func, object } = PropTypes;
import { Map } from 'immutable';
// Selectors
import { getTeamOrOrgRepos } from 'selectors/repositories';
import {
  getIsCurrentUserPageOrgAdmin
} from 'selectors/organizations';
import { createStructuredSelector } from 'reselect';
import { connect } from 'react-redux';
// Components
import RepositoryList from 'components/common/repositoryList';
import Button from 'components/common/button';
import Spinner from 'components/common/spinner';
import AddRepoForm from './addRepoForm';
// Actions
import {
  listOrgOrTeamRepos,
  changeTeamAccessToRepo,
  revokeTeamAccessToRepo
} from 'actions/repositories';
import autoaction from 'autoaction';
import { mapActions } from 'utils';
// Misc
import ui from 'redux-ui';
import consts from 'consts';
import css from 'react-css-modules';
import styles from './repos.css';

const mapState = createStructuredSelector({
  repos: getTeamOrOrgRepos,
  isOrgAdmin: getIsCurrentUserPageOrgAdmin
});

@connect(mapState, mapActions({ changeTeamAccessToRepo, revokeTeamAccessToRepo }))
@ui({
  state: {
    isFormVisible: false
  },
  reducer: (state, action) => {
    if (action.type === consts.repositories.GRANT_TEAM_ACCESS_TO_REPO) {
      return state.set('isFormVisible', false);
    }
    return state;
  }
})
@autoaction({
  listOrgOrTeamRepos: (props) => ({ orgName: props.params.org, teamName: props.params.team || '' })
}, { listOrgOrTeamRepos })
@css(styles)
export default class OrgRepos extends Component {

  static propTypes = {
    actions: object,
    repos: instanceOf(Map),
    isOrgAdmin: bool,
    ui: object,
    updateUI: func,
    params: object
  }

  onRemovePermissions(repoRecord) {
    const {
    org,
    team
    } = this.props.params;
    this.props.actions.revokeTeamAccessToRepo({
      orgName: org,
      teamName: team,
      repo: repoRecord.name
    });
  }

  onEditPermissions = (repoRecord, accessLevel) => {
    const {
      org,
      team
    } = this.props.params;
    this.props.actions.changeTeamAccessToRepo({
      orgName: org,
      teamName: team,
      repo: repoRecord.name,
      accessLevel
    });
  }

  showAddForm() {
    this.props.updateUI('isFormVisible', true);
  }

  hideAddForm() {
    this.props.updateUI('isFormVisible', false);
  }

  render() {
    const {
      isOrgAdmin,
      repos,
      ui: {
        isFormVisible
      },
      params: {
        org,
        team
      }
    } = this.props;

    // We wait for either the org repos or team repos to finish loading
    const status = [
      [consts.repositories.LIST_TEAM_ACCESS_TO_REPO, org, team],
      [consts.repositories.LIST_REPOSITORIES, org]
    ];

    return (
      <div>
        { isOrgAdmin ?
          <div styleName='actions'>
            <Button variant='secondary' onClick={ ::this.showAddForm } disabled={ isFormVisible }>Add repository</Button>
          </div>
        : undefined }

        { isFormVisible ? <AddRepoForm onHide={ ::this.hideAddForm } params={ this.props.params } /> : undefined }

        <Spinner loadingStatus={ status }>
          <RepositoryList
            context='team'
            canEditPermissions={ isOrgAdmin }
            canRemovePermissions={ isOrgAdmin }
            onEditPermissions={ team ? ::this.onEditPermissions : undefined }
            onRemovePermissions={ ::this.onRemovePermissions }
            repositories={ repos.toArray() } />
        </Spinner>
      </div>
    );
  }

}
