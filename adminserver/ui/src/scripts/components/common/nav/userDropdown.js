'use strict';

import React, { Component, PropTypes } from 'react';
const { object } = PropTypes;
import { connect } from 'react-redux';
import { Link } from 'react-router';
import { currentUserSelector } from 'selectors/users';
import { createStructuredSelector as struct } from 'reselect';
import { mapActions } from 'utils';
import * as authActions from 'actions/auth';
import FontAwesome from 'components/common/fontAwesome';
import styles from './userDropdown.css';

@connect(struct({ user: currentUserSelector }), mapActions(authActions))
export default class UserDropdown extends Component {
  static propTypes = {
    user: object.isRequired,
    actions: object.isRequired
  }

  logout() {
    this.props.actions.logOut();
  }

  render = () => (
    <div id='user-menu' className={ styles.container }>
      <Link to={ '/users/' + this.props.user.name }>
        <span className={ styles.avatar }>
          <FontAwesome icon='fa-user' />
        </span>
        <span id='username' className={ styles.name }>{ this.props.user.name }
          <FontAwesome icon='fa-caret-down' />
        </span>
      </Link>
      <ul className={ styles.dropdown }>
        <li><Link to='/docs/api'>API docs</Link></li>
        <li>
          <a href='https://docs.docker.com/docker-trusted-registry/' target='_blank'>
            Docs <FontAwesome icon='fa-external-link' />
          </a>
        </li>
        <li><Link to='/support'>Support</Link></li>
        <li><a href='#' onClick={ ::this.logout }>Logout</a></li>
      </ul>
    </div>
  )

}
