'use strict';

import React, {PropTypes, Component} from 'react';
import saveTeamProfileAction from '../../../actions/saveTeamProfile';
import deleteTeamProfileAction from '../../../actions/removeTeam';
import DUXInput from 'common/DUXInput.jsx';
import styles from './TeamProfile.css';
import Button from '@dux/element-button';
import _ from 'lodash';

const { string, object, bool, func } = PropTypes;

export default class TeamProfile extends Component {

  static contextTypes = {
    executeAction: func.isRequired
  }

  static propTypes = {
    JWT: string,
    history: object.isRequired,
    error: string,
    success: string,
    orgname: string,
    clearError: func,
    onDeleteTeam: func,
    onUpdateTeam: func,
    team: object
  }

  state = {
    teamName: this.props.team.name,
    teamDesc: this.props.team.description
  }

  onSubmit = (e) => {
    e.preventDefault();
    var updatedTeam = {
      name: this.state.teamName,
      description: this.state.teamDesc
    };
    this.context.executeAction(saveTeamProfileAction, {
      jwt: this.props.JWT,
      orgname: this.props.orgname,
      teamname: this.props.team.name,
      team: updatedTeam
    });
  }

  onDelete = (e) => {
    e.preventDefault();
    this.context.executeAction(deleteTeamProfileAction, {
      jwt: this.props.JWT,
      orgname: this.props.orgname,
      teamname: this.props.team.name
    });
  }

  teamNameChange = (e) => {
    if (this.state.teamName !== 'owners') {
      this.setState({teamName: e.target.value});
    } else {
      this.setState({teamName: 'owners'});
    }
  }

  teamDescChange = (e) => {
    this.setState({teamDesc: e.target.value});
  }

  render() {
    var maybeError = <span />;
    var maybeSuccess = <span />;
    if (this.props.error) {
      maybeError = <span className='alert-box alert radius'>{this.props.error}</span>;
    } else if(this.props.success) {
      //TODO: this could be an alert box with a close icon that can be closed when the user wants to dismiss
      //Or time out after some time (ideally i would like to see this as a notification and that's it)
      maybeSuccess = <div><br /><span className='alert-box success radius'>{this.props.success}</span></div>;
    }
    var maybeDeleteBtn;
    if (this.state.teamName !== 'owners' && this.props.team.name === this.state.teamName) {
      maybeDeleteBtn = (<Button variant='alert' size='small' onClick={this.onDelete}>Delete</Button>);
    }
    return (
      <form onSubmit={this.onSubmit}>
        <DUXInput type="text" label="Team Name" styleType="light"
                  value={this.state.teamName}
                  onChange={this.teamNameChange} />
        <DUXInput type="text" label="Description" styleType="light"
                  value={this.state.teamDesc}
                  onChange={this.teamDescChange} />
        <div className={styles.addTeamButtonGroup}>
          <Button type="submit" size='small'>Save</Button>
          {maybeDeleteBtn}
        </div>
        {maybeError}
        {maybeSuccess}
      </form>
    );
  }
}
