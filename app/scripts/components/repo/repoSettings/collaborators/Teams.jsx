'use strict';

import React, {
  Component,
  PropTypes
} from 'react';
import { findDOMNode } from 'react-dom';
const { number, string, array, arrayOf, object, shape, func } = PropTypes;

import _ from 'lodash';
import FA from 'common/FontAwesome';
import Button from '@dux/element-button';
import Card, { Block } from '@dux/element-card';
import addTeamCollaborators from 'actions/addTeamCollaborator';
import delTeamCollaborator from 'actions/delTeamCollaborator';
import {STATUS} from 'stores/collaborators/Constants.js';
import styles from './Teams.css';
import {
  FlexTable,
  FlexHeader,
  FlexRow,
  FlexItem
} from 'common/FlexTable';
import Pagination from 'common/Pagination.jsx';

const debug = require('debug')('RepoCollaborators');

class Teams extends Component {

  static propTypes = {
    count: number.isRequired,
    next: string,
    previous: string,
    results: array.isRequired,
    allTeams: arrayOf(
      shape({
        group_id: number,
        name: string
      })
    ),
    JWT: string.isRequired,
    name: string.isRequired,
    namespace: string.isRequired,
    requests: object,
    location: object,
    history: object
  }

  static contextTypes = {
    executeAction: func.isRequired
  }

  removeTeam = (group_id) => {
    const { namespace, name, JWT } = this.props;
    return (e) => {
      debug('Remove Team ', group_id, 'from', `${namespace}/${name}`);
      const args = {
        JWT,
        namespace,
        name,
        group_id
      };
      this.context.executeAction(delTeamCollaborator, args);
    };
  }

  onSubmit = (e) => {
    e.preventDefault();
    const { JWT, namespace, name, allTeams } = this.props;
    const { permissions, team } = this.refs;
    const permission = findDOMNode(permissions).value;
    const teamName = findDOMNode(team).value;
    const { id } = _.find(allTeams, { name: teamName });
    const args = {
      JWT,
      namespace,
      name,
      id,
      permission
    };
    this.context.executeAction(addTeamCollaborators, args);
  }

  renderTable = () => {
    return (
      <FlexTable>
        <FlexHeader>
          <FlexItem>Team</FlexItem>
          <FlexItem>Access</FlexItem>
          <FlexItem>Action</FlexItem>
        </FlexHeader>
        {this.props.results.map(this.renderTr)}
      </FlexTable>
    );
  }

  renderTr = ({ group_id, group_name, permission }) => {
    //Need better way to do validations for deleting collaborators
    const status = this.props.requests[group_id];
    const isLoading = status === STATUS.ATTEMPTING;
    const hasError = status === STATUS.ERROR;

    let icon = <FA icon='fa-times'/>;
    if(isLoading) {
      icon = <FA icon='fa-circle-o-notch' animate='spin'/>;
    }
    return (
      <FlexRow key={group_name}>
        <FlexItem>{group_name}</FlexItem>
        <FlexItem>{permission}</FlexItem>
        <FlexItem>
          <div className={styles.button}>
            <Button variant='alert'
                    size='tiny'
                    disabled={isLoading}
                    onClick={this.removeTeam(group_id)}>
              {icon} Remove
            </Button>
          </div>
        </FlexItem>
      </FlexRow>
    );
  }

  renderAddTeamForm = () => {
    const { allTeams } = this.props;
    return (
      <Card>
        <Block>
          <form onSubmit={this.onSubmit}>
            <select ref='team'>
              {allTeams.map((team) => <option key={team.name} value={team.group_id}>{team.name}</option>)}
            </select>
            <select ref='permissions'>
              <option value='read'>Read</option>
              <option value='write'>Write</option>
              <option value='admin'>Admin</option>
            </select>
            <Button type='submit'>Add Team</Button>
          </form>
        </Block>
      </Card>
    );
  }

  onChangePage = (pageNumber) => {
    debug('onChangePage', pageNumber);
    this.props.history.pushState(null, this.props.location.pathname, {page: pageNumber});
  }

  render() {
    const { next, prev } = this.props;
    const currentPageNumber = parseInt(this.props.location.query.page, 10);

    return (
      <div className={styles.teamsWrapper}>
        <div className='row'>
          <div className='large-8 columns'>
            {this.renderTable()}

            <Pagination next={next} prev={prev}
                        onChangePage={this.onChangePage}
                        currentPage={currentPageNumber || 1}
                        pageSize={10} />
          </div>
          <div className='large-4 columns'>
            {this.renderAddTeamForm()}
          </div>
        </div>
      </div>
    );
  }
}

export default Teams;
