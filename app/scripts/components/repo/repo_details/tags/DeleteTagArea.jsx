'use strict';

import React, { PropTypes, Component } from 'react';
import { ATTEMPTING, ERROR } from 'reduxConsts';
const { func, instanceOf } = PropTypes;
import styles from './DeleteTagArea.css';
import FA from 'common/FontAwesome';
import { StatusRecord } from 'records';
const debug = require('debug')('hub:deleteTagArea');

export default class DeleteTagArea extends Component {
  // default status comes from StatusRecord
  static propTypes = {
    status: instanceOf(StatusRecord),
    deleteTag: func
  }

  state = {
    confirmingDelete: false
  }

  confirmDelete = (e) => {
    this.setState({
      confirmingDelete: true
    });
  }

  cancelConfirm = () => {
    this.setState({
      confirmingDelete: false
    });
  }

  render() {
    const { deleteTag, status: { status } } = this.props;
    const { confirmingDelete } = this.state;
    let content;
    if (status === ERROR) {
      content = (
        <div className={styles.error}>
          Deletion Failed
        </div>
      );
    } else if (status === ATTEMPTING) {
      content = (
        <div>
          Deleting&nbsp;<FA icon="fa-spinner" animate="spin"/>
        </div>
      );
    } else if (confirmingDelete) {
      content = (
        <div>
          <a className={styles.confirmDelete} onClick={deleteTag}>
            Confirm
          </a> or&nbsp;
          <a onClick={this.cancelConfirm}>
            Cancel
          </a>
        </div>
      );
    } else {
      content = (
        <a className={styles.delete} onClick={this.confirmDelete}>
          <FA icon="fa-trash-o" size="lg"/>
        </a>
      );
    }
    return content;
  }
}
