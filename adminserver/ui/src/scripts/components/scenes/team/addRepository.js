'use strict';

import React, { Component, PropTypes } from 'react';
const { bool, func, object, string } = PropTypes;
import { connect } from 'react-redux';
import { reduxForm } from 'redux-form';
import { reset as resetReduxForm } from 'redux-form';
import { createStructuredSelector } from 'reselect';

import Button from 'components/common/button';
import {
  ErrorMessage,
  Input,
  Label,
  Radio,
  RadioGroup,
  RepositoryName,
  Select,
  TextArea
} from 'components/common';
import { permissionTooltip, visibilityTooltip } from 'components/common/docSnippets';
import {
  createValidator,
  required,
  maxLength,
  oneOf,
  onlyIf
} from 'validation';
import { nonVisibleTeamRepositoriesSelector } from 'selectors/teams';
import * as RepoActions from 'actions/repositories';
import * as TeamActions from 'actions/teams';
import consts from 'consts';
import styles from './team.css';
import { fieldsWithErrorOnlyIfTouched, mapActions } from 'utils';
import ui from 'redux-ui';

const { teams } = consts;

const addRepoValidator = createValidator({
  permissions: [required, oneOf(['admin', 'read-write', 'read-only'])],
  repository: [onlyIf(data => data.repoType === teams.ADD_EXISTING_REPO, required)],
  name: [onlyIf(data => data.repoType === teams.ADD_NEW_REPO, required)],
  visibility: [onlyIf(data => data.repoType === teams.ADD_NEW_REPO, required)],
  shortDescription: [onlyIf(data => data.repoType === teams.ADD_NEW_REPO, maxLength(140))]
});
@reduxForm({
  form: 'addRepoToTeam',
  fields: [
    'permissions',
    'repoType', // either existing or new repo

    // existing repos fields
    'repository',
    'allRepositories',

    // new repo fields
    'name',
    'visibility',
    'shortDescription'
  ],
  validate: addRepoValidator
})
class AddRepoForm extends Component {
  static propTypes = {
    actions: object.isRequired,
    fields: object.isRequired,
    closeForm: func.isRequired,
    handleSubmit: func.isRequired,
    invalid: bool,
    onSubmit: func.isRequired,
    error: string,
    resetForm: func.isRequired,

    organization: object.isRequired,
    repos: object.isRequired
  }

  renderRepositories() {
    const { repos } = this.props;
    return (repos || []).map( item => {
      const name = item.get('repo').name;
      return <option value={ name } key={ name }>{ name }</option>;
    });
  }

  renderAddExistingRepo() {
    const { repository, allRepositories } = fieldsWithErrorOnlyIfTouched(this.props.fields);
    let possiblyDisabledRepository = { ...repository };
    if (allRepositories.checked) {
      possiblyDisabledRepository.disabled = true;
      possiblyDisabledRepository.value = '';
    }
    return (
      <RepositoryName>
        <Label text='Account'>
          { this.props.organization.name }
        </Label>
        <Select
          labelText='Repository name'
          { ...possiblyDisabledRepository } defaultValue=''>
          <option value='' disabled>Select repository</option>
          { this.renderRepositories() }
        </Select>
      </RepositoryName>
    );
  }

  renderAddNewRepo() {
    const { visibility, shortDescription, name } = fieldsWithErrorOnlyIfTouched(this.props.fields);
    return (
      <div>
        <RepositoryName>
          <Label text='Account'>
            { this.props.organization.name }
          </Label>
          <Input labelText='Repository name' {...name} />
        </RepositoryName>

        <RadioGroup
          labelText='Visibility'
          tooltip={ <p>{ visibilityTooltip }</p> }
          {...visibility}
        >
          <Radio value='private'>Private</Radio>
          <Radio value='public'>Public</Radio>
        </RadioGroup>

        <TextArea labelText='Description (optional)' {...shortDescription} />
      </div>
    );
  }

  render() {
    const { handleSubmit, invalid, error, repos, resetForm } = this.props;
    const { permissions, repoType } = fieldsWithErrorOnlyIfTouched(this.props.fields);

    return (
      <form onSubmit={ handleSubmit }>
        <Select
          labelText='Permission'
          tooltip={ <p>{ permissionTooltip }</p> }
          { ...permissions }
          defaultValue=''>
          <option value='' disabled>Select permission</option>
          <option value='admin'>Admin</option>
          <option value='read-write'>Read-write</option>
          <option value='read-only'>Read only</option>
        </Select>

        <RadioGroup labelText='Repository' { ...repoType }>
          <Radio value={ teams.ADD_EXISTING_REPO } disabled={ repos && repos.size === 0 }>Existing</Radio>
          <Radio value={ teams.ADD_NEW_REPO }>New</Radio>
        </RadioGroup>

        {
          repoType.value === teams.ADD_EXISTING_REPO
            ? this.renderAddExistingRepo()
            : this.renderAddNewRepo()
        }


        <div className={ styles.addRepoActions }>
          <Button variant='secondary' onClick={ () => resetForm() && this.props.closeForm() }>Cancel</Button>
          { ' ' }
          <Button disabled={ invalid } type='submit'>Done</Button>
        </div>
        <ErrorMessage error={ error } />
      </form>
    );
  }
}


/**
 * This component allows us to add access rights to a new or existing repo for
 * the current team in an organization.
 *
 * Because this depends on the org to load all repositories in the org's
 * namespace, the organization **must** be fully loaded and passed in as
 * a property to use this form.
 *
 * ## Displaying the correct form
 *
 * This connects to 'state.form' directly to inspect the 'repoType' value; when
 * this value changes we render either the new or existing repo form accordingly
 *
 */
const mapState = createStructuredSelector({
  repos: nonVisibleTeamRepositoriesSelector,
  form: (state) => state.form
});

@ui()
@connect(
  mapState,
  mapActions({
    form: {resetReduxForm},
    repos: RepoActions,
    teams: TeamActions
  })
)
export class AddRepository extends Component {

  static propTypes = {
    actions: object.isRequired,
    form: object.isRequired,

    organization: object.isRequired,
    repos: object.isRequired,
    teamName: string.isRequired,

    updateUI: func
  }

  closeForm() {
    this.props.updateUI('isFormVisible', false);
  }

  onSubmit(data) {
    const teamName = this.props.teamName;

    let promise;

    if (data.repoType === teams.ADD_NEW_REPO) {
      promise = this.props.actions.repos.createRepoAndGrantTeamAccess({
        orgName: this.props.organization.name,
        teamName
      }, {
        name: data.name,
        visibility: data.visibility,
        shortDescription: data.shortDescription,
        accessLevel: data.permissions
      });
    } else if (data.allRepositories) {
      promise = this.props.actions.repos.grantTeamAccessToRepoNamespace({
        orgName: this.props.organization.name,
        teamName,
        accessLevel: data.permissions
      });
    } else {
      promise = this.props.actions.repos.grantTeamAccessToRepo({
        orgName: this.props.organization.name,
        teamName,
        repo: data.repository,
        accessLevel: data.permissions
      });
    }
    return promise.then(() => {
      // On successful submission, reset the form.
      // On error, an {_error: ... } object will be returned from action creator instead
      this.props.actions.form.resetReduxForm('addRepoToTeam');
    });
  }

  render() {
    const { repos } = this.props;
    return (
      <div className={ styles.addRepoForm }>
          <AddRepoForm
            { ...this.props }
            onSubmit={ ::this.onSubmit }
            closeForm={ ::this.closeForm }
            initialValues={ {
              permissions: '',
              repoType: (repos && repos.size === 0) ? teams.ADD_NEW_REPO : teams.ADD_EXISTING_REPO,

              // existing repos fields
              repository: '',

              // new repo fields
              name: '',
              visibility: 'public',
              shortDescription: ''
            } }
          />
      </div>
    );
  }
}
