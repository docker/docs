'use strict';

import React, { Component, PropTypes } from 'react';
const { string, object, func } = PropTypes;
import { connectModal } from 'components/common/modal';
import DeleteModal from 'components/common/deleteModal';
import FontAwesome from 'components/common/fontAwesome';
import { revokeTeamAccessToRepo } from 'actions/repositories';
import { connect } from 'react-redux';
import { mapActions } from 'utils';

@connectModal()
@connect((() => {})(), mapActions({
  revokeTeamAccessToRepo
}))
export default class RevokeTeamAccess extends Component {

  static propTypes = {
    repo: string,
    org: string,
    id: string,
    team: object,
    actions: object,
    params: object,
    showModal: func,
    hideModal: func
  };

  removeTeam = () => {

    const {
      org,
      team,
      repo,
      id,
      actions,
      hideModal
    } = this.props;

    actions.revokeTeamAccessToRepo(
      {
        orgName: org,
        teamName: team.get('name'),
        repo: repo
      },
      id
    );

    hideModal();
  };

  confirmDelete() {
    this.props.showModal((
      <DeleteModal
        resourceType='team access'
        resourceName={ this.props.team.get('name') }
        onDelete={ ::this.removeTeam }
        hideModal={ this.props.hideModal }/>
    ));
  }
  render () {
    return (<FontAwesome icon='fa-close' onClick={ ::this.confirmDelete } />);
  }
}
