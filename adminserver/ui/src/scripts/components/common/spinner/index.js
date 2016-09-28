'use strict';

import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import Immutable from 'immutable';
const { any, array, bool, instanceOf, string, oneOfType } = PropTypes;
import FontAwesome from 'components/common/fontAwesome';
import consts from 'consts';
import styles from './styles.css';

const { loading } = consts;

const mapState = (state) => { return { status: state.status }; };
@connect(mapState)
export default class Spinner extends Component {
  static propTypes = {
    className: string,
    // This can either be a status string or an array of statuses or an array of promise key paths
    loadingStatus: oneOfType([array, string]),
    children: any,
    showChildrenWhileSpinning: bool,
    status: instanceOf(Immutable.Map),
    text: string
  }

  static defaultProps = {
    showChildrenWhileSpinning: false
  }

  render() {
    const { children, className, loadingStatus, showChildrenWhileSpinning, text } = this.props;
    // Normalize loadingStatus into an array
    const statuses = !Array.isArray(loadingStatus) ? [loadingStatus] : loadingStatus;
    // For each status in statuses, if it is an array we need to read it from this.props.status
    const statusStrings = statuses.map((keypathOrStatus) => {
      if (Array.isArray(keypathOrStatus)) {
        return this.props.status.getIn([...keypathOrStatus, 'status']);
      }
      return keypathOrStatus;
    });
    const errorMessages = statuses.map((keypathOrStatus) => {
      if (Array.isArray(keypathOrStatus)) {
        const error = this.props.status.getIn([...keypathOrStatus, 'error']);
        if (error && error.data && error.data.errors) {
          return <p>{ error.data.errors[0].detail }</p>;
        }
      }
    }).filter(Boolean);
    const shouldSpinner = statusStrings.some(status => (status === loading.PENDING));
    return (
      <div className={ className }>
        { (!shouldSpinner || showChildrenWhileSpinning) &&
          (
            errorMessages.length ?
              errorMessages[0]
            :
              children
          )
        }
        { shouldSpinner &&
          <div className={ styles.spinner }>
            <FontAwesome icon='fa-spinner' animate='spin' />
            <p>{ text || 'Loading...' }</p>
          </div>
        }
      </div>
    );
  }
}
