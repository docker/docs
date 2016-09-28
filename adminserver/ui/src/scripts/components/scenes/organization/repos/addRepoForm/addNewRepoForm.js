'use strict';

import React, { Component, PropTypes } from 'react';
const { func, instanceOf, object } = PropTypes;
// components
import InputLabel from 'components/common/inputLabel';
import Input from 'components/common/input';
import RepoVisibility from 'components/common/repoVisibility';
import CancelAndDone from 'components/common/cancelAndDone';
import Dropdown from 'components/common/dropdown';
// actions
import { mapActions } from 'utils';
import {
  createRepoAndGrantTeamAccess
} from 'actions/repositories';
// selectors
import { connect } from 'react-redux';
// misc
import { reduxForm } from 'redux-form';
import { OrganizationRecord } from 'records';
import styles from './addRepoForm.css';
import css from 'react-css-modules';
// validation
import {
  createValidator,
  required
} from 'validation';
import { listOrgOrTeamRepos } from 'actions/repositories';

@connect((() => {})(), mapActions({ createRepoAndGrantTeamAccess, listOrgOrTeamRepos }))
@reduxForm({
  form: 'addTeamRepo',
  fields: ['repo', 'accessLevel', 'shortDescription', 'visibility'],
  validate: createValidator({
    repo: [required],
    accessLevel: [required],
    visibility: [required]
  }),
  initialValues: {
    accessLevel: 'read-only',
    visibility: 'public'
  }
})
@css(styles)
export default class AddExistingRepoForm extends Component {

  static propTypes = {
    actions: object,
    onHide: func,
    handleSubmit: func,
    org: instanceOf(OrganizationRecord),
    fields: object,
    params: object
  }

  /**
   * This function is passed to redux-form's handleSubmit function to call our
   * action for adding the repo to the team.
   */
  onSubmit(data) {
    const {
      params: {
        org,
        team
      }
    } = this.props;

    const { repo, visibility, shortDescription, accessLevel } = data;

    this.props.actions.createRepoAndGrantTeamAccess({
      namespace: org,
      teamName: team || '',
      repo,
      visibility,
      shortDescription,
      accessLevel
    });

    this.props.onHide();

  }

  render() {
    const {
      onHide,
      handleSubmit,
      fields: { repo, accessLevel, shortDescription, visibility },
      org: { name }
    } = this.props;

    return (
      <form onSubmit={ handleSubmit(::this.onSubmit) }>
        <div styleName='row'>
          <div styleName='account'>
            <InputLabel>Account</InputLabel>
            <Input placeholder={ name } disabled />
          </div>

          <div styleName='repo'>
            <InputLabel>Repository name</InputLabel>
            <Input type='text' formfield={ repo } />
          </div>

          <div styleName='permissions'>
            <InputLabel>Permissions</InputLabel>
            <Dropdown { ...accessLevel }>
              <option disabled selected>Permission</option>
              <option value='admin'>Admin</option>
              <option value='read-write'>Read-write</option>
              <option value='read-only'>Read-only</option>
            </Dropdown>
          </div>
        </div>

        <div styleName='row'>
          <div styleName='description'>
            <InputLabel isOptional>Description</InputLabel>
            <Input type='text' formfield={ shortDescription } />
          </div>
          <div styleName='visibility'>
            <RepoVisibility formField={ visibility } />
          </div>
        </div>
        <CancelAndDone onCancel={ onHide } />
      </form>
    );
  }
}
