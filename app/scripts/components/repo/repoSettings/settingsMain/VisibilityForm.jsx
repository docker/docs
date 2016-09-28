'use strict';

import React, { PropTypes, Component } from 'react';
import { Link } from 'react-router';
import connectToStores from 'fluxible-addons-react/connectToStores';
import toggleVisibilityRepoNameConfirmBox from 'actions/toggleVisibilityRepoNameConfirmBox.js';
import updateFormField from 'actions/common/updateFormField';
import updateRepoVisibility from 'actions/updateRepoSettingsVisibilityField';
import RepoDetailsVisibilityFormStore from 'stores/RepoDetailsVisibilityFormStore.js';
import { STATUS as COMMONSTATUS } from 'stores/repovisibilitystore/Constants.js';
import classnames from 'classnames';
import has from 'lodash/object/has';
import capitalize from 'lodash/string/capitalize';
import styles from './VisibilityForm.css';
import Card, { Block } from '@dux/element-card';
import Button from '@dux/element-button';
const { func, bool, shape, number, string, object } = PropTypes;
const debug = require('debug')('VisibilityForm');

class VisibilityForm extends Component {

  static propTypes = {
    JWT: string.isRequired,
    name: string.isRequired,
    namespace: string.isRequired,
    user: object.isRequired,
    visibilityFormStore: shape({
      badRequest: string,
      error: string,
      success: string,
      isPrivate: bool,
      numPrivateReposAvailable: number,
      privateRepoLimit: number,
      STATUS: string,
      values: shape({
        confirmRepoName: string
      })
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
    const { STATUS } = this.props.visibilityFormStore;
    if (STATUS !== COMMONSTATUS.DEFAULT && e.keyCode === 27) {
      this.toggleConfirmBox();
    }
  }

  updateVisibility = (visibility, isDisabled) => {
    return (e) => {
      e.preventDefault();
      const { JWT, namespace, name } = this.props;
      if (!isDisabled) {
        this.context.executeAction(updateRepoVisibility,
        {
          jwt: JWT,
          repoShortName: namespace + '/' + name,
          isPrivate: (visibility === 'private')
        });
      }
    };
  }

  toggleConfirmBox = (e) => {
    this.context.executeAction(toggleVisibilityRepoNameConfirmBox);
  }

  onChange = (fieldKey) => {
    return (e) => {
      this.context.executeAction(updateFormField({
        storePrefix: 'REPO_DETAILS_VISIBILITY'
      }), {
        fieldKey,
        fieldValue: e.target.value
      });
    };
  }

  isOrg = () => {
    const { user, namespace } = this.props;
    // username is the logged in user. namespace is the current org / user
    // if namespace != username, then it is an org
    if (has(user, 'username')) {
      return user.username !== namespace;
    } else {
      //if there is no user.username, there is an user.orgname
      return true;
    }
  }

  getBillingLink = () => {
    const { namespace } = this.props;
    const linkToBilling = this.isOrg() ? `/u/${namespace}/dashboard/billing/` : '/account/billing-plans/';
    return (
      <Link to={ linkToBilling }>
        Get more private repositories.
      </Link>
    );
  }

  renderConfirmBox = () => {
    const { badRequest,
            isPrivate,
            error,
            STATUS } = this.props.visibilityFormStore;
    const { confirmRepoName } = this.props.visibilityFormStore.values;
    const { name } = this.props;
    const isCorrect = confirmRepoName === name;
    const newVisibility = isPrivate ? 'public' : 'private';
    let button, maybeError;
    if (isCorrect) {
      const btnText = STATUS === COMMONSTATUS.ATTEMPTING ? 'Submitting...'
                                                         : 'Confirm';
      button = (
        <Button variant='primary'
                ghost
                type='submit'>
          {btnText}
        </Button>
      );
    }
    if (error || badRequest) {
      maybeError = <div className={styles.error}>{error || badRequest}</div>;
    }
    const confirmTitle = (
      <div className={styles.text}>
        {`Please type the name of your repository to make it ${newVisibility}:`}
        <strong> {name}</strong>
      </div>
    );
    const headingActions = [{
        key: 'cancel',
        icon: 'fa-remove',
        action: this.toggleConfirmBox
      }];
    return (
      <Card heading='Visibility Settings'
            headingActions={headingActions}>
        <Block>
          {confirmTitle}
          <form onSubmit={this.updateVisibility(newVisibility, !isCorrect)}>
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
    const { badRequest,
            error,
            isPrivate,
            numPrivateReposAvailable,
            privateRepoLimit,
            STATUS } = this.props.visibilityFormStore;

    if (STATUS !== COMMONSTATUS.DEFAULT) {
      return this.renderConfirmBox();
    } else {
      const newVisibility = isPrivate ? 'public' : 'private';
      let privateDisabled = numPrivateReposAvailable === 0;
      let helpText;
      const button = (
        <Button variant='primary'
                disabled={!isPrivate && privateDisabled}
                onClick={this.toggleConfirmBox}>
          Make {capitalize(newVisibility)}
        </Button>
      );
      const title = `Make this Repository ${capitalize(newVisibility)}`;

      if (isPrivate) {
        helpText = `Public repositories are available to anyone and will ` +
           `appear in public search results. `;
      } else {
        let numLeftText;
        if (numPrivateReposAvailable !== null && privateRepoLimit !== null) {
          numLeftText = `You are using ` +
            `${privateRepoLimit - numPrivateReposAvailable} of ` +
            `${privateRepoLimit} private repositories. `;
        }
        if (privateDisabled) {
          helpText = (<div>{numLeftText}{this.getBillingLink()}</div>);
        } else {
          const privateText = `Private repositories are only available to you or members of your organization. `;
          helpText = <div>{privateText}<br />{numLeftText}</div>;
        }
      }
      return (
        <Card heading='Visibility Settings'>
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

export default connectToStores(
  VisibilityForm,
  [ RepoDetailsVisibilityFormStore ],
  function({ getStore }, props) {
    return {visibilityFormStore: getStore(RepoDetailsVisibilityFormStore).getState()};
  }
);
