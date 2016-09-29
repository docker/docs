'use strict';

import React, { PropTypes, Component } from 'react';
import connectToStores from 'fluxible-addons-react/connectToStores';
import attemptChangeLongDescription from '../../../../actions/attemptChangeLongDescription';
import longDescriptionUpdateFormField from '../../../../actions/longDescriptionUpdateFormField';
import toggleLongDescriptionEdit from '../../../../actions/toggleLongDescriptionEdit';
import RepoDetailsLongDescriptionFormStore from '../../../../stores/RepoDetailsLongDescriptionFormStore';
import Card, { Block } from '@dux/element-card';
import { Button } from 'dux';
import classnames from 'classnames';
import styles from './RepoFullDescription.css';
import Markdown from '@dux/element-markdown';

const { string, object, shape, bool, number, func } = PropTypes;
const debug = require('debug')('RepoLongDescription');


export default class RepoFullDescription extends Component {
  static propTypes = {
    isPrivate: bool.isRequired,
    isAutomated: bool.isRequired,
    name: string.isRequired,
    namespace: string.isRequired,
    status: number.isRequired,
    user: object,
    JWT: string,
    isEditing: bool.isRequired,
    successfulSave: bool,
    values: shape({
      longDescription: string.isRequired
    })
  }

  static contextTypes = {
    executeAction: func.isRequired
  }

  static defaultProps = {
    successfulSave: false
  }

  toggleEditMode = (e) => {
    this.context.executeAction(toggleLongDescriptionEdit, { isEditing: !this.props.isEditing });
  }

  onChange = (fieldKey) => {
    return (e) => {
      this.context.executeAction(longDescriptionUpdateFormField, {
        fieldKey,
        fieldValue: e.target.value
      });
    };
  }

  updateLongDescription = (event) => {
    event.preventDefault();
    const { JWT,
            name,
            namespace
          } = this.props;
    const { longDescription } = this.props.values;
    this.context.executeAction(attemptChangeLongDescription,
                               {
                                 jwt: JWT,
                                 repoShortName: namespace + '/' + name,
                                 longDescription: longDescription
                               });
  }

  renderForm = (cardHeading) => {
    let isError = false;
    let maybeError = <span />;
    const maxChars = 25000;
    const longDesc = this.props.fields.longDescription;
    const currentLongDesc = this.props.values.longDescription;
    if (longDesc.hasError) {
      /* check for error from bad request (post submit). success is handled by non-edit card */
      maybeError = (<div className={styles.errorText}>{longDesc.error}</div>);
      isError = true;
    } else if (currentLongDesc.length > maxChars) {
      /* check for input that's too long - not caught until submit otherwise */
      const tooMany = currentLongDesc.length - maxChars;
      const msg = `Maximum ${maxChars} characters. You are over the limit by: ${tooMany}`;
      maybeError = (<div className={styles.errorText}>{msg}</div>);
      isError = true;
    }
    const intent = isError ? 'alert' : 'primary';
    const textAreaClass = classnames({
      [styles.textArea]: true,
      [styles.error]: isError
    });

    const input = (
      <textarea type='textarea'
                rows='8'
                cols='200'
                className={textAreaClass}
                placeholder='Enter your full description for this repository.'
                onChange={this.onChange('longDescription')}
                autoFocus>
        {currentLongDesc}
      </textarea>
      );
    const submit = this.updateLongDescription;
    return (
      <Card heading={cardHeading}>
        <Block>
          <form className='row' onSubmit={submit}>
            <div className='large-12 columns'>
              {input}
              {maybeError}
              <div className='right'>
                <ul className='button-group radius'>
                  <li>
                    <Button type='button'
                            size='small'
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
    const { longDescription } = this.props.values;
    let cardHeading = `Full Description`;

    if (isEditing) {
      cardHeading = `Full Description (Optional, Limit 25,000 Characters)`;
      return this.renderForm(cardHeading);
    } else {
      const headingActions = canEdit ? [{
          key: 'edit',
          icon: 'fa-edit',
          action: this.toggleEditMode
        }] : [];
      let formattedDescription;
      if (!longDescription) {
        formattedDescription = (
          <div>
            <p>
              Full description is empty for this repo.
            </p>
          </div>);
      } else {
        formattedDescription = (<Markdown>{longDescription}</Markdown>);
      }
      const maybeSuccessClass = classnames({
        [styles.successfulSave]: successfulSave
      });

      let maybeSuccess = <span />;
      if (successfulSave) {
        maybeSuccess = (
          <div className={styles.successText}>
            {fields.longDescription.success}
          </div>
          );
      }
      return (
        <Card heading={cardHeading}
              headingActions={headingActions}
              className={maybeSuccessClass}>
          <Block>
            <div>
              {maybeSuccess}
              {formattedDescription}
            </div>
          </Block>
        </Card>
      );
    }
  }
}

export default connectToStores(RepoFullDescription,
                               [RepoDetailsLongDescriptionFormStore],
                               function({ getStore }, props) {
                                 return getStore(RepoDetailsLongDescriptionFormStore).getState();
                               });
