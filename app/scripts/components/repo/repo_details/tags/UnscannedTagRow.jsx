'use strict';
import React, { PropTypes, Component } from 'react';
const { func, object, shape, string, number, bool, instanceOf } = PropTypes;
import DeleteTagArea from './DeleteTagArea';
import { FlexRow, FlexItem } from 'common/FlexTable.jsx';
import bytesToSize from '../../../utils/bytesToSize';
const debug = require('debug')('UnscannedTagRow');
import { StatusRecord } from 'records';
import moment from 'moment';

//Renders a tag row for the Tags table that does not have Nautilus Scan information
export default class UnscannedTagRow extends Component {
  static propTypes = {
    actions: shape({
      deleteRepoTag: func
    }),

    tag: shape({
      name: string.isRequired,
      full_size: number
    }),
    JWT: string.isRequired,
    name: string.isRequired,
    namespace: string.isRequired,
    canEdit: bool.isRequired,
    status: instanceOf(StatusRecord)
  }

  static defaultProps = {
    canEdit: false
  }

  deleteTag = (e) => {
    const { JWT, name, namespace, tag } = this.props;
    const tagName = tag.name;
    this.props.actions.deleteRepoTag({ JWT, namespace, name, tagName });
  }

  render() {
    const {
      tag,
      canEdit,
      namespace,
      name: repoName,
      status
    } = this.props;
    const { name, full_size, last_updated } = tag;

    let deleteArea;
    if (canEdit) {
      deleteArea = (
        <DeleteTagArea
          {...this.props}
          status={status}
          deleteTag={this.deleteTag} />
      );
    }

    const lastUpdated = last_updated ? moment(last_updated).fromNow() : `Unavailable`;

    return (
      <FlexRow>
        <FlexItem grow={2}>{name}</FlexItem>
        <FlexItem>{bytesToSize(full_size)}</FlexItem>
        <FlexItem>{lastUpdated}</FlexItem>
        <FlexItem end>
          {deleteArea}
        </FlexItem>
      </FlexRow>
    );
  }
}
