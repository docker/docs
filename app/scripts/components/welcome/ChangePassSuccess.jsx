'use strict';

import React, { PropTypes } from 'react';
import { Button } from 'dux';
import { Link } from 'react-router';
import ChangePasswordStore from '../../stores/ChangePasswordStore.js';
import styles from './ChangePassSuccess.css';

var changePassSuccess = React.createClass({
  contextTypes: {
    executeAction: PropTypes.func.isRequired
  },
  PropTypes: {
    changePassStore: PropTypes.shape({
      reset: PropTypes.bool.isRequired,
      resetErr: PropTypes.bool.isRequired
    })
  },
  render: function() {
    return (
      <div className={'columns large-4 large-centered ' + styles.passwordReset}>
        <h3>Your password has been reset</h3>
        <p>You may now login with your new password</p>
        <Link to="/login/" className="resetPassBack columns large-5">Back to Login</Link>
      </div>
    );
  }
});

module.exports = changePassSuccess;
