'use strict';

import React, { Component, PropTypes } from 'react';
const { bool } = PropTypes;
import { Link } from 'react-router';
import Logo from 'components/common/logo';

import Search from './search.js';
import UserDropdown from './userDropdown.js';
import styles from './nav.css';

export default class Nav extends Component {

  static propTypes = {
    isLoggedIn: bool,
    navExpanded: bool
  }

  render() {
    const {
      isLoggedIn,
      navExpanded
    } = this.props;

    return (
      <nav className={ navExpanded ? styles.navbarExpanded : styles.navbar }>
        <div className={ styles.logoContainer }>
          <Link to='/'>
            <div className={ styles.logo }>
              <Logo scale={ 0.8 } pathClassName={ styles.logoFill } />
            </div>
            <span className={ styles.text }> Trusted Registry</span>
          </Link>
        </div>

        <div className={ styles.searchContainer }>
          { isLoggedIn ? <Search /> : undefined }
        </div>
        <div className={ styles.userContainer }>
          { isLoggedIn ? <UserDropdown /> : undefined }
        </div>
      </nav>
    );
  }
}
