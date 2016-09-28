'use strict';

import styles from './OrganizationSummary.css';
import classnames from 'classnames';
import React, { Component, PropTypes } from 'react';
const { bool, string } = PropTypes;
import { Link } from 'react-router';
import selectOrganizationAction from '../../actions/selectOrganization';
import { mkAvatarForNamespace } from 'utils/avatar';

const debug = require('debug')('OrganizationSummary');


class OrgItem extends Component {

  static propTypes = {
    orgname: string.isRequired,
    end: bool
  }

  render() {
    const { orgname, end } = this.props;
    const classes = classnames(
      styles.orgListItem,
      'medium-3',
      'columns',
      {
        'end': end
      }
    );
    return (
      <li className={classes}>
        <Link to={`/u/${orgname}/dashboard/`}>
          <div className={styles.orgGridItem}>
            <img className={styles.orgAvatar} src={mkAvatarForNamespace(orgname)} />
            <h6 className={styles.orgSummary}>{orgname}</h6>
          </div>
        </Link>
      </li>
    );
  }
}

export default class OrganizationSummary extends Component {
  static contextTypes = {
    executeAction: React.PropTypes.func.isRequired
  }
  static propTypes = {
    user: React.PropTypes.object.isRequired,
    JWT: React.PropTypes.string.isRequired,
    orgs: React.PropTypes.array.isRequired
  }
  render() {
    const orgItems = this.props.orgs.map((org, i, array) => { return <OrgItem {...org} key={org.orgname} end={i === array.length - 1}/>; });

    return (
      <ul className="row">
        {orgItems}
      </ul>
    );
  }
}
