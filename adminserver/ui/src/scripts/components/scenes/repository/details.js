'use strict';

import React, { Component, PropTypes } from 'react';
const { object, string, func } = PropTypes;
import { connect } from 'react-redux';
import ReactMarkdown from 'react-markdown';

import * as RepositoriesActions from 'actions/repositories';
import autoaction from 'autoaction';
// components
import CancelAndDone from 'components/common/cancelAndDone';
import Input from 'components/common/input';
import InputLabel from 'components/common/inputLabel';
import { permissionsByAccessLevel } from 'components/common/docSnippets';
import QTip from 'components/common/qtip';
import CopyToClipboard from 'components/common/copyToClipboard';
import Button from 'components/common/button';

import styles from './repository.css';
import css from 'react-css-modules';
import { humanPermissions, mapActions } from 'utils';
import {
  getAccessLevel,
  getRawRepoState,
  getRepoForName
} from 'selectors/repositories';
import { reduxForm } from 'redux-form';
import ui from 'redux-ui';

const mapRepositoryState = (state, props) => ({
  repositories: getRawRepoState,
  accessLevel: getAccessLevel(state, props),
  repo: getRepoForName(state, props)
});

@connect(mapRepositoryState, mapActions(RepositoriesActions))
@ui({
  state: {
    showEditTools: false,
    descriptionEditMode: false
  }
})
@reduxForm({
  form: 'repoLongDescription',
  fields: ['longDescription']
}, (_, props) => {
  return {
    initialValues: {
      'longDescription': props.repo.longDescription
    }
  };
})
@autoaction({
  getRepository: (props) => ({
    namespace: props.params.namespace,
    repo: props.params.name
  })
}, RepositoriesActions)
@css(styles, {allowMultiple: true})
export class RepositoryDetailsTab extends Component {

  static propTypes = {
    actions: object,
    params: object.isRequired,
    repositories: object,
    accessLevel: string,
    repo: object,
    handleSubmit: func
  };

  onSubmit = (data) => {
    this.props.actions.updateRepository({
      namespace: this.props.params.namespace,
      repo: this.props.params.name
    }, data);
    this.toggleEditMode();
  };

  toggleEditTools = () => {
    this.props.updateUI({
      showEditTools: !this.props.ui.showEditTools
    });
  };

  toggleEditMode = () => {
    this.props.updateUI({
      descriptionEditMode: !this.props.ui.descriptionEditMode
    });
  };

  renderAccessLevel() {

    const {
      accessLevel
    } = this.props;

    return (
      <div>
        <InputLabel>Your permission</InputLabel>
        <p>{ humanPermissions(accessLevel) } <QTip tooltip={
            <div>
              <p>Your repository permissions:</p>
              <ul>{ permissionsByAccessLevel[accessLevel] }</ul>
            </div> } />
        </p>
      </div>
    );
  }

  render() {

    const {
      repo,
      params: {
        namespace,
        name
      },
      handleSubmit,
      ui: {
        descriptionEditMode,
        showEditTools
      },
      fields
    } = this.props;

    /** TODO: add this below CopyToClipboard:
     * <InputLabel>Owners</InputLabel>
     * <ul><FontAwesome icon='fa-person'/>{ TODO: LIST OF OWNERS HERE! }</ul>
     **/
    return (
      <div>{
        repo && (
          <div styleName='repositoryInfo'>
            <form styleName='longDescriptionForm' method='post' onSubmit={ handleSubmit(::this.onSubmit) }>
              <div>
                <div styleName='row'>
                  <h2 styleName='readmeHeader'>README</h2>
                  { !descriptionEditMode ?
                    <div styleName='readmeButton'>
                      <Button type='button' variant='primary outline slim' onClick={ ::this.toggleEditMode }>Edit</Button>
                    </div>
                    : undefined
                  }
                </div>
                <span
                  style={ descriptionEditMode ? { display: 'none' } : undefined }
                  onClick={ ::this.toggleEditMode }
                  styleName={ showEditTools ? 'showEditTools' : undefined }>
                  <p styleName='longDescription'>
                    <ReactMarkdown
                      source={ repo.longDescription ? repo.longDescription : '*README is empty for this repository.*' }
                    />
                  </p>
                </span>
                <div
                  style={ descriptionEditMode ? undefined : { display: 'none' } }>
                  <Input isTextarea formfield={ fields.longDescription } />
                  <CancelAndDone onCancel={ ::this.toggleEditMode } />
                </div>
              </div>
            </form>

            <div styleName='detailColumn'>
              <InputLabel>Docker pull command</InputLabel>
              <CopyToClipboard>
                { `docker pull ${location.host}/${namespace}/${name}` }
              </CopyToClipboard>
              { this.renderAccessLevel() }
            </div>
          </div>
        )
      }</div>
    );
  }

}
