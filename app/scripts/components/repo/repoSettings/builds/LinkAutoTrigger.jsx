'use strict';

import React, {
  PropTypes,
  createClass
  } from 'react';

import { Link } from 'react-router';
import addTriggerLink from 'actions/addTriggerLink.js';
import removeTriggerLink from 'actions/removeTriggerLink.js';
import updateAutoBuildSettingsStore from 'actions/updateAutoBuildSettingsStore.js';
import SimpleInput from 'common/SimpleInput.jsx';
import includes from 'lodash/collection/includes';
import map from 'lodash/collection/map';
import { SplitSection } from '../../../common/Sections.jsx';
import Button from '@dux/element-button';
import Card, { Block } from '@dux/element-card';
import FA from '../../../common/FontAwesome.jsx';
import classnames from 'classnames';
import styles from './LinkAutoTrigger.css';

var debug = require('debug')('BuildOptions');

var _mkLinkDisplayRow = function(clickHandler) {
  return (repoLink) => {
    let username, reponame, link;
    if (includes(repoLink.to_repo, '/')) {
      [username, reponame] = repoLink.to_repo.split('/');
      link = (
        <Link to={`/r/${username}/${reponame}/`}>
          {repoLink.to_repo}
        </Link>
      );
    } else {
      reponame = repoLink.to_repo;
      link = (
        <Link to={`/r/_/${reponame}/`}>
          {reponame}
        </Link>
      );
    }

    return (
      <div key={repoLink.to_repo} className={styles.tag}>
        <div className={styles.repoName}>
          {link}
        </div>
        <div className={styles.removeLink} onClick={clickHandler(repoLink.id)}>
          <FA icon='fa-times' />
        </div>
      </div>
    );
  };
};

var LinkAutoTrigger = createClass({
  displayName: 'LinkAutoTrigger',
  propTypes: {
    autoBuildStore: PropTypes.object.isRequired,
    autoBuildLinks: PropTypes.array.isRequired,
    triggerLinkForm: PropTypes.object.isRequired,
    namespace: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired,
    validations: PropTypes.object.isRequired
  },
  contextTypes: {
    executeAction: React.PropTypes.func.isRequired
  },
  _addTriggerLink(e) {
    e.preventDefault();
    debug('ADDING TRIGGER LINK');
    let repoName = this.props.triggerLinkForm.repoName;
    if (!includes(repoName, '/')) {
      repoName = 'library/' + repoName;
    }
    this.context.executeAction(addTriggerLink, {JWT: this.props.JWT,
      namespace: this.props.namespace,
      name: this.props.name,
      to_repo: repoName});
  },
  _deleteLink: function(repo_id) {
    return (e) => {
      e.preventDefault();
      this.context.executeAction(removeTriggerLink, {JWT: this.props.JWT,
        namespace: this.props.namespace,
        name: this.props.name,
        repo_id: repo_id});
    };
  },
  _updateFormField(field, key) {
    return (e) => {
      e.preventDefault();
      this.context.executeAction(updateAutoBuildSettingsStore, {field: field, key: key, value: e.target.value});
    };
  },
  render() {
    const {
      autoBuildLinks,
      triggerLinkForm,
      validations
    } = this.props;

    let success = validations.links.success ? 'Link successfully added' : null;
    let maybeError = null;
    const errorClasses = classnames(['large-12 column', styles.error]);
    if (validations.links.hasError) {
      maybeError = (
        <div className={errorClasses}>
          {validations.links.error}
        </div>
      );
    }
    let linkedRepos = 'You have not linked any repositories to this Automated Build.';
    if (autoBuildLinks.length) {
      linkedRepos = map(autoBuildLinks, _mkLinkDisplayRow(this._deleteLink));
    }
    return (
      <Card heading='Repository Links'>
        <Block>
          <form onSubmit={this._addTriggerLink}>
            <div className='row'>
              <div className='large-9 columns'>
                <div>
                  Link your Automated Build to another Docker Hub repository,
                  and when that repository is updated, it will automatically
                  trigger a rebuild of this Automated Build.
                </div>
                <div>
                  <div className={styles.label}>Repository Name</div>
                  <div className='row'>
                    <div className='large-9 columns'>
                      <SimpleInput type="text"
                                   label="Repository Full Name"
                                   placeholder='Ex. ubuntu or username/reponame'
                                   hasError={validations.links.hasError}
                                   success={success}
                                   value={triggerLinkForm.repoName}
                                   onChange={this._updateFormField('triggerLinkForm', 'repoName')}/>
                    </div>
                    <div className="large-3 columns">
                      <Button variant='primary' size='tiny' type='submit'>
                        Add Repository Link
                      </Button>
                    </div>
                    {maybeError}
                  </div>
                </div>
                <div>
                  <div className={styles.label}>Linked Repositories</div>
                  <div className={styles.tagsWrapper}>{linkedRepos}</div>
                </div>
              </div>
            </div>
          </form>
        </Block>
      </Card>
    );
  }
});

export default LinkAutoTrigger;
