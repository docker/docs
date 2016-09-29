'use strict';

import React, {
  PropTypes,
  createClass
  } from 'react';
import SimpleInput from 'common/SimpleInput.jsx';
import addPushTriggerItem from 'actions/addAutoBuildPushTriggerItem.js';
import deletePushTriggerItem from 'actions/deleteAutoBuildPushTriggerItem.js';
import saveRepoPushSettings from 'actions/savePushAutoBuildSettings.js';
import updatePushTriggerItem from 'actions/updateAutoBuildPushTriggerItem.js';
import baseClone from 'lodash/internal/baseClone';
import filter from 'lodash/collection/filter';
import map from 'lodash/collection/map';
import isEmpty from 'lodash/lang/isEmpty';
import isString from 'lodash/lang/isString';
import remove from 'lodash/array/remove';
import { Row, Header, Item, FlexTable as Table } from 'common/TagsFlexTable';
import Button from '@dux/element-button';
import FA from 'common/FontAwesome.jsx';
import Card, { Block } from '@dux/element-card';
import BuildTriggerByTag from './BuildTriggerByTag';

import {
  SOURCE_NAME_TAG,
  SOURCE_NAME_BRANCH,
  DOCKER_TAG,
  SOURCE_NAME_PLACEHOLDER_BRANCH,
  SOURCE_NAME_PLACEHOLDER_TAG,
  DOCKER_TAG_PLACEHOLDER_BRANCH,
  DOCKER_TAG_PLACEHOLDER_TAG
} from 'common/enums/BuildTagDefault';

import styles from './RepoPushTrigger.css';
var debug = require('debug')('BuildOptions');

const _isRegexString = (str) => {
  if (isString(str)) {
    return (str.charAt(0) === '/' && str.charAt(str.length - 1) === '/');
  } else {
    return false;
  }
};

var RepoPushTrigger = createClass({
  displayName: 'RepoPushTrigger',
  propTypes: {
    autoBuildStore: PropTypes.object.isRequired,
    newTags: PropTypes.array.isRequired,
    autoBuildBlankSlate: PropTypes.object.isRequired,
    JWT: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired,
    namespace: PropTypes.string.isRequired,
    validations: PropTypes.shape({
      buildTags: PropTypes.object.isRequired
    })
  },
  contextTypes: {
    executeAction: React.PropTypes.func.isRequired
  },
  _addPushItem(e) {
    e.preventDefault();
    this.context.executeAction(addPushTriggerItem);
  },
  _preprocessTags(tags) {
    //Remove empty/null values from the tags list
    tags = filter(tags, (val) => {
      return !isEmpty(val);
    });
    //Preprocess tags to send default name/tags in case user doesn't specify
    return map(tags, (val) => {
      val.source_name = val.source_name || (val.source_type === 'Branch' ? SOURCE_NAME_BRANCH : SOURCE_NAME_TAG);
      val.name = val.name || DOCKER_TAG;
      return val;
    });
  },
  _savePushTrigger(e){
    e.preventDefault();
    const newStripped = remove(this.props.newTags, function(item) {
      return !item.toDelete;
    }); // remove 'new' build items that have been deleted
    let tags = this.props.autoBuildStore.build_tags.concat(newStripped);

    this.context.executeAction(saveRepoPushSettings, {
      JWT: this.props.JWT,
      name: this.props.name,
      namespace: this.props.namespace,
      tags: this._preprocessTags(tags)
    });
  },
  _updatePushItem(isNew, index, fieldkey) {
    // action will handle updates for each field of the build_tag item.
    // index comes from mapping through sorted array
    return (e) => {
      e.preventDefault();
      if (fieldkey === 'delete') {
        this.context.executeAction(deletePushTriggerItem, {isNew, index});
      } else {
        this.context.executeAction(updatePushTriggerItem, {
          isNew,
          index,
          fieldkey,
          value: e.target.value});
      }
    };
  },
  _updateSourceType(isNew, index, source_name) {
    return (e) => {
      e.preventDefault();
      if (e.target.value === 'Tag' && source_name === SOURCE_NAME_BRANCH) {
        this.context.executeAction(updatePushTriggerItem, {
          isNew,
          index,
          fieldkey: 'source_name',
          value: SOURCE_NAME_TAG});
      } else if (e.target.value === 'Branch' && source_name === SOURCE_NAME_TAG) {
        this.context.executeAction(updatePushTriggerItem, {
          isNew,
          index,
          fieldkey: 'source_name',
          value: SOURCE_NAME_BRANCH});
      }
      this.context.executeAction(updatePushTriggerItem, {
        isNew,
        index,
        fieldkey: 'source_type',
        value: e.target.value});
    };
  },
  _renderBuildTagRow(isNew) {
    const {
      JWT,
      name,
      namespace
    } = this.props;
    // isNew needed to differentiate between newly added items and tags from api
    let addOrRemove, triggerBuildItem;
    return (buildTag, index) => {
      if (!buildTag || buildTag.toDelete) {
        return (<div/>);
      } else if (index !== 0 || isNew) {
        addOrRemove = (
          <span className={styles.removeBtn} onClick={this._updatePushItem(isNew, index, 'delete')}>
            <FA icon='fa-minus' />
          </span>
        );
      } else if (index === 0 && !isNew) {
        addOrRemove = (
          <span className={styles.addBtn} onClick={this._addPushItem}>
            <FA icon='fa-plus' />
          </span>
        );
      }

      // Set defaults
      let sourceNamePlaceholder;
      let tagNamePlaceholder;

      if (buildTag.source_type === 'Branch') {
        sourceNamePlaceholder = SOURCE_NAME_PLACEHOLDER_BRANCH;
        tagNamePlaceholder = DOCKER_TAG_PLACEHOLDER_BRANCH;
      } else {
        sourceNamePlaceholder = SOURCE_NAME_PLACEHOLDER_TAG;
        tagNamePlaceholder = DOCKER_TAG_PLACEHOLDER_TAG;
      }

      if (buildTag.dockerfile_location === '') {
        buildTag.dockerfile_location = '/';
      }

      let displayedTagName;
      if (buildTag.name === DOCKER_TAG) {
        displayedTagName = '';
      } else {
        displayedTagName = buildTag.name;
      }

      let displayedSourceName;
      let sourceName = SOURCE_NAME_BRANCH;
      if ((buildTag.source_name === SOURCE_NAME_BRANCH && buildTag.source_type === 'Branch') ||
          (buildTag.source_name === SOURCE_NAME_TAG && buildTag.source_type === 'Tag')) {
        displayedSourceName = '';
      } else {
        displayedSourceName = buildTag.source_name;
      }

      let buildTriggerButton = (<div key={index} className={styles.noBtn}/>);
      if (!isNew || buildTag.id) {
        buildTriggerButton = (
          <BuildTriggerByTag isNew={isNew}
                             triggerId={buildTag.id}
                             isRegexRule={_isRegexString(buildTag.source_name)}
                             repoInfo={{JWT, name, namespace}}
                             buildTag={buildTag} />
        );
      }

      return (
        <Row key={index}>
          <Item>
            <select value={buildTag.source_type}
                    onChange={this._updateSourceType(isNew, index, buildTag.source_name)}>
              <option value='Branch'>Branch</option>
              <option value='Tag'>Tag</option>
            </select>
          </Item>
          <Item grow={3}>
            <SimpleInput value={displayedSourceName}
                         placeholder={sourceNamePlaceholder}
                         onChange={this._updatePushItem(isNew, index, 'source_name')}/>
          </Item>
          <Item grow={2}>
            <SimpleInput value={buildTag.dockerfile_location}
                         onChange={this._updatePushItem(isNew, index, 'dockerfile_location')}/>
          </Item>
          <Item grow={2}>
            <SimpleInput value={displayedTagName}
                         placeholder={tagNamePlaceholder}
                         onChange={this._updatePushItem(isNew, index, 'name')}/>
          </Item>
          <Item grow={2}>
            {addOrRemove}
            {buildTriggerButton}
          </Item>
        </Row>
      );
    };
  },
  render() {
    const {
      autoBuildStore,
      JWT,
      name,
      namespace,
      newTags,
      validations
    } = this.props;

    const tags = autoBuildStore.build_tags;
    let tagRows;
    // _renderBuildTagRow sets index in _updatePushItem
    // bool differentiates between new/old build_tags arrays
    const newTagRows = map(newTags, this._renderBuildTagRow(true));
    if (isEmpty(tags) && isEmpty(newTags)) {
      tagRows = (
        <Row>
          <Item>No push triggers to show</Item>
        </Row>
      );
    } else {
      // _renderBuildTagRow sets index in _updatePushItem
      // bool differentiates between new/old build_tags arrays
      tagRows = map(tags, this._renderBuildTagRow(false));
    }

    let saveButtonText = 'Save Changes';
    let saveButtonVariant = 'primary';

    if (validations.buildTags.success) {
      saveButtonText = 'Saved';
      saveButtonVariant = 'success';
    } else if (validations.buildTags.hasError) {
      saveButtonText = 'Error Saving';
      saveButtonVariant = 'alert';
    }

    let maybeErrors = '';
    if (validations.buildTags.errors) {
      maybeErrors = <pre className={styles.clumpedErrors}>{validations.buildTags.errors.join('\n')}</pre>;
    }

      return (
      <form onSubmit={this._savePushTrigger}>
        <div className="row">
          <div className='columns large-12'>
            <Card>
              <Table isWrappedInACard={true}
                         className={styles.table}>
                <Header>
                  <Item header={true}>Type</Item>
                  <Item header={true} grow={3}>Name</Item>
                  <Item header={true} grow={2}>Dockerfile Location</Item>
                  <Item header={true} grow={2}>Docker Tag Name</Item>
                  <Item header={true} grow={2}/>
                </Header>
                { tagRows }
                { newTagRows }
              </Table>
              <div className="row">
                <div className={styles.error}>
                  { maybeErrors }
                </div>
                <div className={styles.button}>
                  <Button variant={saveButtonVariant}
                          type='submit'
                          size='small'>
                    {saveButtonText}
                  </Button>
                </div>
              </div>
            </Card>
          </div>
        </div>
      </form>
    );
  }
});

export default RepoPushTrigger;
