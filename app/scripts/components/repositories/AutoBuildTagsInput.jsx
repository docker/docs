'use strict';

import React from 'react';
import { FlexTable, Header, Item } from 'common/TagsFlexTable';
import AutoBuildTagsInputItem from './AutoBuildTagsInputItem.jsx';

const AutoBuildTagsInput = React.createClass({
  contextTypes: {
    getStore: React.PropTypes.func.isRequired
  },
  propTypes: {
    repo: React.PropTypes.string.isRequired,
    onTagRemoved: React.PropTypes.func.isRequired,
    onTagAdded: React.PropTypes.func.isRequired,
    onSourceTypeChange: React.PropTypes.func.isRequired,
    onSourceNameChange: React.PropTypes.func.isRequired,
    onDockerfileLocationChange: React.PropTypes.func.isRequired
  },
  getInitialState: function() {
    const masterDefaultRow = {
      row: 'row-0',
      sourceType: 'Branch',
      sourceName: 'master',
      sign: 'plus',
      tag: 'latest',
      fileLocation: '/',
      tagId: 'tag-0'
    };
    const dynamicTagRow = {
      row: 'row-1',
      sourceType: 'Branch',
      sourceName: '',
      sign: 'minus',
      tag: '',
      fileLocation: '/',
      tagId: 'tag-1'
    };
    return {
      currentRows: [
        this.makeRow(masterDefaultRow),
        this.makeRow(dynamicTagRow)
      ]
    };
  },
  _handleBtnClickAdd: function(e) {
    e.preventDefault();
    let rows = this.state.currentRows;
    let idx = this.state.currentRows.length;
    const row = {
      row: 'row-' + idx,
      sourceType: 'Branch',
      sourceName: '',
      sign: 'minus',
      tag: '',
      fileLocation: '/',
      tagId: 'tag-' + idx
    };
    const newRow = this.makeRow(row);
    this.setState({
      currentRows: rows.concat(newRow)
    });
    //Make an empty build tag that will get updated on edit
    const bTag = {
      id: 'tag-' + idx,
      sourceType: 'Branch',
      sourceName: '',
      fileLocation: '/',
      tagName: ''
    };
    this.props.onTagAdded(this.makeBuildTag(bTag));
  },
  makeRow: function(rowObj) {
    const {
      row,
      sourceType,
      sourceName,
      sign,
      tag,
      fileLocation,
      tagId
    } = rowObj;

    return (
      <AutoBuildTagsInputItem row={row}
                              key={row}
                              repoName={this.props.repo}
                              tagId={tagId}
                              sourceName={sourceName}
                              sign={sign}
                              tagName={tag}
                              sourceType={sourceType}
                              fileLocation={fileLocation}
                              handleBtnClick={this._handleBtnClickAdd}
                              removeItem={this.props.onTagRemoved.bind(null, tagId)}
                              onSourceTypeChange={this.props.onSourceTypeChange.bind(null, tagId)}
                              onSourceNameChange={this.props.onSourceNameChange.bind(null, tagId)}
                              onDockerfileLocationChange={this.props.onDockerfileLocationChange.bind(null, tagId)}
                              onTagChange={this.props.onTagChange.bind(null, tagId)}/>
    );
  },
  makeBuildTag: function(newBuildTag) {
    const {
      id,
      sourceType,
      sourceName,
      fileLocation,
      tagName
    } = newBuildTag;
    /*eslint-disable camelcase */
    /* By default, if the input is empty for source tag/branch name, send the string: '{sourceref}'*/
    /* By default, if the input is empty for docker tag name, send the string with the regex for all matches */
    return {
      id: id,
      name: tagName === '' ? '{sourceref}' : tagName,
      source_type: sourceType,
      source_name: sourceName === '' ? '/^([^m]|.[^a]|..[^s]|...[^t]|....[^e]|.....[^r]|.{0,5}$|.{7,})/' : sourceName,
      dockerfile_location: fileLocation
    };
    /*eslint-enable camelcase */
  },
  componentDidMount: function() {
    const masterBuildTag = {
      id: 'tag-0',
      sourceType: 'Branch',
      sourceName: 'master',
      fileLocation: '/',
      tagName: 'latest'
    };
    const dynamicBuildTag = {
      id: 'tag-1',
      sourceType: 'Branch',
      sourceName: '',
      fileLocation: '/',
      tagName: ''
    };
    this.props.onTagAdded(
      this.makeBuildTag(masterBuildTag)
    );
    this.props.onTagAdded(
      this.makeBuildTag(dynamicBuildTag)
    );
  },
  render: function() {
    return (
      <FlexTable>
        <Header>
          <Item header={true}>Push Type</Item>
          <Item header={true} grow={3}>Name</Item>
          <Item header={true} grow={2}>Dockerfile Location</Item>
          <Item header={true} grow={2}>Docker Tag</Item>
          <Item />
        </Header>
        { this.state.currentRows }
      </FlexTable>
    );
  }
});

module.exports = AutoBuildTagsInput;
