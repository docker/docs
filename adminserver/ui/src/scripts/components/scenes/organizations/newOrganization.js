'use strict';

import React, { Component, PropTypes } from 'react';
const { bool, string, func, object } = PropTypes;
import { reduxForm } from 'redux-form';
import Button from 'components/common/button';
import InputLabel from 'components/common/inputLabel';
import Input from 'components/common/input';
import FontAwesome from 'components/common/fontAwesome';
import styles from './list.css';
import css from 'react-css-modules';

@reduxForm({
  form: 'createOrganizationForm',
  fields: ['name']
})
@css(styles)
export default class NewOrganization extends Component {
  static propTypes = {
    dismissForm: func.isRequired,
    error: string,
    fields: object.isRequired,
    handleSubmit: func.isRequired,
    invalid: bool,
    resetForm: func.isRequired
  }

  static defaultProps = {
    invalid: true
  }

  render() {
    const {
      dismissForm,
      error,
      handleSubmit,
      invalid,
      resetForm
    } = this.props;
    const { name } = this.props.fields;

    return (
      <form onSubmit={ handleSubmit } id='new-org-form'>
        <div styleName='container'>
          <div styleName='avatar'>
            <FontAwesome styleName='orgAvatar' icon='fa-users' />
          </div>

          <div styleName='namespace'>
            <InputLabel
              tip='Where repositories will be created. Your organization&apos;s name needs to be unique.'>
              Organization name
            </InputLabel>
            <Input
              required
              formfield={ name } />
          </div>
        </div>

        <div styleName='actions'>
          <Button variant='primary outline' type='button' onClick={ () => resetForm() && dismissForm() }>Cancel</Button>
          <Button variant='primary' disabled={ invalid } type='submit'>Save</Button>
        </div>

        <div style={ {clear: 'both'} }>
          <p>{ error }</p>
        </div>
      </form>
    );
  }
}
