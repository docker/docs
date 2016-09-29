'use strict';

import React, {
  Component,
  PropTypes
  } from 'react';
const { object, number, string, array, func} = PropTypes;

import _ from 'lodash';
import FA from 'common/FontAwesome';
import Button from '@dux/element-button';
import addCollaborator from 'actions/addCollaborator';
import delCollaborator from 'actions/delCollaborator';
import { STATUS } from 'stores/collaborators/Constants.js';
import onAddCollaboratorChange from 'actions/onAddCollaboratorChange';
import Card, { Block } from '@dux/element-card';
import SimpleInput from 'common/SimpleInput';
import styles from './Collaborators.css';
import {
  FlexTable,
  FlexHeader,
  FlexRow,
  FlexItem
} from 'common/FlexTable';
import Pagination from 'common/Pagination.jsx';

const debug = require('debug')('RepoCollaborators');

class Collaborators extends Component {

  static propTypes = {
    count: number.isRequired,
    next: string,
    previous: string,
    results: array.isRequired,
    requests: object,
    STATUS: string,
    newCollaborator: string,
    namespace: string.isRequired,
    name: string.isRequired,
    JWT: string.isRequired,
    error: string,
    location: object,
    history: object
  }

  static contextTypes = {
    executeAction: func.isRequired
  }

  removeUser = (username) => {
    const { namespace, name, JWT } = this.props;
    return (e) => {
      const args = {
        JWT,
        namespace,
        name,
        username
      };
      this.context.executeAction(delCollaborator, args);
    };
  }

  onSubmit = (e) => {
    e.preventDefault();
    const { namespace, name, JWT } = this.props;
    const args = {
      JWT,
      namespace,
      name,
      user: this.props.newCollaborator
    };
    this.context.executeAction(addCollaborator, args);
  }

  renderTable = () => {
    return (
      <FlexTable>
        <FlexHeader>
          <FlexItem>Username</FlexItem>
          <FlexItem>Access</FlexItem>
          <FlexItem>Action</FlexItem>
        </FlexHeader>
        {this.props.results.map(this.renderTr)}
      </FlexTable>
    );
  }

  renderTr = ({user}) => {
     //TODO: Need better way to do validations for deleting collaborators
    const status = this.props.requests[user];
    const isLoading = status === 'LOADING';
    const hasError = status === 'ERROR';

    let icon = <FA icon='fa-times'/>;
    if(isLoading) {
      icon = <FA icon='fa-spinner' animate='spin'/>;
    }
    return (
      <FlexRow key={user}>
        <FlexItem>{user}</FlexItem>
        <FlexItem>Collaborator</FlexItem>
        <FlexItem>
          <div className={styles.button}>
            <Button variant='alert'
                    size='tiny'
                    disabled={isLoading}
                    onClick={this.removeUser(user)}>
              {icon} Remove
            </Button>
          </div>
        </FlexItem>
      </FlexRow>
    );
  }

  onAddCollaboratorChange = function(e) {
    e.preventDefault();
    this.context.executeAction(onAddCollaboratorChange, { collaborator: e.target.value });
  }

  renderAddCollaboratorForm = (collaborator) => {
    const { error } = this.props;
    const isLoading = this.props.STATUS === STATUS.ATTEMPTING;
    const hasError = this.props.STATUS === STATUS.ERROR && error;
    let variant = 'primary';
    let errorMsg;
    if (hasError) {
      variant = 'alert';
      errorMsg = <div className={styles.error}>{this.props.error}</div>;
    }
    let button = <Button variant={variant} type='submit'>Add User</Button>;
    if(isLoading) {
      button = <Button variant={variant} type='submit' disabled>Adding...</Button>;
    }

    return (
      <Card>
        <Block>
          <form onSubmit={this.onSubmit}>
            <label>Username
              <SimpleInput type='text'
                           onChange={this.onAddCollaboratorChange.bind(this)}
                           value={collaborator}/>
            </label>
            {errorMsg}
            {button}
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
      <div className={'row ' + styles.collaboratorsContent}>
        <div className='large-8 columns'>
          {this.renderTable()}

          <Pagination next={next} prev={prev}
                      onChangePage={this.onChangePage}
                      currentPage={currentPageNumber || 1}
                      pageSize={10} />
        </div>
        <div className='large-4 columns'>
          {this.renderAddCollaboratorForm(this.props.newCollaborator)}
        </div>
      </div>
    );
  }
}

export default Collaborators;
