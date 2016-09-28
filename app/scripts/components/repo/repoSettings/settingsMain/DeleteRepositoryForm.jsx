'use strict';

import React, { PropTypes, Component } from 'react';
import connectToStores from 'fluxible-addons-react/connectToStores';
import toggleDeleteRepoNameConfirmBox from 'actions/toggleDeleteRepoNameConfirmBox';
import updateFormField from 'actions/common/updateFormField';
import deleteRepositoryAction from 'actions/deleteRepo';
import DeleteRepoFormStore from 'stores/DeleteRepoFormStore';
import Card, { Block } from '@dux/element-card';
import Button from '@dux/element-button';
import { STATUS as COMMONSTATUS } from 'stores/deleterepostore/Constants';
import styles from './DeleteRepositoryForm.css';
const debug = require('debug')('DeleteRepoForm');
const { string, func, shape, object } = PropTypes;

class DeleteRepositoryForm extends Component {
  static propTypes = {
    JWT: string.isRequired,
    name: string.isRequired,
    params: object.isRequired,
    deleteRepoFormStore: shape({
      error: string,
      STATUS: string
    })
  }

  static contextTypes = {
    executeAction: func.isRequired
  }

  addEscapeKeyHandler = () => {
    document.addEventListener('keydown', this.handleEscapeKey);
  }

  removeEscapeKeyHandler = () => {
    document.removeEventListener('keydown', this.handleEscapeKey);
  }

  handleEscapeKey = (e) => {
    e = e || window.event;
    const { STATUS } = this.props.deleteRepoFormStore;
    if (STATUS !== COMMONSTATUS.DEFAULT && e.keyCode === 27) {
      this.toggleConfirmBox();
    }
  }

  onChange = (fieldKey) => {
    return (e) => {
      this.context.executeAction(updateFormField({
        storePrefix: 'DELETE_REPO'
      }), {
        fieldKey,
        fieldValue: e.target.value
      });
    };
  }

  onDelete = (e) => {
    e.preventDefault();
    const { JWT, name } = this.props;
    const { confirmRepoName } = this.props.deleteRepoFormStore.values;
    if (confirmRepoName === name) {
      const { user, splat } = this.props.params;
      this.context.executeAction(deleteRepositoryAction,
        {
          jwt: JWT,
          repoShortName: `${user}/${splat}`
        }
      );
    }
  }

  toggleConfirmBox = (e) => {
    this.context.executeAction(toggleDeleteRepoNameConfirmBox);
  }

  renderConfirmBox = () => {
    const { error, STATUS } = this.props.deleteRepoFormStore;
    const { confirmRepoName } = this.props.deleteRepoFormStore.values;
    const { name } = this.props;

    const isCorrect = confirmRepoName === name;
    let button, maybeError;
    if (isCorrect) {
      const btnText = STATUS === COMMONSTATUS.ATTEMPTING ? 'Deleting...'
                                                         : 'Delete';
      button = (
        <Button variant='alert'
                ghost
                type='submit'>
          {btnText}
        </Button>
      );
    }
    if (error) {
      maybeError = <div className={styles.error}>{error}</div>;
    }
    const confirmTitle = (
      <div className={styles.text}>
        {`Please type the name of your repository to confirm deletion:`}
        <strong> {name}</strong>
      </div>
    );
    const headingActions = [{
        key: 'cancel',
        icon: 'fa-remove',
        action: this.toggleConfirmBox
      }];
    return (
      <Card heading='Delete Repository'
            headingActions={headingActions}>
        <Block>
          {confirmTitle}
          <form onSubmit={this.onDelete}>
            <input className={styles.textArea}
                   type="text"
                   value={confirmRepoName}
                   onChange={this.onChange('confirmRepoName')}
                   onFocus={this.addEscapeKeyHandler}
                   onBlur={this.removeEscapeKeyHandler}
                   autoFocus />
            {maybeError}
            {button}
          </form>
        </Block>
      </Card>
      );
  }

  render() {
    const { STATUS } = this.props.deleteRepoFormStore;
    if (STATUS !== COMMONSTATUS.DEFAULT) {
      return this.renderConfirmBox();
    } else {
      const button = (
        <Button variant='alert'
                onClick={this.toggleConfirmBox}>
          Delete
        </Button>
      );
      const title = `Delete Repository`;
      const helpText = (
        <div>
          Deleting a repository will <strong>destroy</strong> all images stored within it!
          This action is <strong>not reversible.</strong>
        </div>
      );
      return (
        <Card heading='Delete Repository'>
          <Block>
            <div className='row'>
              <div className='large-9 columns'>
                <h4 className={styles.text}>{title}</h4>
              </div>
              <div className='large-3 columns'>
                <div className='right'>
                  {button}
                </div>
              </div>
            </div>
            <div className='row'>
              <div className='large-12 columns'>
                <div className={styles.text}>
                  {helpText}
                </div>
              </div>
            </div>
          </Block>
        </Card>
      );
    }
  }
}

export default connectToStores(DeleteRepositoryForm,
  [
    DeleteRepoFormStore
  ],
  function({ getStore }, props) {
    return { deleteRepoFormStore: getStore(DeleteRepoFormStore).getState() };
  });
