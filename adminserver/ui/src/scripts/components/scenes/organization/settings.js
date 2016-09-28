'use strict';

import React, { Component, PropTypes } from 'react';
import Button from 'components/common/button';
import styles from 'components/scenes/settings/formstyle.css';
import css from 'react-css-modules';
import DeleteModal from 'components/common/deleteModal';
import { connectModal } from 'components/common/modal';
import { mapActions } from 'utils';
import { connect } from 'react-redux';
// actions
import { deleteTeam, updateTeam, getTeam } from 'actions/teams';
import { deleteOrganization } from 'actions/organizations';
import { getAuthSettings } from 'actions/settings';
import autoaction from 'autoaction';
// selectors
import { createStructuredSelector } from 'reselect';
import { teamDetailsByName } from 'selectors/teams';
import { getAuthMethod } from 'selectors/settings';
import { isAdminSelector } from 'selectors/users';
// components
import InputLabel from 'components/common/inputLabel';
import Input from 'components/common/input';
import RadioGroup from 'components/common/radioGroup';
import Spinner from 'components/common/spinner';
// validation
import {
  createValidator,
  required
} from 'validation';
// misc
import { reduxForm } from 'redux-form';
import consts from 'consts';

const validateForm = createValidator({
  type: [required],
  name: [required]
});

const mapState = createStructuredSelector({
  team: teamDetailsByName,
  authMethod: getAuthMethod,
  isSysAdmin: isAdminSelector
});

@connect(mapState, mapActions({
  deleteTeam,
  updateTeam,
  deleteOrganization
}))
@autoaction({
  getAuthSettings: [],
  getTeam: (props) => [props.params.org, props.params.team]
}, { getTeam, getAuthSettings })
export default class OrgSettings extends Component {

  static propTypes = {
    params: PropTypes.object,
    hideModal: PropTypes.func,
    showModal: PropTypes.func,
    actions: PropTypes.object,
    fields: PropTypes.object,
    handleSubmit: PropTypes.func,
    team: PropTypes.object,
    isSysAdmin: PropTypes.bool
  }

  render() {
    const status = [
      [consts.settings.AUTH_SETTINGS],
      [consts.teams.GET_TEAM]
    ];

    return (
      <Spinner loadingStatus={ status }>
        { this.props.params.team === undefined
          ? <EditOrg { ...this.props } />
          : <EditTeam { ...this.props } /> }
      </Spinner>
    );
  }

}

@reduxForm({
  form: 'editTeam',
  fields: ['type', 'name', 'groupDN', 'groupMemberAttr'],
  validate: validateForm
}, (_, props) => ({
  initialValues: {
    ...props.team.toJS(),
    type: props.team.enableSync === true ? 'ldap' : 'managed'
  }
}))
@connectModal()
@css(styles, { allowMultiple: true })
class EditTeam extends Component {

  static propTypes = {
    params: PropTypes.object,
    hideModal: PropTypes.func,
    showModal: PropTypes.func,
    actions: PropTypes.object,
    fields: PropTypes.object,
    handleSubmit: PropTypes.func,
    team: PropTypes.object
  }


  showModal = () => {
    const { team } = this.props.params;

    this.props.showModal(
      (
        <span id='delete-org-modal'>
          <DeleteModal
            resourceType='team'
            resourceName={ team }
            onDelete={ ::this.delete }
            hideModal={ this.props.hideModal }
          />
        </span>
      )
    );
  }

  delete() {
    const { org, team } = this.props.params;
    this.props.actions.deleteTeam({
      orgName: org,
      teamName: team
    });
    this.props.hideModal();
  }

  showTypeFields() {
    const {
      fields: { type },
      authMethod,
      team
    } = this.props;

    if (authMethod !== 'ldap') {
      return null;
    }

    return (
      <div style={ { marginBottom: '20px' } }>
        <InputLabel>Type</InputLabel>
        <RadioGroup
          formField={ type }
          initialChoice={ team.enableSync === true ? 'ldap' : 'managed' }
          choices={ [
            { label: 'Default', value: 'managed' },
            { label: 'LDAP', value: 'ldap' }
          ] } />
      </div>
    );
  }

  showLdapFields() {
    const {
      fields: {
        groupDN,
        groupMemberAttr
      }
    } = this.props;

    return (
      <div>
        <InputLabel>LDAP DN</InputLabel>
        <Input formfield={ groupDN } />
        <InputLabel>LDAP Group Member Attribute</InputLabel>
        <Input formfield={ groupMemberAttr } />
      </div>
    );
  }

  onSubmit(data) {
    this.props.actions.updateTeam(this.props.params.org, this.props.params.team, data);
  }

  render() {
    const {
      handleSubmit,
      fields: { name, type }
    } = this.props;

    return (
      <span>
        <form onSubmit={ handleSubmit(::this.onSubmit) }>
          <div styleName='formbox'>
            { this.showTypeFields() }
            { type.value === 'ldap' && this.showLdapFields() }
            <InputLabel>Team name</InputLabel>
            <Input formfield={ name } />
            <Button variant='primary' type='submit'>Save</Button>
          </div>
        </form>

        <div styleName='row'>
          <div styleName='formbox halfColumn'>
            <div styleName='row'>
              <p styleName='halfColumn'>Delete this team?</p>
              <div styleName='halfColumn right'>
                <Button
                  id='delete-org-button'
                  variant='alert'
                  onClick={ ::this.showModal }>
                  Delete
                </Button>
              </div>
            </div>
          </div>
        </div>
      </span>
    );
  }
}

@connectModal()
@css(styles, { allowMultiple: true })
class EditOrg extends Component {

  static propTypes = {
    params: PropTypes.object,
    hideModal: PropTypes.func,
    showModal: PropTypes.func,
    actions: PropTypes.object,
    isSysAdmin: PropTypes.bool
  }

  delete = () => {
    const { org } = this.props.params;
    this.props.actions.deleteOrganization({
      name: org
    });
    this.props.hideModal();
  }

  showModal() {
    const { org } = this.props.params;
    this.props.showModal(
      (
        <span id='delete-org-modal'>
          <DeleteModal
            resourceType={ 'org' }
            resourceName={ org }
            onDelete={ ::this.delete }
            hideModal={ this.props.hideModal }
          />
        </span>
      )
    );
  }

  render() {
    const { org } = this.props.params;

    if (!this.props.isSysAdmin) {
      return <span />;
    }

    return (
      <span>
        <div styleName='row'>
          <div styleName='formbox halfColumn'>
            <div styleName='row'>
              <p styleName='halfColumn'>Delete this organization?</p>
              <div styleName='halfColumn right'>
                <Button
                  id='delete-org-button'
                  variant='alert'
                  onClick={ ::this.showModal }
                  disabled={ org === 'docker-datacenter' }>
                  Delete
                </Button>
              </div>
            </div>
          </div>
          { org === 'docker-datacenter' && <DeleteDDCWarning /> }
        </div>
      </span>
    );
  }
}

@css(styles, { allowMultiple: true })
class DeleteDDCWarning extends Component {
  render () {
    return (
      <div styleName='halfColumn'>
        <h3>Cannot delete <strong>docker-datacenter!</strong></h3>
        <p>You need this organization for Datacenter's<br/> authentication to work.</p>
      </div>
    );
  }
}
