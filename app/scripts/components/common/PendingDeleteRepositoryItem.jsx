'use strict';

import styles from './PendingDeleteRepositoryItem.css';
import React, { PropTypes, Component } from 'react';
let { func, string } = PropTypes;
import classnames from 'classnames';
import FA from 'common/FontAwesome';
import { mkAvatarForNamespace, isOfficialAvatarURL } from 'utils/avatar';

/**
 * PendingDeleteRepositoryItem will show only the name and namespace
 */

export default class PendingDeleteRepositoryItem extends Component {

  static propTypes = {
    namespace: string.isRequired,
    name: string.isRequired
  }

  render() {

    const {
      name,
      namespace
      } = this.props;

    const repoDisplayName = namespace + '/' + name;
    const avatar = mkAvatarForNamespace(namespace, name);
    const avatarClass = classnames({
      [styles.avatar]: true,
      [styles.officialAvatar]: isOfficialAvatarURL(avatar)
    });

    return (
      <li key={repoDisplayName}
          className={styles.pendingDeleteItem}
          title={`A request has been made to delete this repository and is in progress at the moment.
You cannot push to/pull this repository and autobuilds setup on this repository will be inactive.`}>
        <div className={styles.flexible}>
          <div className={styles.head}>
            <div className={avatarClass}><img src={avatar} /></div>
            <div className={styles.title}>
              <div className={styles.labels}>
                <div className={styles.repoName}>{ repoDisplayName }</div>
                <span className={styles.pendingDelete}><FA icon='fa-spin fa-spinner' /> Deleting...</span>
              </div>
            </div>
          </div>
          <div className={styles.action}>
            <FA icon='fa-info' size='lg'/>
            <div className={styles.text}>DETAILS</div>
          </div>
        </div>
      </li>
    );
  }
}
