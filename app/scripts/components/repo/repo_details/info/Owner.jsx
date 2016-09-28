'use strict';

import React, { PropTypes, Component } from 'react';
const { string } = PropTypes;
import Card, { Block } from '@dux/element-card';
import isEqual from 'lodash/lang/isEqual';
import { mkAvatarForNamespace } from 'utils/avatar';
import styles from './Owner.css';

export default class Owner extends Component {
  static propTypes = {
    namespace: string.isRequired
  }

  render() {
    const { namespace } = this.props;
    if(isEqual(namespace, 'library')) {
      return null;
    } else {
      const avatar = mkAvatarForNamespace(namespace);
      return (
        <Card heading='Owner'>
          <Block>
            <div className={styles.detail}>
              <img src={avatar} className={styles.avatar} />
              <div className={styles.username}>{namespace}</div>
            </div>
          </Block>
        </Card>
      );
    }
  }
}
