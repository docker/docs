'use strict';

import React, { Component, PropTypes } from 'react';
const { array, func, object, string } = PropTypes;
import { reduxForm } from 'redux-form';
// components
import InputLabel from 'components/common/inputLabel';
import Input from 'components/common/input';
import RepoVisibility from 'components/common/repoVisibility';
import CancelAndDone from 'components/common/cancelAndDone';
import Dropdown from 'components/common/dropdown';
// validation
import {
  createValidator,
  required,
  maxLength,
  regex
} from 'validation';
// misc
import styles from './createRepoForm.css';
import css from 'react-css-modules';
import { Repositories } from 'dtr-js-sdk';

@reduxForm({
  form: 'addRepo',
  fields: ['namespace', 'name', 'shortDescription', 'visibility'],
  validate: createValidator({
    namespace: [required],
    name: [required, maxLength(30), regex(/^[a-zA-Z0-9]+(?:[._-][a-zA-Z0-9]+)*$/, `Must start with alphanumeric characters and use '._-' as separators`)],
    visibility: [required]
  }),
  asyncValidate: (data) => {
    return Repositories.default.getRepository({
      repo: data.name,
      namespace: data.namespace
    })
      .then(() => {
        // Error if we found the repo with this name
        return {
          name: 'Repository name already exists'
        };
      })
      .catch(() => {
        // Otherwise we're good
        return {};
      });
  }
}, (state) => ({
  initialValues: {
    namespace: state.users.get('currentUser').name,
    visibility: 'public'
  }
}))
@css(styles)
export default class AddRepoForm extends Component {

  static propTypes = {
    fields: object,
    // handleSubmit is a reduxForm-specific prop which calls onSubmit
    handleSubmit: func,
    onSubmit: func,
    onCancel: func,

    // We can only add a repository to your account or any org you're a member
    // of.
    orgNames: array,
    // The name of the current logged-in user
    username: string.isRequired
  }

  static defaultProps = {
    orgNames: []
  }

  render() {
    const {
      username,
      orgNames,
      fields: {
        namespace,
        name,
        shortDescription,
        visibility
      },
      handleSubmit,
      onCancel
    } = this.props;

    return (
      <form onSubmit={ handleSubmit } id='create-repo-form'>
        <div styleName='row'>
          <div styleName='account'>
            <InputLabel>Account</InputLabel>
            <Dropdown values={ namespace }>
              <option value={ username }>{ username }</option>
              { orgNames.map(n => <option value={ n } key={ n }>{ n }</option> ) }
            </Dropdown>
          </div>
          <div styleName='name'>
            <span styleName='divider' >/</span>
            <InputLabel>Repository name</InputLabel>
            <Input required formfield={ name } />
          </div>
        </div>

        <div styleName='row'>
          <div styleName='description'>
            <InputLabel isOptional={ true }>Description</InputLabel>
            <Input formfield={ shortDescription } />
          </div>
          <div styleName='visibility'>
            <RepoVisibility formField={ visibility } />
          </div>
        </div>

        <CancelAndDone onCancel={ onCancel } />
      </form>
    );
  }

}
