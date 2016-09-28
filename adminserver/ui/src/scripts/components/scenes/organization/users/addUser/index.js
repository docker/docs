'use strict';

import React, { Component, PropTypes } from 'react';
const { object, func, array, oneOfType } = PropTypes;
import css from 'react-css-modules';
import styles from './addUser.css';
import RadioGroup from 'components/common/radioGroup';
import Input from 'components/common/input';
import AddUserForm from 'components/scenes/users/addUserForm';
import FontAwesome from 'components/common/fontAwesome';
import ui from 'redux-ui';
import autoaction from 'autoaction';
import { listUsers, searchUsers } from 'actions/users';
import { usersSelector } from 'selectors/users';
import { getOrgMembers } from 'selectors/organizations';
import { createStructuredSelector } from 'reselect';
import { connect } from 'react-redux';
import Checkbox from 'components/common/checkbox';
import Button from 'components/common/button';
import {
  listOrgOrTeamMembers,
  deleteOrgOrTeamMember
} from 'actions/orgsAndTeams';
import { addTeamMembers } from 'actions/teams';
import { addOrganizationMember, createUserAndAddToOrg } from 'actions/organizations';
import { mapActions } from 'utils';
import { reset } from 'redux-form';
import { getCurrentTeamMembers } from 'selectors/teams';

const mapState = createStructuredSelector({
  users: usersSelector,
  orgMembers: getOrgMembers,
  teamMembers: getCurrentTeamMembers
});

@connect(
  mapState,
  mapActions({
    addOrganizationMember,
    deleteOrgOrTeamMember,
    createUserAndAddToOrg,
    reset,
    addTeamMembers,
    searchUsers
  })
)
@autoaction({
  listOrgOrTeamMembers: (props) => {
    return {
      orgName: props.params.org,
      teamName: props.params.team || '',
      limit: 100
    };
  },
  listUsers: []
}, {
  listUsers,
  listOrgOrTeamMembers
})
@ui({
  state: {
    choices: [
      {
        label: 'Existing',
        value: 'existing'
      },
      {
        label: 'New',
        value: 'new'
      }
    ],
    selected: 'existing',
    // combine orgMembers with manually selected users
    selectedUsers: (props) => {
      if (!props.params.team) {
        return Object.keys(props.orgMembers) || [];
      }
      return Object.keys(props.teamMembers) || [];
    },
    usersToRemove: [],
    createAnotherUser: false
  }
})
@css(styles, { allowMultiple: true })
export default class AddUser extends Component {

  static propTypes = {
    ui: object,
    updateUI: func,
    users: object,
    orgMembers: oneOfType([array, object]),
    teamMembers: oneOfType([array, object]),
    actions: object,
    params: object,
    cancel: func
  };

  // called when Save & Add Another is clicked
  shouldCreateAnotherUser = () => {
    this.props.updateUI({ createAnotherUser: true });
  }

  chooseScene = (sceneChoice) => {
    this.props.updateUI({
      selected: sceneChoice
    });
  };

  searchUsers = (evt) => {
    this.props.actions.searchUsers(evt.target.value);
  };

  onSubmit = (data) => {
    data.name = data.username;
    data.type = 'user';

    const { org, team } = this.props.params;

    this.props.actions.createUserAndAddToOrg(
      data,
      org,
      data.username,
      team
    );

    this.props.actions.reset('newUser');

    if (this.props.ui.createAnotherUser) {
      this.props.updateUI({
        createAnotherUser: false
      });
    } else {
      this.props.updateUI({
        isAddVisible: false
      });
    }

  }

  selectUser = (user) => () => {

    const {
      ui: {
        selectedUsers
      },
      updateUI
    } = this.props;

    if (selectedUsers.includes(user)) {
      this.deselectUser(user)();
    } else {
      updateUI({
        selectedUsers: selectedUsers.concat([user])
      });
    }

  };

  deselectUser = (user) => () => {

    const {
      updateUI,
      orgMembers,
      teamMembers,
      ui: {
        selectedUsers,
        usersToRemove
      }
    } = this.props;

    updateUI({
      selectedUsers: selectedUsers.filter((username) => {
        return username !== user;
      })
    });

    if (this.refs.user) {
      // uncheck the related box
      this.refs.user.checked = false;
    }

    // which list to use?
    const checkList = this.props.params.team ? teamMembers : orgMembers;
    // check if the user was already in the org/team
    // if so, we have to add them to the list of users to remove

    // I'm not really happy with how I'm checking team/org all over the place
    // TODO figure out a better way to check once
    if (Object.keys(checkList).includes(user) && !usersToRemove.includes(user)) {
      updateUI({
        usersToRemove: usersToRemove.concat([user])
      });
    }
  };

  syncUsers = () => {

    const {
      ui: {
        selectedUsers,
        usersToRemove
      },
      orgMembers,
      teamMembers,
      params: {
        org,
        team
      }
    } = this.props;

    if (!team) {
      selectedUsers
      .filter((username) => !Object.keys(orgMembers).includes(username))
      .map((username) => {
        this.props.actions.addOrganizationMember({
          name: org,
          member: username
        });
      });
    } else {
      this.props.actions.addTeamMembers(
        org,
        team,
        selectedUsers.filter((username) => !Object.keys(teamMembers).includes(username))
      );
    }

    usersToRemove.map((username) => {
      this.props.actions.deleteOrgOrTeamMember({
        orgName: org,
        teamName: team,
        memberName: username
      });
    });

    this.props.updateUI({
      isAddVisible: false
    });

  };

  render() {

    const {
      ui: {
        choices,
        selected,
        selectedUsers
      }
    } = this.props;

    return (
      <div styleName='formbox' id='add-user-ui'>
        <RadioGroup
          id='user-create-toggle'
          initialChoice={ selected }
          choices={ choices }
          onChange={ ::this.chooseScene } />
        {
          selected === 'existing' ?
          <span>
            <div styleName='row'>
              <div styleName='existingUsersColumn'>
                <div styleName='searchInput'>
                  <Input
                    onChange={ ::this.searchUsers }
                    type='text'
                    placeholder='Search users'
                  />
                </div>
                <div styleName='userRows'>
                  {
                    Object.keys(this.props.users).map((user, i) => {
                      return (
                          <div key={ i }>
                            <Checkbox
                              ref={ user }
                              isChecked={ selectedUsers.includes(user) ? true : false }
                              styleName='check'
                              onChange={ ::this.selectUser(user) }
                            />
                            <FontAwesome icon='fa-user' />
                            { user }
                          </div>
                      );
                    })
                  }
                </div>
              </div>
              <div styleName='selectedUsersColumn'>
                { selectedUsers.length } Selected
                { selectedUsers.map((user, i) => {
                  return (
                    <div key={ i }>
                      <FontAwesome icon='fa-user' />
                      { user }
                      <FontAwesome icon='fa-close' onClick={ ::this.deselectUser(user) } />
                    </div>
                  );
                }) }
              </div>
            </div>
            <div styleName='row buttonRow'>
              <Button variant='primary simple' onClick={ this.props.cancel }>Cancel</Button>
              <Button variant='primary' onClick={ ::this.syncUsers }>Save</Button>
            </div>
          </span>
          :
          <AddUserForm
            submitHandler={ ::this.onSubmit }
            addAnotherUserHandler={ ::this.shouldCreateAnotherUser }
            cancel={ this.props.cancel }
          />
        }
      </div>
    );
  }
}
