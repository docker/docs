'use strict';

import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { listRepoTeamAccess } from 'actions/repositories';
import autoaction from 'autoaction';
import PaginatedTable from 'components/common/paginatedTable';
import css from 'react-css-modules';
import styles from './repository.css';
import TeamRow from 'components/common/teamRow';
import { teamRepositoryAccessSelector } from 'selectors/repositories';


import { createStructuredSelector } from 'reselect';

const mapRepositoryState = createStructuredSelector({
    teams: teamRepositoryAccessSelector
});
@connect(mapRepositoryState)
@autoaction({
  listRepoTeamAccess: (props) => ({
    namespace: props.params.namespace,
    repo: props.params.name
  })
}, {
  listRepoTeamAccess
})
@css(styles)
export class RepositoryTeamsTab extends Component {
  static propTypes = {
    params: PropTypes.object.isRequired,
    teams: PropTypes.object
  }

  makeTeamRows = () => {

    const {
      teams,
      params: {
        namespace,
        name
      }
    } = this.props;

    if (teams && Object.keys(teams).length) {
      return Object.keys(teams).map((teamId, i) => {
        return (
          <TeamRow
            repo={ name }
            permission={ teams[teamId].accessLevel }
            key={ i }
            id={ teamId }
            org={ namespace } />
        );
      });
    }
    return [];
  };

  render() {

    return (
      <span styleName='teamsTable'>
        <PaginatedTable
          perPage={ 5 }
          headers={ ['Team', 'Permission', ''] }
          rows={ this.makeTeamRows() } />
      </span>
    );

  }
}
