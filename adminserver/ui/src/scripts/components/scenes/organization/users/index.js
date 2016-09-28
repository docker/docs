'use strict';

import React, { Component, PropTypes } from 'react';
const { func, object, array, oneOfType, bool } = PropTypes;
import { connect } from 'react-redux';
import uiDecorator from 'redux-ui';
import consts from 'consts';
import { getOrgMembers, getIsCurrentUserPageOrgAdmin } from 'selectors/organizations';
import { getCurrentTeamMembers } from 'selectors/teams';
import { createStructuredSelector } from 'reselect';
// Actions
import {
  listOrgOrTeamMembers,
  deleteOrgOrTeamMember
} from 'actions/orgsAndTeams';
import autoaction from 'autoaction';
import { mapActions } from 'utils';
// Components
import UserList from './list';
import Button from 'components/common/button';
import AddUser from './addUser';
import NoUsers from './noUsers';
import Spinner from 'components/common/spinner';
// CSS
import styles from './users.css';
import css from 'react-css-modules';

const mapState = createStructuredSelector({
  members: getOrgMembers,
  teamMembers: getCurrentTeamMembers,
  isOrgAdmin: getIsCurrentUserPageOrgAdmin
});

/**
 * OrgUserList shows a list of users for the currently selected organization or
 * team.
 *
 */
@connect(
  mapState,
  mapActions({
    deleteOrgOrTeamMember
  })
)
@autoaction({
  listOrgOrTeamMembers: (props) => {
    return {
      orgName: props.params.org,
      teamName: props.params.team || '',
      limit: 100
    };
  }
}, {
  listOrgOrTeamMembers
})
@uiDecorator({
  state: {
    isAddVisible: false
  }
})
@css(styles)
export default class OrgUserList extends Component {

  static propTypes = {
    actions: object,
    ui: object,
    updateUI: func,
    members: oneOfType([object, array]),
    teamMembers: oneOfType([object, array]),
    params: object,
    isOrgAdmin: bool
  }

  addUser() {
    this.props.updateUI('isAddVisible', true);
  }

  cancelAddUser(evt) {
    evt.preventDefault();
    this.props.updateUI('isAddVisible', false);
  }

  deleteUser = (user) => {
    const {
      actions,
      params: {
        org,
        team
      }
    } = this.props;
    actions.deleteOrgOrTeamMember({
      orgName: org,
      teamName: team,
      memberName: user.member.name
    });
  };

  render() {
    const {
      members,
      teamMembers,
      ui,
      params: {
        org,
        team
      },
      isOrgAdmin
    } = this.props;

    // decide which list of users to use
    const memberList = team ? teamMembers : members;

    // When either of these two finishes the content will load.  Because we only
    // make one of the calls the other will be undefined; our spinner assumes
    // undefined means that the call was never made and treats undefined as
    // a success case.
    const status = [
      [consts.organizations.LIST_ORGANIZATION_MEMBERS, org],
      [consts.teams.LIST_MEMBERS, org, team]
    ];

    return (
      <div styleName='wrapper'>
        {
          isOrgAdmin ?

          <span>
            <div styleName='actions'>
              <Button
                id='add-user-button'
                disabled={ ui.isAddVisible }
                variant='secondary'
                onClick={ ::this.addUser }>
                Add user
              </Button>
            </div>

            {
              ui.isAddVisible ?
              <AddUser
                params={ this.props.params }
                cancel={ ::this.cancelAddUser }
              />
              : undefined
            }
          </span>

          :

          undefined
        }

        <Spinner loadingStatus={ status } >
          { Object.keys(memberList).length > 0
              ? <UserList
                  params={ this.props.params }
                  onDelete={ ::this.deleteUser }
                  canEdit={ isOrgAdmin }
                  members={ memberList } />
              : <NoUsers noun={ team === undefined ? 'organization' : 'team' } /> }
        </Spinner>
      </div>
    );
  }
}
