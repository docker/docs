'use strict';

import React, { PropTypes } from 'react';
import findIndex from 'lodash/array/findIndex';
import includes from 'lodash/collection/includes';
import omit from 'lodash/object/omit';
import map from 'lodash/collection/map';
import AutoBuildTagsInput from './AutoBuildTagsInput.jsx';
import connectToStores from 'fluxible-addons-react/connectToStores';
import RepositoryNameInput from 'common/RepositoryNameInput.jsx';
import SimpleTextArea from 'common/SimpleTextArea.jsx';
import AutobuildStore from '../../stores/AutobuildStore';
import AutobuildConfigStore from '../../stores/AutobuildConfigStore';
import AutobuildSourceRepositoriesStore from '../../stores/AutobuildSourceRepositoriesStore';
import RepoStore from '../../stores/RepositoryPageStore';
import UserStore from '../../stores/UserStore';
import createAutobuild from '../../actions/createAutobuild';
import updateAutobuildFormField from '../../actions/updateAutobuildFormField.js';
import getSettingsData from 'actions/getSettingsData';
import { PageHeader } from 'dux';
import AlertBox from 'common/AlertBox';
import Card, { Block } from '@dux/element-card';
import Button from '@dux/element-button';
import { validateRepositoryName } from '../utils/validateRepositoryName';
import { STATUS as COMMONSTATUS } from '../../stores/common/Constants';
import Markdown from '@dux/element-markdown';

const {
  ATTEMPTING
} = COMMONSTATUS;

const buildTagsClientSideError = 'No empty strings allowed for docker tag (or) source tag/branch name specification.';

import styles from './AutoBuildSetupForm.css';

const {
  array,
  bool,
  func,
  number,
  object,
  oneOf,
  shape,
  string
} = PropTypes;

var AutoBuildSetupForm = React.createClass({
  contextTypes: {
    executeAction: func.isRequired
  },
  propTypes: {
    user: object.isRequired,
    JWT: string.isRequired,
    ownedNamespaces: array.isRequired,
    configStore: shape({
      description: string,
      isPrivate: oneOf(['private', 'public']).isRequired,
      name: string.isRequired,
      namespace: string.isRequired,
      sourceRepoName: string.isRequired,
      STATUS: string.isRequired
    }),
    sourceRepositories: shape({
      type: string.isRequired
    })
  },
  getInitialState: function() {
    return {
      isActive: true,
      buildTags: this.defaultBuildTags,
      clientSideError: '',
      advancedMode: false
    };
  },
  /*eslint-disable camelcase*/

  /*
   * By default, if the input is empty for source tag/branch name, send the string: '{sourceref}' & 'master'
   * By default, if the input is empty for docker tag name, send the string with regex for all matches & 'latest'
   */
  defaultBuildTags: [
    {
      id: 'tag-0',
      name: 'latest',
      source_type: 'Branch',
      source_name: 'master',
      dockerfile_location: '/'
    },
    {
      id: 'tag-1',
      name: '{sourceref}',
      source_type: 'Branch',
      source_name: '/^([^m]|.[^a]|..[^s]|...[^t]|....[^e]|.....[^r]|.{0,5}$|.{7,})/',
      dockerfile_location: '/'
    }
  ],
  getBuildTagsToSend: function() {
    let bTags = this.state.buildTags;
    return map(bTags, (tag) => {
      return omit(tag, 'id');
    });
  },
  /*eslint-enable camelcase */
  _handleCreate: function(evt) {
    evt.preventDefault();
    const { username } = this.props.user;
    const { buildTags, isActive } = this.state;
    const { description, isPrivate, name, namespace, sourceRepoName } = this.props.configStore;
    const { type } = this.props.sourceRepositories;

    const params = this.props.params;
    const sourceRepoFallback = `${params.sourceRepoNamespace}/${params.sourceRepoName}`;

    if (!validateRepositoryName(name.toLowerCase())) {
      //check if the repo name is valid | client side check
      this.setState({
        clientSideError: `No spaces and special characters other than '.' and '-' are allowed.
Repository names should not begin/end with a '.' or '-'.`
      });
    } else {

      var newAutobuild = {
        user: username,
        namespace: namespace,
        name: name.toLowerCase(),
        description: description,
        is_private: isPrivate === 'private',
        build_name: sourceRepoName.toLowerCase() || sourceRepoFallback.toLowerCase(),
        provider: type.toLowerCase(),
        active: isActive,
        tags: this.getBuildTagsToSend()
      };

      this.context.executeAction(createAutobuild, {JWT: this.props.JWT, autobuildConfig: newAutobuild});
    }
  },
  _onActiveStateChange: function(e) {
    this.setState({isActive: !this.state.isActive});
  },
  _getTagIndex: function(id) {
    return findIndex(this.state.buildTags, function(tag) {
      return (tag.id === id);
    });
  },
  _setTagState: function(id, prop, value) {
    let bTags = this.state.buildTags;
    bTags[this._getTagIndex(id)][prop] = value;  //find tag and update property
    this.setState({
      buildTags: bTags
    });
  },
  _onTagRemoved: function(id) {
    //Remove tag, when user removes it from the form
    let bTags = this.state.buildTags;
    bTags.splice(this._getTagIndex(id), 1);
    this.setState({
      buildTags: bTags
    });
  },
  _onTagAdded: function(tag) {
    let bTags = this.state.buildTags;
    bTags.push(tag);
    this.setState({
      buildTags: bTags
    });
  },
  _resetBuildTagsError: function() {
    //Reset clientSideError on change
    if (this.state.clientSideError === buildTagsClientSideError) {
      this.setState({
        clientSideError: ''
      });
    }
  },
  _onSourceNameChange: function(tagId, e) {
    let sourceName = e.target.value;
    let sourceType = this.state.buildTags[this._getTagIndex(tagId)].sourceType;
    if (sourceName === '' && sourceType && sourceType === 'Branch') {
      sourceName = '/^([^m]|.[^a]|..[^s]|...[^t]|....[^e]|.....[^r]|.{0,5}$|.{7,})/';
    } else if (sourceName === '' && sourceType && sourceType === 'Tag') {
      sourceName = '/.*/';
    } else {
      this._resetBuildTagsError();
      //Strip trailing and leading spaces. If we end up with empty string, throw an error.
      if (sourceName.trim() === '') {
        this.setState({
          clientSideError: buildTagsClientSideError
        });
      }
    }
    this._setTagState(tagId, 'source_name', sourceName.trim());
  },
  _onSourceTypeChange: function(tagId, e) {
    this._setTagState(tagId, 'source_type', e.target.value);
  },
  _onDockerfileLocationChange: function(tagId, e) {
    this._setTagState(tagId, 'dockerfile_location', e.target.value);
  },
  _onTagChange: function(tagId, e) {
    let tagName = e.target.value;
    if (tagName === '') {
      tagName = '{sourceref}';
    }
    this._setTagState(tagId, 'name', tagName.trim());
  },
  _updateForm(fieldKey) {
    return (e) => {
      if (fieldKey === 'namespace') {
        this.setState({
          currentNamespace: e.target.value
        });
        this.context.executeAction(getSettingsData, {
          JWT: this.props.JWT,
          username: e.target.value,
          repoType: 'autobuild'
        });
      } else if (fieldKey === 'name' && this.state.clientSideError) {
        this.setState({
          clientSideError: ''
        });
      }
      this.context.executeAction(updateAutobuildFormField, {
        fieldKey,
        fieldValue: e.target.value
      });
    };
  },
  componentWillReceiveProps: function(nextProps) {
    const { name, namespace, success, STATUS } = nextProps.configStore;
    //If autobuild was created successfully
    if (STATUS.SUCCESSFUL || success) {
      this.props.history.pushState(null, `/r/${namespace}/${name.toLowerCase()}/`);
    }
  },
  customTagsConfig: function() {
    this.setState({
      advancedMode: true,
      buildTags: []
    });
  },
  defaultTagsConfig: function() {
    this.setState({
      advancedMode: false,
      buildTags: this.defaultBuildTags
    });
  },
  render: function() {

    const {
      description,
      error,
      isPrivate,
      name,
      namespace,
      success,
      STATUS
    } = this.props.configStore;

    /* start error/success handling */
    let maybeSuccess = <span />;
    if (success) {
      maybeSuccess = <span className='alert-box success radius'>{success}</span>;
    }

    let nameError;
    let nameErrorContent = error.dockerhub_repo_name;
    if(nameErrorContent) {
      nameError = nameErrorContent;
    }

    let descriptionError;
    if(error.description) {
      descriptionError = error.description;
    }

    let privateRepoError;
    if(error.is_private) {
      privateRepoError = error.is_private;
    }

    let buildTagsError;
    if (error.buildTags) {
      buildTagsError = error.buildTags;
    }

    let maybeError = null;
    let errorDetail = error.detail || this.state.clientSideError;
    if (errorDetail) {
      maybeError = (
        <div className='row'>
          <div className={styles.globalError}>
            {errorDetail}
          </div>
        </div>
      );
    }
    /* end error handling */

    //Check if user has passed in namespace as query | verify if they have access to it
    let currentUserNamespace = this.props.location.query.namespace;
    if (!includes(this.props.ownedNamespaces, currentUserNamespace)) {
      //If they don't have access to the namespace set in the query param ?<ns> then fallback to default namespace
      currentUserNamespace = this.props.user.namespace;
    }
    let tagsConfigList = null;

    if (this.state.advancedMode) {
      tagsConfigList = (
        <div>
          <br />
          <h5>Customize Autobuild Tags</h5>
          <label className={styles.customizeLabel}>
            Your image will build automatically when your source repository is pushed based on the following rules.
            <a onClick={this.defaultTagsConfig}> Revert to default settings</a>
          </label>
          <div className={styles.error}>
            {buildTagsError}
          </div>
          <AutoBuildTagsInput repo={name}
                              onTagRemoved={this._onTagRemoved}
                              onTagAdded={this._onTagAdded}
                              onSourceNameChange={this._onSourceNameChange}
                              onSourceTypeChange={this._onSourceTypeChange}
                              onDockerfileLocationChange={this._onDockerfileLocationChange}
                              onTagChange={this._onTagChange}/>
        </div>
      );
    }

    return (
      <div>
        <PageHeader title='Create Automated Build'/>
        <div className={'row ' + styles.formContainer}>
          <div className='columns large-9 end'>
            <Card>
              <Block>
                <form onSubmit={this._handleCreate} className="create-autobuild-form">
                    <div className="row">
                      <div className="columns large-8">
                        <label className={styles.label}>Repository Namespace & Name<sup>*</sup></label>
                        <div className={styles.error}>
                          {nameError}
                        </div>
                        <RepositoryNameInput namespaces={this.props.ownedNamespaces}
                                           selectedNamespace={currentUserNamespace || namespace}
                                           repoName={name}
                                           onRepoNameChange={this._updateForm('name').bind(this)}
                                           onNamespaceChange={this._updateForm('namespace').bind(this)}/>
                      </div>
                      <div className='columns large-offset-1 large-3'>
                        <label className={styles.label}>Visibility</label>
                        <div className={styles.error}>
                          {privateRepoError}
                        </div>
                        <select className={styles.select}
                                value={isPrivate}
                                onChange={this._updateForm('isPrivate').bind(this)}>
                          <option value="public">public</option>
                          <option value="private">private</option>
                        </select>
                      </div>
                    </div>
                    <label className={styles.label}>Short Description<sup>*</sup></label>
                    <div className={styles.error}>
                      {descriptionError}
                    </div>
                    <SimpleTextArea value={description}
                                    placeholder="Max 100 Characters"
                                    onChange={this._updateForm('description').bind(this)}
                                    rows={2}
                                    cols={100} />
                    <label className={styles.customizeLabel}>
                      By default Automated Builds will match branch names to Docker build tags.
                      <a onClick={this.customTagsConfig}> Click here to customize</a> behavior.
                    </label>
                    {/* advanced mode */}
                    {tagsConfigList}
                    <br />
                    {maybeError}
                    <div className='row'>
                      <div className={styles.floatRight}>
                        <Button variant='primary' type='submit'>
                          {STATUS === ATTEMPTING ? 'Creating...' : 'Create'}
                        </Button>
                      </div>
                    </div>
                </form>
                {maybeSuccess}
              </Block>
            </Card>
          </div>
        </div>
      </div>
    );
  }
});

export default connectToStores(AutoBuildSetupForm,
                [
                  AutobuildSourceRepositoriesStore,
                  AutobuildConfigStore,
                  UserStore
                ],
                function({ getStore }, props){
                  return {
                    sourceRepositories: getStore(AutobuildSourceRepositoriesStore).getState(),
                    configStore: getStore(AutobuildConfigStore).getState(),
                    ownedNamespaces: getStore(UserStore).getNamespaces()
                  };
                });
