'use strict';

import React, { PropTypes, Component } from 'react';
const { string } = PropTypes;
import { findDOMNode } from 'react-dom';
import isEqual from 'lodash/lang/isEqual';
import Card, { Block } from '@dux/element-card';
import styles from './PullCommand.css';
const debug = require('debug')('PullCommand');


export default class PullCommand extends Component {
  static propTypes = {
    namespace: string.isRequired,
    name: string.isRequired
  }

  /* When you click on the pull command, automatically select all of it */
  selectPullCommand = (e) => {
    findDOMNode(this.refs.pull).select();
  }

  copyPullCommandToClipboard = (e) => {
    /* Works in Chrome and Firefox 41.x
     * TODO: Does not work in Safari
     */
    this.selectPullCommand(e);
    try {
      const success = document.execCommand('copy');
      debug(`Copy worked: ${success}`);
    } catch (err) {
      debug('Cannot copy.');
    }
  }

  render() {
    const { namespace, name } = this.props;
    const repoName = isEqual(namespace, 'library') ? name : `${namespace}/${name}`;
    return (
      <Card heading='Docker Pull Command'
            headingActions={[{
              key: 'copy',
              icon: 'fa-clipboard',
              action: this.copyPullCommandToClipboard
            }]} >
        <Block>
          <div className={styles.pullCommand}>
            <input className={styles.pullCommand}
                   value={`docker pull ${repoName}`}
                   onClick={this.selectPullCommand}
                   ref="pull"
                   readOnly/>
          </div>
        </Block>
      </Card>
    );
  }
}
