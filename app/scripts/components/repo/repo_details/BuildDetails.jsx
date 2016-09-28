'use strict';
import React, { PropTypes, Component } from 'react';
import connectToStores from 'fluxible-addons-react/connectToStores';
import RepoDetailsBuildsStore from '../../../stores/RepoDetailsBuildsStore';
import BuildStatus from './BuildStatus';
import SourceRepositoryCard from 'common/SourceRepositoryCard';
import moment from 'moment';
import trunc from 'lodash/string/trunc';
import {
  FlexTable,
  FlexHeader,
  FlexRow,
  FlexItem
} from 'common/FlexTable';
const { func } = PropTypes;
import styles from './BuildDetails.css';
const debug = require('debug')('RepositoryDetailsTags');

class TableItem extends Component {

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

  showBuildLogs = (e) => {
    const { build_code, params, history } = this.props;
    const { user, splat } = params;
    history.pushState(null, `/r/${user}/${splat}/builds/${build_code}/`);
  };

  render() {
    debug(this.props);
    const { created_date, last_updated, dockertag_name } = this.props;
    return (
      <FlexRow onClick={this.showBuildLogs} selectable={true}>
        <FlexItem>
          {this.renderBuildStatus()}
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
    return (
      <TableItem
        key={item.build_code}
        params={this.props.params}
        history={this.props.history}
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

export default connectToStores(BuildDetails,
                               [RepoDetailsBuildsStore],
                               function({ getStore }, props) {
                                 return getStore(RepoDetailsBuildsStore).getState();
                               });
