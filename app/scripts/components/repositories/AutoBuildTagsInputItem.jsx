'use strict';

import React, {PropTypes} from 'react';
import { Row, Header, Item } from 'common/TagsFlexTable';
import FA from 'common/FontAwesome';
import styles from './AutoBuildTagsInputItem.css';
import Button from '@dux/element-button';
const {func, string} = PropTypes;

const AutoBuildTagsInputItem = React.createClass({
  propTypes: {
    row: string.isRequired,
    sign: string.isRequired,
    tagId: string.isRequired,
    tagName: string.isRequired,
    removeItem: func.isRequired,
    sourceName: string.isRequired,
    sourceType: string.isRequired,
    fileLocation: string.isRequired,
    handleBtnClick: func.isRequired,
    onSourceTypeChange: func.isRequired,
    onSourceNameChange: func.isRequired,
    onTagChange: func.isRequired,
    onDockerfileLocationChange: func.isRequired
  },
  getInitialState: function() {
    return {
      hidden: false,
      sourceType: 'Branch'
    };
  },
  _handleHide: function(e) {
    e.preventDefault();
    e.stopPropagation();
    this.setState({
      hidden: true
    });
    this.props.removeItem(this.props.tagId);
  },
  onSourceTypeChange: function(e) {
    this.setState({
      sourceType: e.target.value
    });
    this.props.onSourceTypeChange(e);
  },
  componentDidMount: function() {
    this.setState({
      tag: this.props.tagName
    });
  },
  render: function() {

    let sourceNamePlaceholder;
    let tagNamePlaceholder;

    //Support to show default rule for each type of rule in the UI
    const branchesOrTags = this.state.sourceType === 'Branch' ? 'branches' : 'tags';
    const branchOrTag = this.state.sourceType === 'Branch' ? 'branch' : 'tag';
    if (this.state.sourceType === 'Branch') {
      sourceNamePlaceholder = 'All branches except master';
    } else {
      sourceNamePlaceholder = '/.*/ This targets all tags';
    }
    tagNamePlaceholder = `Same as ${branchOrTag}`;

    const {
      fileLocation,
      handleBtnClick,
      onDockerfileLocationChange,
      onTagChange,
      onSourceNameChange,
      row,
      sign,
      sourceName,
      sourceType,
      tagName
    } = this.props;

    let item = null;
    if (!this.state.hidden) {
      item = (
      <Row key={this.props.row}>
        <Item grow={1}>
          <select className={styles.select} defaultValue={sourceType} onChange={this.onSourceTypeChange}>
            <option value="Branch">Branch</option>
            <option value="Tag">Tag</option>
          </select>
        </Item>
        <Item grow={3}>
          <input className={styles.rounded}
                 type="text"
                 placeholder={sourceNamePlaceholder}
                 defaultValue={sourceName}
                 onChange={onSourceNameChange} />
        </Item>
        <Item grow={2}>
          <input className={styles.rounded}
                 type="text"
                 defaultValue={fileLocation}
                 onChange={onDockerfileLocationChange}/>
        </Item>
        <Item grow={2}>
          <input className={styles.rounded}
                 type="text"
                 placeholder={tagNamePlaceholder}
                 defaultValue={tagName}
                 onChange={onTagChange} />
        </Item>
        <Item>
          <span id={row} className={sign === 'plus' ? styles.addBtn : styles.removeBtn}
                onClick={sign === 'plus' ? handleBtnClick : this._handleHide}>
            <FA icon={sign === 'plus' ? 'fa-plus' : 'fa-minus'} />
          </span>
        </Item>
      </Row>
      );
    }

    return item;
  }
});

module.exports = AutoBuildTagsInputItem;
