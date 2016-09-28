'use strict';

// CSS
import styles from './Login.css';

import React, { createClass, PropTypes } from 'react';
import { Link } from 'react-router';
const { func } = PropTypes;
import LoginForm from './welcome/LoginForm.jsx';
import clearLoginForm from '../actions/clearLoginForm.js';
import attemptLoginAction from '../actions/attemptLogin';
var debug = require('debug')('LoginPage');

export default createClass({
  displayName: 'LoginPage',
  contextTypes: {
    executeAction: func.isRequired
  },
  statics: {
    willTransitionFrom(transition, component) {
      debug('transitioning', component);
      component.context.executeAction(clearLoginForm);
    }
  },
  render() {

    return (
      <div className={styles.loginPage}>
        <header className={styles.header}>
          <div className='row'>
            <div className={'small-10 small-centered large-5 large-centered columns'}>
              <Link to='/'>
                <img src="/public/images/logos/mini-logo.svg" alt='docker logo' className={styles.logo}/>
              </Link>
              <div className={styles.head}>Welcome to Docker Hub</div>
              <div className={styles.subHead}>Login with your Docker ID</div>
            </div>
          </div>
          <div className='row'>
            <div className='small-6 small-centered large-3 large-centered columns'>
                <LoginForm loginAction={attemptLoginAction}
                           variant='white'/>
            </div>
          </div>
        </header>
        <div className={styles.footer}>
          <Link to='/reset-password/'>Can't Login?</Link> | <Link to='/register/'>Create Account</Link>
        </div>
      </div>
    );
  }
});
