'use strict';

import React, { Component, PropTypes } from 'react';

import PaginatedTable from 'components/common/paginatedTable';
import autoaction from 'autoaction';
import { getTeamsForUser } from 'actions/teams';
import { createStructuredSelector } from 'reselect';
import { connect } from 'react-redux';
import { teamsForUser } from 'selectors/teams';

const mapState = createStructuredSelector({
  orgsWithTeams: teamsForUser
});

@connect(mapState)
@autoaction({
  getTeamsForUser: (props) => {
    return {
      name: props.params.username
    };
  }
}, {
  getTeamsForUser
})
export default class UserTeams extends Component {

  static propTypes = {
    orgsWithTeams: PropTypes.array,
    params: PropTypes.object
  };

  makeTeamRows = () => {

    const {
      orgsWithTeams
    } = this.props;

    if (orgsWithTeams.length > 0) {
      // orgsWithTeams is an array of orgs
      // containing an array of teams
      return orgsWithTeams.map((org) => {
        return org.map((team) => {
          return (
            <tr>
              <td>
                { team.team.orgName }
              </td>
              <td>
                { team.team.name }
              </td>
              <td>
                {
                  team.isAdmin ? 'Admin' : 'Member'
                }
              </td>
            </tr>
          );
        });
      });
    }

    return [];
  };

  render () {
    return (
    <span styleName='teamsTable'>
        <PaginatedTable
          perPage={ 5 }
          headers={ ['Organization', 'Team', ''] }
          rows={ this.makeTeamRows() } />
      </span>
    );
  }
}
