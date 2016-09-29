'use strict';
import React, {Component, PropTypes} from 'react';
import DUXInput from 'common/DUXInput.jsx';
import createOrgTeamAction from '../../../actions/createOrgTeam';
import OrgTeamStore from '../../../stores/OrgTeamStore';
import connectToStores from 'fluxible-addons-react/connectToStores';
import Button from '@dux/element-button';
import styles from './TeamProfile.css';

import { STATUS as ORG_STATUS } from '../../../stores/orgteamstore/Constants';
const { string, array, func, object } = PropTypes;

class OrganizationAddTeam extends Component {
  static contextTypes = {
    executeAction: func.isRequired
  }

  static propTypes = {
    JWT: string.isRequired,
    onCancel: func,
    name: string,
    description: string,
    members: array,
    errorDetails: object,
    success: string,
    STATUS: string
  }

  state = {
    teamname: '',
    teamdesc: '',
    members: []
  }

  _handleCreateTeam = (e) => {
    e.preventDefault();
      this.context.executeAction(createOrgTeamAction,
        {
          jwt: this.props.JWT,
          orgName: this.props.params.user,
          team: {name: this.state.teamname, description: this.state.teamdesc} });
  }

  _handleReset = (e) => {
    e.preventDefault();
  }

  teamNameChange = (e) => {
    this.setState({teamname: e.target.value});
  }

  teamDescChange = (e) => {
    this.setState({teamdesc: e.target.value});
  }

  render() {
    let maybeError;
    if (this.props.errorDetails.detail) {
      maybeError = <span className='alert-box alert radius'>{this.props.errorDetails.detail}</span>;
    }
    return (
      <div className="dux-form add-team-container">
        <form onSubmit={this._handleCreateTeam}>
          <DUXInput type="text" styleType="light"
                    label="Team Name" onChange={this.teamNameChange} value={this.state.teamname}/>
          <DUXInput type="text" styleType="light"
                    label="Description" onChange={this.teamDescChange} value={this.state.teamdesc}/>
          <div className={styles.addTeamButtonGroup}>
            <Button type="submit" size='small'>Add</Button>
            <Button variant='alert' size='small' onClick={this.props.onCancel}>Cancel</Button>
          </div>
          {maybeError}
        </form>
      </div>

    );
  }
}

export default connectToStores(OrganizationAddTeam,
                               [
                                 OrgTeamStore
                               ],
                               function({ getStore }, props) {
                                 return getStore(OrgTeamStore).getState();
                               });
