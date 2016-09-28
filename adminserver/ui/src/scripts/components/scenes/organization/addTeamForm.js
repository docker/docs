'use strict';

import React, { Component, PropTypes } from 'react';
const { bool, object, func } = PropTypes;
// components
import InputLabel from 'components/common/inputLabel';
import Input from 'components/common/input';
import Button from 'components/common/button';
import RadioGroup from 'components/common/radioGroup';
// actions
import { connect } from 'react-redux';
import { mapActions } from 'utils';
import {
  createTeam
} from 'actions/teams';
// validation
import {
  createValidator,
  required
} from 'validation';
// misc
import { reduxForm } from 'redux-form';
import css from 'react-css-modules';
import styles from './addTeamForm.css';

const validateForm = createValidator({
  type: [required],
  name: [required]
});

@connect(() => ({}), mapActions({ createTeam }))
@reduxForm({
  form: 'addTeam',
  fields: ['type', 'name', 'ldapDN', 'ldapGroupMemberAttribute'],
  initialValues: {
    type: 'managed'
  },
  validate: validateForm
})
@css(styles)
export default class AddTeamForm extends Component {

  static propTypes = {
    actions: object,
    handleSubmit: func,
    onHide: func,
    isLdapEnabled: bool,
    fields: object,
    params: object
  };

  onSubmit() {
    const { fields } = this.props;

    this.props.actions.createTeam({
      orgName: this.props.params.org,
      type: fields.type.value,
      name: fields.name.value,
      ldapDN: fields.ldapDN.value,
      ldapGroupMemberAttribute: fields.ldapGroupMemberAttribute.value
    });
  }

  showTypeFields() {
    const {
      fields: { type }
    } = this.props;

    return (
      <div styleName='type'>
        <InputLabel>Type</InputLabel>
        <RadioGroup
          formField={ type }
          vertical
          initialChoice='managed'
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
        ldapDN,
        ldapGroupMemberAttribute
      }
    } = this.props;

    return (
      <div styleName='ldap'>
        <InputLabel>LDAP DN</InputLabel>
        <Input formfield={ ldapDN } />
        <InputLabel>LDAP Group Member Attribute</InputLabel>
        <Input formfield={ ldapGroupMemberAttribute } />
      </div>
    );
  }

  onCancel(evt) {
    evt.preventDefault();
    this.props.onHide();
  }

  render() {
    const {
      isLdapEnabled,
      handleSubmit,
      fields: { name, type }
    } = this.props;

    return (
      <form id='create-team-form' styleName='wrapper' onSubmit={ handleSubmit(::this.onSubmit) }>
        <div styleName='fields'>
          { isLdapEnabled && this.showTypeFields() }
          { type.value === 'ldap' && this.showLdapFields() }
          <InputLabel>Team name</InputLabel>
          <Input formfield={ name } />
        </div>

        <div styleName='buttons'>
          <div>
            <Button
              type='button'
              variant='primary simple'
              onClick={ ::this.onCancel }>Cancel</Button>
            <Button variant='primary' type='submit'>Save</Button>
          </div>
        </div>
      </form>
    );
  }

}
