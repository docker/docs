'use strict';
import React, { Component, PropTypes } from 'react';
import moment from 'moment';
import getPipelineHistory from 'actions/getPipelineHistory';

import Status from './Status.jsx';
import styles from './History.css';

export default class PipelineHistory extends Component {
  static propTypes = {
    namespace: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired,
    slug: PropTypes.string.isRequired,
    attempts: PropTypes.arrayOf(PropTypes.shape({
      created: PropTypes.string,
      status: PropTypes.string,
      uuid: PropTypes.string
    }))
  }

  static contextTypes = {
    executeAction: PropTypes.func.isRequired
  }

  state = {
    loading: true
  }

  componentDidMount() {
    const { JWT: jwt, namespace, name, slug } = this.props;
    this.context.executeAction(getPipelineHistory, {
      jwt, namespace, name, slug
    });
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.attempts) {
      this.setState({loading: false});
    }
  }

  renderPipelineAttempt({ created, status, uuid }) {
    return (
      <tr key={uuid} className={styles.attempt}>
        <td className={styles.uuid}>
          {uuid}
        </td>
        <td>
          <Status status={status} />
        </td>
        <td className={styles.created}>
          {moment(created).format('MM/DD/YY HH:mm')}
        </td>
      </tr>
    );
  }

  renderEmptyMessage() {
    return (
      <p className={styles.empty}>No history for this webhook</p>
    );
  }

  renderLoadingMessage() {
    return (
      <p className={styles.loading}>Loading history...</p>
    );
  }

  renderLoadingOrHistoryOrEmptyMessage() {
    if (this.state.loading) {
      return this.renderLoadingMessage();
    } else if (this.props.attempts && this.props.attempts.length) {
      return (
        <table className={styles.attempts}>
          <thead>
            <tr>
              <th>ID</th>
              <th>Status</th>
              <th>Date &amp; Time</th>
            </tr>
          </thead>
          <tbody>
            {this.props.attempts.map(this.renderPipelineAttempt, this)}
          </tbody>
        </table>
      );
    } else {
      return this.renderEmptyMessage();
    }
  }

  render() {
    return (
      <div className={styles.history}>
        <div className={styles.title}>History</div>
        {this.renderLoadingOrHistoryOrEmptyMessage()}
      </div>
    );
  }
}
