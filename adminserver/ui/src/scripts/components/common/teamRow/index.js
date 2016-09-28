'use strict';

import React, { Component, PropTypes } from 'react';
const { string, object } = PropTypes;
import autoaction from 'autoaction';
import { getTeam } from 'actions/teams';
import { teamDetails } from 'selectors/teams';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { changeTeamAccessToRepo } from 'actions/repositories';
import { mapActions } from 'utils';
import RevokeTeamAccess from './revokeTeamAccess';


const mapState = createStructuredSelector({
  team: teamDetails
});

@connect(mapState, mapActions({
  changeTeamAccessToRepo
}))
@autoaction({
  getTeam: (props) => {
    return [props.org, `id:${props.id}`];
  }
}, {
  getTeam
})
export default class TeamRow extends Component {

  static propTypes = {
    repo: string,
    org: string,
    id: string,
    permission: string,
    team: object,
    actions: object,
    params: object
  };

  changePermission = (evt) => {

    const {
      org,
      team,
      repo,
      actions
    } = this.props;

    actions.changeTeamAccessToRepo({
      orgName: org,
      teamName: team.get('name'),
      repo: repo,
      accessLevel: evt.target.value
    });
  };

  render () {

    const {
      team,
      org,
      id,
      repo,
      permission
    } = this.props;

    return (
      <tr>
        <td>
          { team.get('name') }
        </td>
        <td>
          <select value={ permission } onChange={ ::this.changePermission }>
            <option value='admin'>Admin</option>
            <option value='read-write'>Read-Write</option>
            <option value='read-only'>Read Only</option>
          </select>
        </td>
        <td>
          <RevokeTeamAccess
            repo={ repo }
            org={ org }
            id={ id }
            team={ team }
          />
        </td>
      </tr>
    );
  }
}
