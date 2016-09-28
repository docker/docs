'use strict';

import React, { Children, Component, PropTypes, cloneElement } from 'react';
const { bool, func, instanceOf, node, object, shape, string } = PropTypes;
import { pushState, replaceState } from 'redux-router';
import * as TeamActions from 'actions/teams';
import * as OrgActions from 'actions/organizations';
import * as StatusActions from 'actions/status';
import { connect } from 'react-redux';
import { reduxForm } from 'redux-form';
import Immutable from 'immutable';
import consts from 'consts';
import Button from 'components/common/button';
import Spinner from 'components/common/spinner';
import createTeamValidator from 'components/scenes/organizations/resource/createTeam/createTeamValidator';
import styles from './team.css';
import {
  ActionBar,
  BreadCrumbs,
  ErrorMessage,
  Input,
  Radio,
  RadioGroup,
  Tabs,
  Tab,
  TextArea,
  DeleteModal
} from 'components/common';
import { Link } from 'react-router';
import { connectModal } from 'components/common/modal';

import { fieldsWithErrorOnlyIfTouched, mapActions } from 'utils';

// Exports from other local components
export { TeamRepositories } from './repositories.js';
export TeamMembers from './members.js';

import { createStructuredSelector } from 'reselect';
import { isAdminOrOrgOwnerSelector, teamSelector } from 'selectors/teams';
import { locationSelector } from 'selectors/router';
const mapState = createStructuredSelector({
  isAdminOrOrgOwner: isAdminOrOrgOwnerSelector,
  location: locationSelector,
  status: (state) => state.status,
  team: teamSelector
});

/**
 * What data does this team component need?
 *
 *  - Organization in which this team belongs
 *  - Repo list
 *  - Team info
 *  - Member list (summary)
 */
@connect(mapState, mapActions({
  teams: TeamActions,
  orgs: OrgActions,
  router: {pushState, replaceState},
  status: StatusActions
}))
export default class Team extends Component {
  static propTypes = {
    actions: object.isRequired,
    children: node.isRequired,
    params: object.isRequired,

    isAdminOrOrgOwner: bool,
    location: object,
    status: instanceOf(Immutable.Map).isRequired,
    team: object
  }

  componentWillMount() {
    // We must fetch the core team information as soon as this component loads
    const { team, org } = this.props.params;
    this.props.actions.status.resetStatus();
    this.props.actions.orgs.getOrganization(org);
    this.props.actions.teams.getTeam(org, team);

    if (org === '_global') {
      this.props.actions.router.replaceState({}, `/orgs/_global/teams/${this.props.params.team}/members`);
      return;
    }
  }

  /**
   * When the form subcomponent calls the delete action we need to detect the
   * 'success' status to redirect to the org teams page.
   *
   * We do this in this particular component because this receives the team name
   * from the URL parameters; once the team is deleted this is the only place we
   * can find the team name for checking delete status (as the team is no longer
   * in our state and passed via @connect once deleted).
   *
   */
  componentWillReceiveProps(next) {
    const { team, org } = this.props.params;
    const deleteTeamStatus = next.status.getIn([consts.teams.DELETE_TEAM, org, team, 'status']);
    if (deleteTeamStatus === consts.loading.SUCCESS) {
      this.props.actions.router.pushState({}, `/orgs/${org}/teams`);
    }

    const pathName = this.props.location.pathname;
    const getOwnershipStatus = this.props.status.getIn([consts.teams.GET_TEAM_MEMBER, org, 'owners', window.user.name, 'status']);
    if (getOwnershipStatus && getOwnershipStatus !== consts.loading.PENDING &&
        !next.isAdminOrOrgOwner &&
        next.team &&
        !next.team.clientUserIsMember &&
        pathName.indexOf('repositories') === pathName.length - 'repositories'.length) {
      // checking owners and getting team is done loading.
      // We are on team repos tab but current user isn't an admin, org owner, or team member
      // Hide tab and redirect to members tab
      this.props.actions.router.replaceState({}, `/orgs/${org}/teams/${team}/members`);
    }
  }

  render() {
    const {
      children,
      isAdminOrOrgOwner,
      team,
      params: {
        team: teamName,
        org: orgName
      }
    } = this.props;
    let isMember = false;
    if (team) {
      isMember = team.clientUserIsMember;
    }

    const getOrgStatus = [consts.organizations.GET_ORGANIZATION, orgName];
    const getTeamStatus = [consts.teams.GET_TEAM, orgName, teamName];
    const getOwnershipStatus = [consts.teams.GET_TEAM_MEMBER, orgName, 'owners', window.user.name];

    return (
      <div>
        <BreadCrumbs header items={ [
          ['Organizations', '/orgs'],
          [orgName, `/orgs/${orgName}/teams`],
          teamName
        ] } />
        <div className={ styles.container }>
          <Spinner loadingStatus={ [getOrgStatus, getTeamStatus, getOwnershipStatus] }>
            <div className={ styles.teamWrapper }>
              <div className={ styles.teamProfile }>
                <TeamInfo
                  canDelete={ teamName !== 'owners' && orgName !== '_global' }
                  canEdit={ isAdminOrOrgOwner }
                  orgName={ orgName }
                  team={ team || { name: teamName } }
                />
              </div>
              <div className={ styles.teamTabs }>
                <Tabs sidebar>
                  { (isAdminOrOrgOwner || isMember) && teamName !== 'owners' && orgName !== '_global' &&
                    // Owners team has admin access to everything so this tab isn't needed
                    <Tab><Link to={ `/orgs/${orgName}/teams/${teamName}/repositories` }>Repositories</Link></Tab>
                  }
                  <Tab><Link to={ `/orgs/${orgName}/teams/${teamName}/members` }>Members</Link></Tab>
                </Tabs>
                {
                  Children.map(
                    children,
                    child => cloneElement(child, {
                      isAdminOrOrgOwner,
                      team
                    })
                  )
                }
              </div>
            </div>
          </Spinner>
        </div>
      </div>
    );
  }
}

@connectModal()
@connect((state) => ({ status: state.status }), mapActions(TeamActions))
class TeamInfo extends Component {

  static propTypes = {
    actions: shape({
      deleteTeam: func.isRequired
    }),
    canDelete: bool,
    canEdit: bool,
    orgName: string.isRequired,
    team: shape({
      name: string.isRequired,
      description: string
    }).isRequired,

    showModal: func,
    hideModal: func
  }

  state = {
    isEditVisible: false
  }

  componentWillReceiveProps(next) {
    const { orgName, team: { name: teamName } } = this.props;
    const status = next.status.getIn([consts.teams.UPDATE_TEAM, orgName, teamName, 'status']);
    if (status === consts.loading.SUCCESS) {
      this.closeEditForm();
    }
  }

  showEditForm() {
    this.setState({ isEditVisible: true });
  }

  closeEditForm() {
    this.setState({ isEditVisible: false });
  }

  deleteTeam() {
    const { orgName, team: {name: teamName} } = this.props;
    this.props.actions.deleteTeam({ orgName, teamName });
    this.props.hideModal();
  }

  confirmDelete() {
    this.props.showModal(
      (
        <DeleteModal
          resourceType='team'
          resourceName={ this.props.team.name }
          onDelete={ ::this.deleteTeam }
          hideModal={ this.props.hideModal } />
      )
    );
  }


  /**
   * This is called when the EditTeamForm form is submitted
   */
  handleSubmit(data) {
    const { orgName, team } = this.props;
    return this.props.actions.updateTeam(orgName, { name: team.name, ...data });
  }

  render() {
    const { canDelete, canEdit, orgName, team } = this.props;
    const { isEditVisible } = this.state;

    let bottom;
    if (canEdit) {
      bottom = (
        <ActionBar>
          <Button onClick={ ::this.showEditForm } ghost>Edit</Button>
        </ActionBar>
      );
    }

    if (isEditVisible) {
      bottom = (
        <EditTeamForm
          canDelete={ canDelete }
          closeForm={ ::this.closeEditForm }
          onDelete={ ::this.confirmDelete }
          initialValues={ team.toJS() }
          formKey={ `editTeamForm/${orgName}/${team.id}` }
          onSubmit={ ::this.handleSubmit } />
      );
    }

    return (
      <div>
        <div>
          <div className={ styles.teamNameClass }>
            <h1 style={ { wordBreak: 'break-word' } }>{ team.name }</h1>
          </div>
          { !this.state.isEditVisible && <TeamDetails team={ team } /> }
        </div>
        { bottom }
      </div>
    );
  }

}

@reduxForm({
  form: 'editTeamForm',
  fields: [
    'description',
    'type',
    'ldapDN',
    'ldapGroupMemberAttribute'
  ],
  validate: createTeamValidator // TODO separate out the validations
})
class EditTeamForm extends Component {

  static propTypes = {
    canDelete: bool,
    closeForm: func.isRequired,
    error: string,
    fields: object.isRequired,
    handleSubmit: func.isRequired,
    invalid: bool,
    onDelete: func.isRequired,
    pristine: bool
  }

  render() {
    const {
      canDelete,
      error,
      handleSubmit,
      invalid,
      closeForm,
      onDelete,
      pristine
    } = this.props;
    const { description, type, ldapDN, ldapGroupMemberAttribute } = fieldsWithErrorOnlyIfTouched(this.props.fields);

    return (
      <form onSubmit={ handleSubmit }>
        <div>
          <TextArea labelText='Description' { ...description } />

          {
            window.authMethod === 'ldap' ? (
              <RadioGroup { ...type }>
                <Radio value='managed'>Managed</Radio>
                <Radio value='ldap'>LDAP</Radio>
              </RadioGroup>
            ) : <input type='hidden' { ...type } value='managed' />
          }

          { type.value === 'ldap' &&
            (
              <div style={ {marginTop: '10px'} }>
                <Input labelText='LDAP DN' { ...ldapDN } />
                <Input labelText='Group Member Attribute' { ...ldapGroupMemberAttribute } />
              </div>
            )
          }
        </div>

        <ActionBar className={ styles.teamProfileActionBar } >
          <ErrorMessage error={ error } />
          { canDelete &&
            <Button ghost variant='alert' onClick={ onDelete }>Delete</Button>
          }
          <Button variant='secondary' onClick={ closeForm }>Cancel</Button>
          <Button disabled={ invalid || pristine } type='submit'>Save</Button>
        </ActionBar>
      </form>
    );
  }

}

class TeamDetails extends Component {
  static propTypes = {
    team: shape({
      name: string.isRequired,
      type: string,
      description: string,
      ldapDN: string,
      ldapGroupMemberAttribute: string
    })
  }

  render() {
    const { team } = this.props;

    return (
      <div>
        <p style={ {wordWrap: 'break-word'} }>{ team.description }</p>
      </div>
    );
  }
}
