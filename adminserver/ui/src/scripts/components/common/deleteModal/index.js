'use strict';

import React, { Component, PropTypes } from 'react';
const { string, func, object } = PropTypes;
import Button from 'components/common/button';
import Input from 'components/common/input';
import FontAwesome from 'components/common/fontAwesome';
import styles from './deleteModal.css';
import ui from 'redux-ui';


/**
 * DeleteModal is a reusable confirmation dialog for destructive deletes.
 *
 */
@ui({
  state: {
    value: ''
  }
})
export default class DeleteModal extends Component {

  static propTypes = {
    resourceType: string.isRequired, // ie. 'organization', 'tag' etc
    resourceName: string.isRequired, // ie. 'my-repo'
    onDelete: func.isRequired,
    // Pass in the hideModal function from @connectModal in your component
    hideModal: func.isRequired,
    // UI
    updateUI: func.isRequired,
    ui: object.isRequired
  }

  onChange(evt) {
    const { value } = evt.target;
    this.props.updateUI('value', value);
  }

  onSubmit() {
    if (this.props.ui.value.toUpperCase() !== 'DELETE') {
      return;
    }
    this.props.onDelete();
    this.props.updateUI('value', '');
  }

  onCancel() {
    this.props.updateUI('value', '');
    this.props.hideModal();
  }

  render() {
    let { resourceType, resourceName } = this.props;
    if (resourceName.length > 30) {
        resourceName = `${resourceName.substr(0, 30)}\u2026`;
    }

    const { value } = this.props.ui;

    return (
      <div className={ styles.deleteModal } id='delete-modal'>
          <div>
            <div className={ styles.container }>
              <div className={ styles.iconContainer }>
                <FontAwesome icon='fa-exclamation-triangle' size='4x' />
              </div>
              <div className={ styles.formContainer }>
                <h2>Delete { resourceType }</h2>
                <p>{ resourceName } will be deleted. This cannot be undone!</p>
                <p>Type 'DELETE' to continue.</p>
                <Input onChange={ ::this.onChange } type='text' value={ value } />
              </div>
            </div>
          </div>
          <div className={ styles.buttons }>
            <div>
              <Button variant='primary simple' onClick={ ::this.onCancel } type='button'>Cancel</Button>
              <Button type='submit' variant='alert' disabled={ value.toUpperCase() !== 'DELETE' } onClick={ ::this.onSubmit }>Delete</Button>
            </div>
          </div>
      </div>
    );
  }

}
