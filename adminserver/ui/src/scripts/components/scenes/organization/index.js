'use strict';

import React, { Component, PropTypes } from 'react';
const { array, bool, func, instanceOf, node, object, string } = PropTypes;
import { OrganizationRecord } from 'records';
import { connect } from 'react-redux';
// Actions
import autoaction from 'autoaction';
import { getOrganization } from 'actions/organizations';
import { listTeams } from 'actions/teams';
import {
  getAuthSettings
} from 'actions/settings';
// Components
import TeamList from './teamList.js';
import AddTeamForm from './addTeamForm.js';
import { Tabs, Tab } from 'components/common/tabs';
import { Link } from 'react-router';
import FA from 'components/common/fontAwesome';
// Selectors
import { createStructuredSelector } from 'reselect';
import { isAdminSelector } from 'selectors/users';
import {
  getPageOrg,
  getIsCurrentUserPageOrgAdmin
} from 'selectors/organizations';
import { getAuthMethod } from 'selectors/settings';
import { getTeamsWithoutOwners } from './selectors.js';
// Misc
import consts from 'consts';
import css from 'react-css-modules';
import styles from './organization.css';
import ui from 'redux-ui';

const mapState = createStructuredSelector({
  org: getPageOrg,
  teams: getTeamsWithoutOwners,
  // For whether we can add a team or not.
  isOrgAdmin: getIsCurrentUserPageOrgAdmin,
  authMethod: getAuthMethod
});

/**
 * This represents the main parent component for listing an organization.
 * Within this we render the organization's teams, users, repositories and
 * settings.
 *
 * Because org teams are shown on every page within the left panel this must
 * request the orgs teams directly.
 *
 */
@ui({
  state: {
    showAddTeam: false
  },
  reducer: (state, action) => {
    if (action.type === consts.teams.CREATE_TEAM && action.ready && !action.error) {
      return state.set('showAddTeam', false);
    }
    return state;
  }
})
@connect(mapState)
@autoaction({
  getAuthSettings: [],
  getOrganization: (props) => props.params.org,
  // We need to know whether the user is an org admin for configuring repos,
  // teams etc.
  listTeams: (props) => ({ orgName: props.params.org, limit: 50 })
}, { getOrganization, listTeams, getAuthSettings })
@css(styles)
export default class Organization extends Component {

  static propTypes = {
    updateUI: func,
    ui: object,
    params: object,

    org: instanceOf(OrganizationRecord),
    teams: array,
    // We can only create teams if we're an org admin
    isOrgAdmin: bool,

    authMethod: PropTypes.string,

    // This represents the content we're going to show when a tab is selected.
    // The children prop is passed via the router.
    children: node
  }

  showAddTeam() {
    this.props.updateUI('showAddTeam', true);
  }

  hideAddTeam() {
    this.props.updateUI('showAddTeam', false);
  }

  render() {
    const {
      org,
      isOrgAdmin,
      teams,
      children,
      params,
      authMethod,
      ui: { showAddTeam }
    } = this.props;

    return (
      <div styleName='wrapper'>

            <div styleName='left'>
              <div styleName='usersIcon'>
                <FA icon='fa-users' />
              </div>
              <h2 styleName='orgName'>{ org.name }</h2>
              <div styleName='teamHeader'>
                <Link to={ `/orgs/${org.name}` }
                  activeClassName={ (params.team === undefined) ? styles.active : undefined }>Everyone</Link>
                <h3 styleName='teamHeading'>Teams</h3>
                <div styleName='add'>
                  { isOrgAdmin
                      ? <FA id='add-team-button' icon='fa-plus' onClick={ ::this.showAddTeam } />
                      : undefined }
                </div>
                { showAddTeam
                    ? (
                        <AddTeamForm
                          params={ this.props.params }
                          onHide={ ::this.hideAddTeam }
                          isLdapEnabled={ authMethod === 'ldap' } />
                      )
                    : undefined }
              </div>
              {
                teams.length ?
                  <TeamList
                    teams={ teams }
                    orgName={ org.name } />
                :
                isOrgAdmin && !showAddTeam &&
                  <span>
                    <p styleName='createTeam'>Create a team to give users more repository permissions.
                    <FA icon='fa-long-arrow-right'/></p>
                  </span>
              }
            </div>
            <div styleName='content'>
              { params.team
                  ? <TeamTabsHeader orgName={ params.org } teamName={ params.team } />
                  : <OrgTabsHeader orgName={ params.org } /> }
              { children }
            </div>

      </div>
    );
  }
}

class TeamTabsHeader extends Component {
  static propTypes = {
    orgName: string,
    teamName: string
  }

  render() {
    const {
      orgName,
      teamName
    } = this.props;

    return (
      <Tabs header>
        <Tab>
          <Link to={ `/orgs/${ orgName }/teams/${ teamName }/users` }>Members</Link>
        </Tab>
        <Tab>
          <Link to={ `/orgs/${ orgName }/teams/${ teamName }/repos` }>Repositories</Link>
        </Tab>
        <Tab id='team-settings-tab'>
          <Link to={ `/orgs/${ orgName }/teams/${ teamName }/settings` }>Settings</Link>
        </Tab>
      </Tabs>
    );
  }
}

@connect(createStructuredSelector({
  isSysAdmin: isAdminSelector
}))
class OrgTabsHeader extends Component {
  static propTypes = {
    orgName: string,
    isSysAdmin: bool
  }

  maybeRenderSettingsTab(orgName) {
    if (!this.props.isSysAdmin) {
      return null;
    }

    return (
      <Tab id='org-settings-tab'>
        <Link to={ `/orgs/${ orgName }/settings` }>Settings</Link>
      </Tab>
    );
  }

  render() {
    const { orgName } = this.props;

    return (
      <Tabs header>
        <Tab>
          <Link to={ `/orgs/${ orgName }/users` }>Members</Link>
        </Tab>
        <Tab>
          <Link to={ `/orgs/${ orgName }/repos` }>Repositories</Link>
        </Tab>
        { this.maybeRenderSettingsTab(orgName) }
      </Tabs>
    );
  }
}
