'use strict';

import styles from './Members.css';
import React, { createClass } from 'react';
import { Link } from 'react-router';
import connectToStores from 'fluxible-addons-react/connectToStores';
import classnames from 'classnames';
import DashboardMembersStore from '../../stores/DashboardMembersStore';
import _ from 'lodash';
import clearMemberErrorAction from '../../actions/clearMemberError';
import createTeamMembersAction from '../../actions/createTeamMembers';
import removeTeamMemberAction from '../../actions/removeTeamMember';
import ListSelector from '../common/ListSelector';
import FA from 'common/FontAwesome';


var debug = require('debug')('DashboardTeams');

function mkMemberListItem(member) {
  const removeButtonStyles = classnames({
    'fa fa-remove right': true,
    [styles.removeButton]: true
  });
  return (
    <li key={member.username}
        className="li-no-hover">
      <div className="delete-member-item">
        <Link to={`/u/${member.username}/`}>
        {member.username}
        </Link>
        <i className={removeButtonStyles}
           onClick={this._removeMember.bind(null, member.username)}></i>
      </div>
    </li>
  );
}

var MembersList = createClass({
  displayName: 'MembersList',
  propTypes: {
    members: React.PropTypes.array,
    JWT: React.PropTypes.string,
    error: React.PropTypes.object,
    location: React.PropTypes.object,
    user: React.PropTypes.string
  },
  contextTypes: {
    executeAction: React.PropTypes.func.isRequired
  },
  getInitialState() {
    return {
      member: ''
    };
  },
  onMemberChange(e) {
    if (this.props.error.response) {
      this.context.executeAction(clearMemberErrorAction);
    }
    this.setState({
      member: e.target.value
    });

  },
  _addMember(e) {
    e.preventDefault();
    //username can be either an email or username
    var membersToAdd = [{username: this.state.member}];
    this.context.executeAction(createTeamMembersAction,
      {
        jwt: this.props.JWT,
        orgname: this.props.user,
        teamname: this.props.location.query.team,
        members: membersToAdd
      });
    this.setState({
      member: ''
    });
  },
  _removeMember(username) {
    this.context.executeAction(removeTeamMemberAction,
      {
        jwt: this.props.JWT,
        orgname: this.props.user,
        teamname: this.props.location.query.team,
        membername: username
      });
  },
  render() {
    var maybeError;
    var err = this.props.error.response;
    var errorText;
    if (err) {
      if (err.unauthorized || err.forbidden) {
        errorText = 'You have no permissions to perform this action.';
      } else if (err.badRequest || err.clientError) {

        if(err.body) {
          errorText = err.body.detail || err.body.__all__;
        } else {
          errorText = 'Error adding member. Please check that this email or username is associated with a Docker Hub account.';
        }
      } else if (err.serverError) {
        errorText = 'An error occurred while adding the member. Please try again shortly.';
      }
      maybeError = (
        <li key="listError">
          <span className='alert-box alert radius'>{errorText}</span>
        </li>
      );
    }
    var headerItem = (
      <form className='row collapse'
                             onSubmit={this._addMember}>
        <div className={'small-10 columns ' + styles.addMemberInput}>
          <input type='text'
                 placeholder='Add new member by username or email'
                 value={this.state.member}
                 onChange={this.onMemberChange}></input>
        </div>
        <div className='small-2 columns'>
          <button className={'button postfix ' + styles.removePadding}>
            <FA icon='fa-plus'
               onClick={this._addMember} />
          </button>
        </div>
      </form>
    );
    var listItems = this.props.members.map(mkMemberListItem, this);
    if (maybeError) {
      listItems.push(maybeError);
    }
    return (
      <ListSelector header={headerItem} items={listItems} />
    );
  }
});

export default connectToStores(MembersList,
  [
    DashboardMembersStore
  ],
  function({ getStore }, props) {
    return getStore(DashboardMembersStore).getState();
  });
