'use strict';
import React, { PropTypes, Component } from 'react';
import connectToStores from 'fluxible-addons-react/connectToStores';
import RepoDetailsBuildsStore from '../../../stores/RepoDetailsBuildsStore';
import BuildStatus from './BuildStatus';
import SourceRepositoryCard from 'common/SourceRepositoryCard';
import moment from 'moment';
import trunc from 'lodash/string/trunc';
import has from 'lodash/object/has';
import {
  FlexTable,
  FlexHeader,
  FlexRow,
  FlexItem
} from 'common/FlexTable';
const { func } = PropTypes;
import styles from './BuildDetails.css';
import { Button } from 'dux';
import cancelBuildAction from 'actions/cancelBuild';

const debug = require('debug')('RepositoryDetailsTags');

class TableItem extends Component {

  isCancelable(status) {
    /* terminal states
    const terminalStates = {
      -4: ‘canceled',
      -2: 'exception',
      -1: 'error',
      10: 'done',
    };
    */
    const cancelableStates = {
      0: 'pending',
      1: 'claimed',
      2: 'started',
      3: 'cloned',
      4: 'readme',
      5: 'dockerfile',
      6: 'built',
      7: 'bundled',
      8: 'uploaded',
      9: 'pushed',
      11: 'queued'
    };

    return has(cancelableStates, status);
  }

  _formatDate(ISODateString) {
    return ISODateString ? moment(ISODateString, moment.ISO_8601).fromNow() : '';
  }

  renderBuildStatus() {
    return (
      <span>
        <BuildStatus status={this.props.status} />
      </span>
    );
  }

  renderCancelBtn() {
    const { status, onCancel, canceling } = this.props;
    const stopPropagation = (e) => e.stopPropagation();

    // hide display cancelation btn if status isn't cancelable
    if (!this.isCancelable(status)) {
      return null;
   }

    // display 'Canceling' if build cancelation api call successful
    if (canceling === 'success') {
      return (
        <button
          className={`${styles.actionBtn} ${styles.disabled} button tiny`}
          onClick={stopPropagation}
        >Canceling</button>
      );
    }

    // display cancelation button
    return (
      <button
        className={`${styles.actionBtn} button tiny`}
        onClick={onCancel}
      >Cancel</button>
    );
  }

  showBuildLogs = (e) => {
    const { build_code, params, history } = this.props;
    const { user, splat } = params;
    history.pushState(null, `/r/${user}/${splat}/builds/${build_code}/`);
  }

  render() {
    debug(this.props);
    const { created_date, last_updated, dockertag_name } = this.props;

    return (
      <FlexRow onClick={this.showBuildLogs} selectable={true} className={styles.row}>
        <FlexItem>
          {this.renderBuildStatus()}
        </FlexItem>
        <FlexItem>
          {this.renderCancelBtn()}
        </FlexItem>
        <FlexItem grow={2}>{dockertag_name}</FlexItem>
        <FlexItem>{this._formatDate(created_date)}</FlexItem>
        <FlexItem>{this._formatDate(last_updated)}</FlexItem>
      </FlexRow>
    );
  }
}

class BuildDetails extends Component {

  _mkItem = (item) => {
    const onClick = (e) => {
      const { JWT, name, namespace, context } = this.props;
      const { build_code, id } = item;
      e.stopPropagation();
      context.executeAction(cancelBuildAction, { JWT, id, name, namespace, build_code });
    };

    return (
      <TableItem
        key={item.build_code}
        params={this.props.params}
        history={this.props.history}
        onCancel={onClick}
        {...item}/>
    );
  }

  render() {
    const { provider, repo_web_url } = this.props.autoBuildStore;
    return (
      <div className='row'>
        <div className='large-8 columns'>
          <div className="repo-details-blank">
            <FlexTable>
              <FlexHeader>
                <FlexItem>Status</FlexItem>
                <FlexItem>Actions</FlexItem>
                <FlexItem grow={2}>Tag</FlexItem>
                <FlexItem>Created</FlexItem>
                <FlexItem>Last Updated</FlexItem>
              </FlexHeader>
              {this.props.results.map(this._mkItem)}
            </FlexTable>
          </div>
        </div>
        <div className='large-4 columns'>
          <SourceRepositoryCard provider={provider} url={repo_web_url} />
        </div>
      </div>
    );
  }
}

export default connectToStores(
  BuildDetails,
  [RepoDetailsBuildsStore],
  ({ getStore }, props) => getStore(RepoDetailsBuildsStore).getState()
);
