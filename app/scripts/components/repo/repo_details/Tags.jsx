'use strict';
import React, { PropTypes, Component } from 'react';
const { array, object, bool, string, number, shape, func } = PropTypes;
import { connect } from 'react-redux';
import Card, { Block } from '@dux/element-card';

import { FlexTable, FlexRow, FlexHeader, FlexItem } from 'common/FlexTable.jsx';
import ScannedTagRow from './tags/ScannedTagRow.jsx';
import UnscannedTagRow from './tags/UnscannedTagRow.jsx';
import styles from './Tags.css';
import FontAwesome from 'common/FontAwesome';
import Tooltip from 'rc-tooltip';
import { createStructuredSelector } from 'reselect';
import {
  getScannedTags,
  getScannedTagCount,
  getUnscannedTags,
  getUnscannedTagCount
} from './tags/selectors';
import { getStatus } from 'selectors/status';
import * as tagActions from 'actions/redux/tags.js';
import { mapActions } from 'reduxUtils';
import Button from '@dux/element-button';
import { StatusRecord } from 'records';
import moment from 'moment';

const debug = require('debug')('RepositoryDetailsTags');

const mapState = createStructuredSelector({
  scannedTags: getScannedTags,
  scannedTagCount: getScannedTagCount,
  unscannedTags: getUnscannedTags,
  unscannedTagCount: getUnscannedTagCount,
  status: getStatus
});

/**
 * TagDisplay is the new UI for listing tags with vulnerability information from
 * nautilus.
 * It connects to the redux store and uses redux actions.
 */
@connect(mapState, mapActions(tagActions))
class TagDisplay extends Component {

  static propTypes = {
    actions: shape({
      deleteRepoTag: func
    }),
    status: object,

    scannedTags: array,
    scannedTagCount: number,
    unscannedTags: array,
    unscannedTagCount: number
  }

  state = {
    //one of 'unknown', 'show'
    showUnscannedTags: 'unknown'
  }

  toggleShowUnscannedTags = (e) => {
    const { showUnscannedTags } = this.state;
    //will be 'unknown' on the first time clicking Show Outdated Tags
    this.setState({
      showUnscannedTags: 'show'
    });
  }

  mkUnscannedTagRow = (tag) => {
    const { status } = this.props;
    const tagName = tag.name;
    const tagStatus = status.getIn(['deleteRepoTag', tagName], new StatusRecord());
    return (
      <UnscannedTagRow
        {...this.props}
        tag={tag}
        key={tagName}
        status={tagStatus} />
    );
  }

  mkUnscannedTagTable = () => {
    const { unscannedTags, unscannedTagCount, scannedTagCount } = this.props;
    const { showUnscannedTags } = this.state;
    if (!unscannedTagCount) {
      return null;
    }
    //Nautilus scan results exist --> show the button instead of the table
    if (showUnscannedTags === 'unknown' && scannedTagCount) {
      return (
        <div className={`columns large-12 ${styles.toggleButtonWrapper}`}>
          <Button onClick={this.toggleShowUnscannedTags} ghost>View Unscanned Tags</Button>
        </div>
      );
    }
    //Nautilus scan results exist and the button has been pressed to show unscanned tags
    if (showUnscannedTags === 'show' && scannedTagCount) {
      return (
        <FlexTable>
          <FlexHeader>
            <FlexItem>
              <div>
                <span className={styles.cardHeader}>Unscanned Tags</span>
              </div>
            </FlexItem>
          </FlexHeader>
          <FlexHeader>
            <FlexItem grow={2} noPadding><div className={styles.secondaryTableHeader}>Tag Name</div></FlexItem>
            <FlexItem noPadding><div className={styles.secondaryTableHeader}>Compressed Size</div></FlexItem>
            <FlexItem noPadding><div className={styles.secondaryTableHeader}>Last Updated</div></FlexItem>
            <FlexItem />
          </FlexHeader>
          {unscannedTags.map(this.mkUnscannedTagRow)}
        </FlexTable>
      );

    }
    return (
      <FlexTable>
        <FlexHeader>
          <FlexItem grow={2}><div className={styles.header}>Tag Name</div></FlexItem>
          <FlexItem><div className={styles.header}>Compressed Size</div></FlexItem>
          <FlexItem><div className={styles.header}>Last Updated</div></FlexItem>
          <FlexItem />
        </FlexHeader>
        {unscannedTags.map(this.mkUnscannedTagRow)}
      </FlexTable>
    );
  }

  mkScannedTagRow = (tag) => {
    const { status } = this.props;
    const tagName = tag.name;
    const tagStatus = status.getIn(['deleteRepoTag', tagName], new StatusRecord());
    return (
      <ScannedTagRow
        {...this.props}
        tag={tag}
        key={tagName}
        status={tagStatus} />
    );
  }

  mkScannedTagTable = () => {
    const { scannedTags, scannedTagCount } = this.props;
    if (!scannedTagCount) {
      return null;
    }
    const tooltipTitle = `Scanned Images`;
    const tooltipText = `Docker supports and maintains a set of tags that are recommended for use based on the most secure and quality images. To view the rest of the tag instances for this repository, view "Unscanned Images".`;
    const tooltipContent = (
      <div style={{width: 300}}>
        <div className={styles.tooltipTitle}>{tooltipTitle}</div>
        <div>{tooltipText}</div>
      </div>
    );
    const questionMark = (
      <Tooltip overlay={ tooltipContent }
        placement='bottom'
        align={ { overflow: { adjustY: 0 } } }
        trigger={ ['click'] }>
        <div className={styles.questionMark}>
          <FontAwesome icon='fa-question-circle' />
        </div>
      </Tooltip>
    );
    //sort scannedTags by `created_at` date, from most recent to oldest
    scannedTags.sort((tag1, tag2) => {
      if (moment(tag1.created_at).isSame(tag2.created_at)) {
        return moment(tag1.last_updated).isAfter(tag2.last_updated) ? -1 : 1;
      }
      return moment(tag1.created_at).isAfter(tag2.created_at) ? -1 : 1;
    });
    return (
      <FlexTable>
        <FlexHeader>
          <FlexItem>
            <div>
              <span className={styles.cardHeader}>Scanned Images &nbsp;</span>
              {questionMark}
            </div>
          </FlexItem>
        </FlexHeader>
        {scannedTags.map(this.mkScannedTagRow)}
      </FlexTable>
    );
  }

  render() {
    const { scannedTagCount, unscannedTagCount } = this.props;

    if (!scannedTagCount && !unscannedTagCount) {
      return (
        <Card heading='Tags'>
          <Block>
            <div className={styles.empty}>No tags for this repository.</div>
          </Block>
        </Card>
      );
    }

    return (
      <div>
        { this.mkScannedTagTable() }
        { this.mkUnscannedTagTable() }
      </div>
    );
  }
}

export default TagDisplay;
