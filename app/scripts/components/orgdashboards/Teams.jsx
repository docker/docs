'use strict';

import styles from './Teams.css';

import React, { Component, PropTypes } from 'react';
import { Link } from 'react-router';
import connectToStores from 'fluxible-addons-react/connectToStores';
import {
  Module,
  PageHeader,
  FontAwesome as FA
} from 'dux';

import _ from 'lodash';
import find from 'lodash/collection/find';

import createOrgTeamAction from 'actions/createOrgTeam';
import DashboardTeamsStore from 'stores/DashboardTeamsStore';
import { STATUS as ORGTEAMSTATUS } from 'stores/orgteamstore/Constants';
import Members from './Members';
import OrganizationAddTeam from '../account/orgs/OrganizationAddTeam';
import TeamProfile from '../account/orgs/TeamProfile';
import TeamListItem from 'common/TeamListItem';
import ListSelector from 'common/ListSelector';

const { string, array, object, func } = PropTypes;
var debug = require('debug')('DashboardTeams');

class TeamsList extends Component {
  static contextTypes = {
    executeAction: func.isRequired
  }

  static propTypes = {
    JWT: string,
    user: object,
    currentUserContext: string,
    teams: array,
    success: string,
    errorDetails: object
  }

  state = {
    addingTeam: false,
    updatingTeam: !!this.props.location.query.team,
    currentTeam: this.props.location.query
  }

  _onAddTeamClick = (e) => {
    this.setState({
      addingTeam: true,
      updatingTeam: false
    });
  }

  onDeleteTeam = () => {
    this.setState({
      updatingTeam: false
    });
  }

  _handleTeamClick = (team) => {
    this.setState({
      currentTeam: {team: team.name},
      addingTeam: false,
      updatingTeam: true
    });
    this.props.history.pushState(null, `/u/${this.props.currentUserContext}/dashboard/teams/`, {team: team.name});
  }

  _onCancel = () => {
    this.setState({
      addingTeam: false,
      updatingTeam: true
    });
  }

  _mkTeamListItem = (team) => {
    let selectedArrow;
    if (this.state.currentTeam.team === team.name) {
      selectedArrow = <span className={'right ' + styles.arrowSelect}><FA icon='fa-chevron-right' /></span>;
    }
    return (
      <li key={team.name}
              onClick={this._handleTeamClick.bind(null, team)}>
        <div className='list-item'>
          {team.name}
          {selectedArrow}
        </div>
      </li>
    );
  }

  componentDidMount() {
    const teamExists = find(this.props.teams, (team) => {
      return team.name === this.props.location.query.team;
    });
    if (!this.props.location.query.team) {
      //If there are teams, on load navigate to the first team and show members
      if (this.props.teams && this.props.teams.length > 0) {
        this.props.history.pushState(null, `/u/${this.props.currentUserContext}/dashboard/teams/`, {team: this.props.teams[0].name});
      }
    } else if (!teamExists) {
      //If the team query isn't one of the orgs teams navigate to first in list
      this.props.history.pushState(null, `/u/${this.props.currentUserContext}/dashboard/teams/`, {team: this.props.teams[0].name});
    }
  }

  render() {
    var addOrUpdateTeamForm;
    if (this.state.addingTeam && !this.props.teamReadOnly) {
      addOrUpdateTeamForm = (
        <Module>
          <h5 className={styles.formHeading}>Create Team</h5>
          <OrganizationAddTeam JWT={this.props.JWT}
                               user={this.props.user}
                               params={this.props.params}
                               onCancel={this._onCancel}/>
        </Module>
      );
    } else if (this.state.updatingTeam && !this.props.teamReadOnly) {
      let currentTeam = find(this.props.teams, (team) => {
        return this.state.currentTeam.team === team.name;
      }) || {};

      addOrUpdateTeamForm = (
        <Module>
          <h5 className={styles.formHeading}>Edit {this.props.currentUserContext}</h5>
          <TeamProfile JWT={this.props.JWT}
                       history={this.props.history}
                       orgname={this.props.currentUserContext}
                       success={this.props.success}
                       onDeleteTeam={this.onDeleteTeam}
                       onUpdateTeam={this.onUpdateTeam}
                       team={currentTeam}
                       error={this.props.errorDetails.detail} />
        </Module>
      );
    }

    var maybeAddTeam;

    if (!this.props.teamReadOnly) {
      maybeAddTeam = (
        <button className='button'
                onClick={this._onAddTeamClick}>Create Team <FA icon='fa-plus'/>
        </button>
      );
    }

    return (
      <div>
        <PageHeader title={`${this.props.currentUserContext}'s teams`}>
          {maybeAddTeam}
        </PageHeader>
        <div className='row'>
          <div className={'columns large-4 ' + styles.marginTop}>
            <ListSelector header={<h5>Choose Team</h5>} items={this.props.teams.map(this._mkTeamListItem)}>
            </ListSelector>
          </div>
          <div className={'columns large-4 ' + styles.marginTop}>
            <Members JWT={this.props.JWT}
              location={this.props.location}
              user={this.props.params.user}/>
          </div>
          <div className='columns large-4'>
            {addOrUpdateTeamForm}
          </div>
        </div>
      </div>
    );
  }
}

export default connectToStores(TeamsList,
                               [
                                 DashboardTeamsStore
                               ],
                               function({ getStore }, props) {
                                 return getStore(DashboardTeamsStore).getState();
                               });
