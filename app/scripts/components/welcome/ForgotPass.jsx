'use strict';

import React from 'react';
import { Link } from 'react-router';
import Button from '@dux/element-button';
import FancyInput from 'common/FancyInput';
import forgotPasswordSubmit from '../../actions/forgotPasswordSubmit.js';
import styles from './ForgotPass.css';
var debug = require('debug')('Password Reset');

module.exports = React.createClass({
  contextTypes: {
    executeAction: React.PropTypes.func.isRequired
  },
  getInitialState: function() {
    return {
      email: '',
      reset: false
    };
  },
  onEmailChange: function(e) {
    e.preventDefault();
    this.setState({email: e.target.value});
  },
  forgotPassSubmit: function(e) {
    e.preventDefault();
    debug(this.state.email);
    if (this.state.email) {
      this.context.executeAction(forgotPasswordSubmit, {email: this.state.email});
      this.setState({reset: true});
    }
  },
  transitionHome: function(e) {
    e.preventDefault();
    this.props.history.pushState(null, '/');
  },
  render: function() {
    var disabledState;
    if (this.state.email) {
      disabledState = false;
    } else {
      disabledState = true;
    }

    if (!this.state.reset) {
      return (
        <div className={styles.forgotPassPage}>
          <div className='row'>
            <div className={'small-6 small-centered large-4 large-centered columns ' + styles.header}>
              <Link to='/'>
                <img src="/public/images/logos/mini-logo.svg" alt='docker logo' className={styles.logo}/>
              </Link>
              <div className={styles.head}>Reset your password</div>
              <div className={styles.subHead}>Enter an email address associated with a Docker ID.</div>
            </div>
          </div>
          <div className='row'>
            <div className='small-6 small-centered large-3 large-centered columns'>
              <form onSubmit={this.forgotPassSubmit} className={styles.forgotPassSubmit}>
                <FancyInput placeholder='Email Address'
                  onChange={this.onEmailChange}
                  type="email"
                  value={this.state.email}
                  variant='white'/>
              </form>
            </div>
            <div className={'small-6 small-centered large-4 large-centered columns ' + styles.header}>
              <div className={styles.subHeadLight}>
                We’ll send a password reset link to the Docker ID’s primary email address.
              </div>
              <div className={styles.subHeadLight}>
                If you don't have access to your primary email address, contact <a href="mailto:support@docker.com">Docker Support</a>.
              </div>
              <div className={styles.btn}>
                <Button className={styles.btn} type="submit" disabled={disabledState}>Reset password</Button>
              </div>
            </div>
          </div>
        </div>
      );
    } else {
      return (
        <div className={styles.forgotPassPage}>
          <div className='row'>
            <div className={'small-6 small-centered large-4 large-centered columns ' + styles.header}>
              <Link to='/'>
                <img src="/public/images/logos/mini-logo.svg" alt='docker logo' className={styles.logo}/>
              </Link>
              <div className={styles.head}>Reset request sent!</div>
              <p className={styles.subHead}> Password reset link sent. This link is valid for 24 hours. If you don't see a password reset link, check your spam folder.</p>
              <Button onClick={this.transitionHome}>Back to Home</Button>
            </div>
          </div>
        </div>
      );
    }
  }
});
