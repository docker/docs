'use strict';

import React, { Component, PropTypes } from 'react';
const { array, bool, object } = PropTypes;
import { connect } from 'react-redux';

import consts from 'consts';
import { mapActions } from 'utils';
import * as teamActions from 'actions/teams';
import * as organizationActions from 'actions/organizations';
import MemberList from './memberList';
import {
  Spinner
} from 'components/common';
import autoaction from 'autoaction';
import { createStructuredSelector } from 'reselect';
import { orgMembersSelector } from 'selectors/organizations';
import styles from './resource.css';

const mapState = createStructuredSelector({
  orgMembers: orgMembersSelector
});

@connect(mapState, mapActions({...organizationActions, ...teamActions}))
@autoaction({
  listOrganizationMembers: (params, state) => state.router.params.org
}, organizationActions)
export default class OrganizationMembers extends Component {

  static propTypes = {
    actions: object,
    // This must be given the organization in order to query for teams.  It
    // does *not* use the org reducer's `organization` resource.
    organization: object, // passed from parent
    isAdminOrOrgOwner: bool, // passed from parent
    orgMembers: array.isRequired
  }

  deleteMemberFromAllTeams(memberName) {
    const orgName = this.props.organization.name;
    this.props.actions.deleteOrganizationMember(orgName, memberName);
  }

  renderZeroState() {
    if (!this.props.orgMembers || this.props.orgMembers.length) {
      return null;
    }
    // Note: Only admins can see this message
    // Org owners and org members will bump the orgMembers list to nonzero
    // Nonorg members won't be able to see this tab
    return (
      <div className={ styles.zero }>
        This organization has no members...
        add some by adding people to a team of this organization!
      </div>
    );
  }

  renderMemberList() {
    const { organization, orgMembers } = this.props;
    if (orgMembers.length > 0) {
      return (
        <MemberList
          members={ orgMembers }
          canDelete={ this.props.isAdminOrOrgOwner && window.authMethod === 'managed' }
          onDelete={ (member) => this.deleteMemberFromAllTeams(member.name) }
          onDeleteTooltip={ (member) => <span>Remove <code>{ member.name }</code> from all teams of <code>{ organization.name }</code></span> } />
      );
    }
  }

  render() {
    const org = this.props.organization.name;
    const listOrgMembersStatus = [consts.organizations.LIST_ORGANIZATION_MEMBERS, org];
    return (
      <Spinner loadingStatus={ [listOrgMembersStatus] } >
        { this.renderZeroState() }
        { this.renderMemberList() }
      </Spinner>
    );
  }
}
