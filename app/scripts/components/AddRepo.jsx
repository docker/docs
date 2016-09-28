'use strict';

import React, { createClass, PropTypes } from 'react';
import _ from 'lodash';
import connectToStores from 'fluxible-addons-react/connectToStores';

import CreateRepositoryFormStore from 'stores/CreateRepositoryFormStore';
import DUXInput from './common/DUXInput.jsx';
import SimpleTextArea from './common/SimpleTextArea.jsx';
import Route404 from './common/RouteNotFound404Page.jsx';
import PrivateRepoUsageStore from 'stores/PrivateRepoUsageStore.js';
import RepositoryNameInput from 'common/RepositoryNameInput.jsx';
import createRepository from 'actions/createRepository';
import updateAddRepositoryFormField from 'actions/updateAddRepositoryFormField';
import getSettingsData from 'actions/getSettingsData';
import { STATUS } from 'stores/common/Constants';
import { SplitSection } from 'common/Sections';

import Markdown from '@dux/element-markdown';
import Card, { Block } from '@dux/element-card';
import { PageHeader, Button } from 'dux';
import styles from './AddRepo.css';

var debug = require('debug')('AddRepositoryForm');

var AddRepositoryForm = createClass({
  displayName: 'AddRepositoryForm',
  contextTypes: {
    executeAction: PropTypes.func.isRequired
  },
  propTypes: {
    JWT: PropTypes.string.isRequired,
    createRepoFormStore: PropTypes.shape({
      fields: PropTypes.object.isRequired,
      values: PropTypes.object.isRequired,
      namespaces: PropTypes.arrayOf(PropTypes.string.isRequired).isRequired,
      globalFormError: PropTypes.string
    }),
    privateRepoUsage: PropTypes.shape({
      privateRepoAvailable: PropTypes.number.isRequired
    })
  },
  getInitialState: function() {
    return {
      currentNamespace: this.props.createRepoFormStore.values.namespace
    };
  },
  _handleCreate: function(e) {
    e.preventDefault();

    var newRepo = {
      namespace: this.state.currentNamespace,
      name: this.props.createRepoFormStore.values.name.toLowerCase(),
      description: this.props.createRepoFormStore.values.description,
      full_description: this.props.createRepoFormStore.values.full_description,
      is_private: this.props.createRepoFormStore.values.is_private === 'private'
    };

    this.context.executeAction(createRepository,
                               {
                                 jwt: this.props.JWT,
                                 repository: newRepo
                               });
  },
  _onChange(fieldKey) {
    let _this = this;
    return (e) => {
      if (fieldKey === 'namespace') {
        this.setState({
          currentNamespace: e.target.value
        });
        this.context.executeAction(getSettingsData, {
          JWT: _this.props.JWT,
          username: e.target.value,
          repoType: 'regular'
        });
      }
      this.context.executeAction(updateAddRepositoryFormField, {
        fieldKey,
        fieldValue: e.target.value
      });
    };
  },
  _getCurrentQueryNamespace() {
    //Check if user has passed in namespace as query | verify if they have access to it
    var currentNamespace = this.props.location.query.namespace;
    if (!_.includes(this.props.createRepoFormStore.namespaces, currentNamespace)) {
      currentNamespace = this.props.createRepoFormStore.values.namespace;
    }
    return currentNamespace;
  },
  componentDidMount() {
    this.setState({
      currentNamespace: this._getCurrentQueryNamespace()
    });
  },
  _renderCreateForm() {
    const { createRepoFormStore } = this.props;

    var globalFormError = (<div />);
    if(createRepoFormStore.globalFormError) {
      globalFormError = (
        <div className='alert-box alert'>{createRepoFormStore.globalFormError}</div>
      );
    }

    var nameError = (<p></p>);
    var nameClass = '';
    if(createRepoFormStore.fields.name.hasError) {
      nameError = (
        <p className='alert-box alert large-12 columns'>
          {createRepoFormStore.fields.name.error +
          ' No spaces and special characters other than - and . are allowed. Repo names should not begin/end with a . or - .'}
        </p>
      );
      nameClass = 'text-error';
    }

    var fullDescriptionError = (<p></p>);
    var fullDescriptionClass = 'large-12 columns';
    if(createRepoFormStore.fields.full_description.hasError) {
      fullDescriptionError = (
        <p className='alert-box alert large-12 columns'>{createRepoFormStore.fields.full_description.error}</p>
      );
      fullDescriptionClass = 'large-12 columns form-error';
    }

    var descriptionError = (<p></p>);
    var descriptionClass = 'large-12 columns';
    if(createRepoFormStore.fields.description.hasError) {
      descriptionError = (
        <p className='alert-box alert large-12 columns'>{createRepoFormStore.fields.description.error}</p>
      );
      descriptionClass = 'large-12 columns form-error';
    }
    var defaultValue = createRepoFormStore.values.is_private ? 'private' : 'public';
    var errText = '';
    if (createRepoFormStore.fields.is_private.hasError) {
      errText =
        <span className='alert-box alert large-12 columns'>
          {createRepoFormStore.fields.is_private.error}
        </span>;
    }
    const subtitleContent = `1. Choose a namespace *(Required)*
2. Add a repository name *(Required)*
3. Add a short description
4. Add markdown to the full description field
5. Set it to be a private or public repository`;

    return (
      <SplitSection title=''
                    subtitle={<Markdown>{subtitleContent}</Markdown>}>
        <div>
          {globalFormError}
          <form onSubmit={this._handleCreate}>
            <div>
              <RepositoryNameInput namespaces={this.props.createRepoFormStore.namespaces}
                                   selectedNamespace={this._getCurrentQueryNamespace()}
                                   repoName={this.props.createRepoFormStore.values.name}
                                   onRepoNameChange={this._onChange('name')}
                                   onNamespaceChange={this._onChange('namespace')}
                                   inputClass={nameClass}/>
              {nameError}
            </div>

            <div className="row">
              <div className='large 12 columns'>
                <SimpleTextArea placeholder="Short Description (100 Characters)"
                                value={this.props.createRepoFormStore.values.description}
                                onChange={this._onChange('description')}
                                rows={2}
                                cols={100}
                                className={descriptionClass}/>
                {descriptionError}
              </div>
            </div>

            <div className="row">
              <div className='large-12 columns'>
                <SimpleTextArea placeholder="Full Description"
                                value={this.props.createRepoFormStore.values.full_description}
                                rows={5}
                                cols={100}
                                onChange={this._onChange('full_description')}
                                className={fullDescriptionClass} />
                {fullDescriptionError}
              </div>
            </div>

            <div className="row">
              <div className='large-12 columns'>
                <span className="text">Visibility</span>
                <select defaultValue={defaultValue}
                        value={this.props.createRepoFormStore.values.is_private}
                        onChange={this._onChange('is_private')}>
                  <option value="public">public</option>
                  <option value="private">private</option>
                </select>
                {errText}
              </div>
            </div>

            <Button type="submit">Create</Button>

          </form>
        </div>
      </SplitSection>
    );
  },
  render() {
    if (!this.props.JWT) {
      return (<Route404 />);
    } else {
      let content = (
        <div className={`row ${styles.contentWrapper}`}>
          <div className="columns large-8 large-centered">
            <Card>
              <Block>
                  <h2>Something went wrong!</h2>
                  <p>
                    We were unable to populate your user namespaces.
                    If this issue persists, please contact <a href="mailto:support@docker.com">support@docker.com</a>
                    &nbsp;or file an issue at <a href="https://github.com/docker/hub-feedback">hub-feedback on github</a>
                  </p>
              </Block>
            </Card>
          </div>
        </div>
      );
      if ( this.props.createRepoFormStore.namespaces && this.props.createRepoFormStore.namespaces.length > 0 ) {
        content = (
          <div className={styles.contentWrapper}>
            {this._renderCreateForm()}
          </div>
        );
      }

      return (
        <div>
          <PageHeader title='Create Repository'/>
          { content }
        </div>
      );
    }
  }
});

export default connectToStores(AddRepositoryForm,
                               [
                                 CreateRepositoryFormStore
                               ],
                               function({ getStore }, props) {
                                 return {
                                   createRepoFormStore: getStore(CreateRepositoryFormStore).getState()
                                 };
                               });
