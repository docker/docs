'use strict';

import React, { PropTypes, Component } from 'react';
import { connect } from 'react-redux';
import { updateRepository } from 'actions/repositories';
import Button from 'components/common/button';

import ButtonSelect from 'components/common/buttonSelect';

import { mapActions } from 'utils';
import styles from '../repository.css';
import css from 'react-css-modules';
// import FontAwesome from 'components/common/fontAwesome';
import { reduxForm } from 'redux-form';
import { getRawRepoState, getRepoForName } from 'selectors/repositories';

const mapRepositoryState = (state, props) => ({
  repositories: getRawRepoState,
  repo: getRepoForName(state, props)
});

@connect(
  mapRepositoryState,
  mapActions({ updateRepository })
)
@reduxForm({
  form: 'repoPrivacyAndShortDescription',
  fields: [
    'visibility',
    'shortDescription'
  ]
}, (_, props) => {
  const {
    visibility,
    shortDescription
    } = props.repo || {};
  return {
    initialValues: {
      'visibility': visibility,
      'shortDescription': shortDescription
    }
  };
})
@css(styles, {allowMultiple: true})
export default class RepositoryEditForm extends Component {
  static propTypes = {
    actions: PropTypes.object,
    params: PropTypes.object,
    fields: PropTypes.object,
    repo: PropTypes.object,
    handleSubmit: PropTypes.func,
    pristine: PropTypes.bool,
    resetForm: PropTypes.func
  };

  onSubmit(data) {
    this.props.actions.updateRepository({
      namespace: this.props.params.namespace,
      repo: this.props.params.name
    }, data);
  }

  render() {
    const {
      fields,
      handleSubmit,
      resetForm,
      pristine
    } = this.props;

    return (
      <div styleName='wrapper'>
        <form
          method='post'
          onSubmit={ handleSubmit(::this.onSubmit) }
          styleName='editRepository'
        >
          <div>
            <div styleName='field-title'>Visibility</div>
            <ButtonSelect
              initialChoice={ fields.visibility.value }
              onChange={ fields.visibility.onChange }
            >
              <div
                icon='fa-globe'
                primaryText='Public'
                secondaryText='Visible to everyone'
                value='public'
              />
              <div
                icon='fa-lock'
                primaryText='Private'
                secondaryText='Hide this repository'
                value='private'
              />
            </ButtonSelect>
          </div>
          <div>
            <div styleName='field-title'>Description</div>
            <div styleName='shortDescriptionInput'>
              <input {...fields.shortDescription} type='text'/>
            </div>
          </div>
          <div styleName='actions'>
            <Button
              type='button'
              variant='outline primary'
              disabled={ pristine }
              onClick={ resetForm }
            >
              Cancel
            </Button>
            <Button
              type='submit'
              variant='primary'
              disabled={ pristine }
            >
              Save
            </Button>
          </div>
        </form>
      </div>
    );
  }
}
