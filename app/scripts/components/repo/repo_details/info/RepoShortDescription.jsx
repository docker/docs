'use strict';

import React, { PropTypes, Component } from 'react';
import connectToStores from 'fluxible-addons-react/connectToStores';
import attemptChangeShortDescription from 'actions/attemptChangeShortDescription';
import RepoDetailsShortDescriptionFormStore from 'stores/RepoDetailsShortDescriptionFormStore';
import shortDescriptionUpdateFormField from 'actions/shortDescriptionUpdateFormField';
import toggleShortDescriptionEdit from 'actions/toggleShortDescriptionEdit';
import Card, { Block } from '@dux/element-card';
import { Button } from 'dux';
import classnames from 'classnames';
import styles from './RepoShortDescription.css';
import SimpleInput from 'common/SimpleInput';
const { string, object, bool, number, func, shape } = PropTypes;
const debug = require('debug')('RepoShortDescription');


class RepoShortDescription extends Component {
  static propTypes = {
    canEdit: bool.isRequired,
    name: string.isRequired,
    namespace: string.isRequired,
    JWT: string,
    isEditing: bool.isRequired,
    successfulSave: bool,
    values: shape({
      shortDescription: string
    })
  }

  static contextTypes = {
    executeAction: func.isRequired
  }

  static defaultProps = {
    successfulSave: false,
    isEditing: false
  }

  toggleEditMode = (e) => {
    this.context.executeAction(toggleShortDescriptionEdit, { isEditing: !this.props.isEditing });
  }

  onChange = (fieldKey) => {
    return (e) => {
      this.context.executeAction(shortDescriptionUpdateFormField, {
        fieldKey,
        fieldValue: e.target.value
      });
    };
  }

  updateShortDescription = (event) => {
    event.preventDefault();
    const { JWT,
            name,
            namespace
          } = this.props;
    const { shortDescription } = this.props.values;
    this.context.executeAction(attemptChangeShortDescription,
                               {
                                 jwt: JWT,
                                 repoShortName: namespace + '/' + name,
                                 shortDescription: shortDescription
                               });
  }
  renderForm = (cardHeading) => {
    let maybeError = <span />;
    const shortDesc = this.props.fields.shortDescription;
    const currentShortDesc = this.props.values.shortDescription;
    const maxChars = 100;
    let hasError = false;
    if (shortDesc.hasError) {
      /* check for error from request (post submit). success is handled by non-edit card */
      maybeError = (<div className={styles.errorText}>{shortDesc.error}</div>);
      hasError = true;
    } else if (currentShortDesc.length > maxChars) {
      /* check for typed input that's too long - not caught until submit otherwise */
      const tooMany = currentShortDesc.length - maxChars;

      const msg = `Maximum ${maxChars} characters. You are over the limit by: ${tooMany}`;
      maybeError = (<div className={styles.errorText}>{msg}</div>);
      hasError = true;
    }
    const intent = hasError ? 'alert' : 'primary';
    const textAreaClass = classnames({
      [styles.textArea]: true,
      [styles.error]: hasError
    });
    const input = (
      <SimpleInput type='text'
                   autoFocus={true}
                   hasError={hasError}
                   onChange={this.onChange('shortDescription')}
                   placeholder='Enter your short description'
                   value={currentShortDesc} />
    );
    const submit = this.updateShortDescription;
    return (
      <Card heading={cardHeading}>
        <Block>
          <form className={'row ' + styles.input} onSubmit={submit}>
            <div className='large-12 columns'>
              {input}
              {maybeError}
              <div className='right'>
                <ul className='button-group radius'>
                  <li>
                    <Button size='small'
                            intent='secondary'
                            onClick={this.toggleEditMode}>
                      Cancel
                    </Button>
                  </li>
                  <li>
                    <Button type='submit'
                            size='small'
                            intent={intent}>
                      Save
                    </Button>
                  </li>
                </ul>
              </div>
            </div>
          </form>
        </Block>
      </Card>
    );
  }

  render() {
    const {
      canEdit,
      isEditing,
      successfulSave,
      fields
    } = this.props;
    const { shortDescription } = this.props.values;
    let cardHeading = `Short Description`;

    if (isEditing) {
      cardHeading = `Short Description (Optional, Limit 100 Characters)`;
      return this.renderForm(cardHeading);
    } else {
      let headingActions = canEdit ? [{
          key: 'edit',
          icon: 'fa-edit',
          action: this.toggleEditMode
        }] : [];
      let maybeSuccessClass = classnames({
        [styles.successfulSave]: successfulSave
      });
      let maybeSuccess = <span />;
      if (successfulSave) {
        maybeSuccess = (<div className={styles.successText}>
                          {fields.shortDescription.success}
                        </div>);
      }
      return (
        <Card heading={cardHeading}
              headingActions={headingActions}
              className={maybeSuccessClass}>
          <Block>
            <div>
              {maybeSuccess}
              {shortDescription || 'Short description is empty for this repo.'}
            </div>
          </Block>
        </Card>
      );
    }
  }
}

export default connectToStores(RepoShortDescription,
                               [RepoDetailsShortDescriptionFormStore],
                               function({ getStore }, props) {
                                 return getStore(RepoDetailsShortDescriptionFormStore).getState();
                               });
