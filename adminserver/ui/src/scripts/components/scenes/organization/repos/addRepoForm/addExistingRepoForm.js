'use strict';

import React, { Component, PropTypes } from 'react';
const { func, instanceOf, array, object } = PropTypes;
// components
import InputLabel from 'components/common/inputLabel';
import Input from 'components/common/input';
import Dropdown from 'components/common/dropdown';
import CancelAndDone from 'components/common/cancelAndDone';
// actions
import autoaction from 'autoaction';
import { mapActions } from 'utils';
import {
  listRepositories,
  grantTeamAccessToRepo
} from 'actions/repositories';
// selectors
import {
  getReposByNamespace
} from 'selectors/repositories';
import {
  createSelector,
  createStructuredSelector
} from 'reselect';
import { connect } from 'react-redux';
// misc
import { get } from 'lodash';
import { reduxForm } from 'redux-form';
import { OrganizationRecord } from 'records';
import { Map } from 'immutable';
import styles from './addRepoForm.css';
import css from 'react-css-modules';

// validation
import {
  createValidator,
  required
} from 'validation';

/**
 * Get the selected namespace from redux-form
 */
const getSelectedNamespaceName = (state) =>
  get(state, 'form.addTeamRepo.namespace.value');

/**
 * Given all repos by namespace and a namespace name, return the namespace's
 * repos as an array of names
 */
const getSelectedNamespaceRepos = createSelector(
  getReposByNamespace,
  getSelectedNamespaceName,
  (reposByNamespace, namespaceName) => {
    return reposByNamespace
      .getIn([namespaceName, 'entities', 'repo'], new Map())
      .map(r => r.name)
      .toArray();
  }
);

const mapState = createStructuredSelector({
  repoNames: getSelectedNamespaceRepos
});

// Connect needs to be top-level as reduxForm depends on the
// getSelectedNamespaceRepos results from mapState.
@connect(mapState, mapActions({ grantTeamAccessToRepo }))
@reduxForm({
  form: 'addTeamRepo',
  fields: ['namespace', 'repo', 'accessLevel'],
  validate: createValidator({
    namespace: [required],
    repo: [required],
    accessLevel: [required]
  })
}, (state, props) => ({
  initialValues: {
    namespace: props.org.name,
    repo: props.repoNames[0],
    accessLevel: 'read-only'
  }
}))
@autoaction({
  listRepositories: (props) => ({ namespace: props.org.name })
}, { listRepositories })
@css(styles)
export default class AddExistingRepoForm extends Component {

  static propTypes = {
    actions: object,
    onHide: func,
    handleSubmit: func,

    repoNames: array,
    org: instanceOf(OrganizationRecord),
    fields: object,
    params: object
  }

  /**
   * This function is passed to redux-form's handleSubmit function to call our
   * action for creating a new repo then adding the repo to the team
   */
  onSubmit(formData) {
    this.props.actions.grantTeamAccessToRepo({
      orgName: this.props.org.name,
      teamName: this.props.params.team,
      repo: formData.repo,
      accessLevel: formData.accessLevel
    });
  }

  /**
   * TODO: Don't add repositories this team already has access to to the repo
   * dropdown
   */
  render() {
    const {
      onHide,
      handleSubmit,
      repoNames,
      fields: { namespace, repo, accessLevel },
      org: { name }
    } = this.props;

    return (
      <form onSubmit={ handleSubmit(::this.onSubmit) }>
        <div styleName='row'>
          <div styleName='account'>
            <InputLabel>Account</InputLabel>
            <Input { ...namespace } placeholder={ name } disabled />
          </div>

          <div styleName='repo'>
            <InputLabel>Repository name</InputLabel>
            <Dropdown values={ { ...repo, defaultValue: repoNames[0] } }>
              { repoNames.length === 0
                ? <option disabled>This account has no repositories</option>
                : repoNames.map(n => <option key={ n } value={ n }>{ n }</option>) }
            </Dropdown>
          </div>

          <div styleName='permissions'>
            <InputLabel>Permissions</InputLabel>
            <Dropdown values={ accessLevel } defaultValue='permission'>
              <option disabled value='permission'>Permission</option>
              <option value='admin'>Admin</option>
              <option value='read-write'>Read-write</option>
              <option selected={ accessLevel.value === 'read-only' } value='read-only'>Read-only</option>
            </Dropdown>
          </div>
        </div>
        <CancelAndDone onCancel={ onHide } />
      </form>
    );
  }
}
