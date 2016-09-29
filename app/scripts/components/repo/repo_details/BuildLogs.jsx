'use strict';
import React, { PropTypes, Component } from 'react';
import { Link } from 'react-router';
import connectToStores from 'fluxible-addons-react/connectToStores';
import Card, { Block } from '@dux/element-card';
import Code from 'common/Code';
import Markdown from '@dux/element-markdown';
import RepoDetailsBuildLogsStore from 'stores/RepoDetailsBuildLogs';
import endsWith from 'lodash/string/endsWith';
import moment from 'moment';
import {
  FlexTable,
  FlexHeader,
  FlexItem,
  FlexRow
} from 'common/FlexTable';
import {
  FlexTable as Table,
  Row,
  Header,
  Item
} from 'common/TagsFlexTable';
import styles from './BuildLogs.css';
const debug = require('debug')('BuildLogs');

class BuildLogs extends Component {

  _formatDate(ISODateString) {
    //TODO: Remove this if statement once Highland 2.0 works for prod
    //fixes a bug where Highland returns a date without the `Z` UTC indicator
    //only in build logs on prod
    if (!endsWith(ISODateString, 'Z')) {
      return ISODateString ? moment(`${ISODateString}Z`, moment.ISO_8601).fromNow() : '';
    }
    return ISODateString ? moment(ISODateString, moment.ISO_8601).fromNow() : '';
  }

  render() {
    const {
      failure,
      readme_contents,
      created_at,
      build_path,
      docker_tag,
      build_code,
      dockerfile_contents,
      source_branch,
      logs
    } = this.props.build_results;

    let failureAlert = null;
    if(failure) {
      failureAlert = <div className='alert-box alert'>{failure}</div>;
    }
    return (
      <div className="repo-details-blank">
        {failureAlert}
        <Table>
          <Header>
            <Item>SourceRef</Item>
            <Item>Dockerfile Location</Item>
            <Item>Docker Tag</Item>
            <Item>Build Created</Item>
            <Item>UTC</Item>
          </Header>
          <Row>
            <Item>{source_branch}</Item>
            <Item>{build_path}</Item>
            <Item>{docker_tag}</Item>
            <Item>{this._formatDate(created_at)}</Item>
            <Item>{created_at}</Item>
          </Row>
        </Table>
        <Card heading='Build Code'>
          <Block>
            <p className={styles.buildCode}>{build_code}</p>
          </Block>
        </Card>
        <Card heading='README'>
          <Block>
            <Markdown>{readme_contents}</Markdown>
          </Block>
        </Card>
        <Card heading='Dockerfile'>
          <Code>{dockerfile_contents}</Code>
        </Card>
        <Card heading='Logs'>
          <Block>
            <p className={styles.logs}>{logs}</p>
          </Block>
        </Card>
      </div>
    );
  }
}

export default connectToStores(BuildLogs,
                               [RepoDetailsBuildLogsStore],
                               function({ getStore }, props) {
                                 return getStore(RepoDetailsBuildLogsStore).getState();
                               });
