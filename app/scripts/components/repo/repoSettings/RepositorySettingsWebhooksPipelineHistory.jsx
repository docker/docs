'use strict';

import React, { PropTypes, Component } from 'react';
import moment from 'moment';
import _ from 'lodash';
import Card, { Block } from '@dux/element-card';
const debug = require('debug')('RepositorySettingsWebhooksPipelineHistory');

const states = [
  'initialized',
  'pending',
  'success',
  'error',
  'failure'
];

export default class PipelineHistory extends Component {
  render() {
    const { historyHistory, params } = this.props;
    debug('historyHistory', historyHistory);
    const hook_id = _.has(params, 'hook') ? params.hook : '';
    return (
      <Card heading={`Hook History: ${hook_id}`}>
        <Block>
          <table>
            <tr>
              <th>id</th>
              <th>state</th>
              <th>reply_ip</th>
              <th>last_updated</th>
            </tr>
            {historyHistory.results.map(this.renderItem)}
          </table>
        </Block>
      </Card>
    );
  }

  renderItem = ({ id, reply_ip, last_updated, state }) => {
    return (
      <tr key={id}>
        <td>{id}</td>
        <td>{state}</td>
        <td>{reply_ip}</td>
        <td>{moment.utc(last_updated).fromNow()}</td>
      </tr>
    );
  }
}
