'use strict';

import styles from './Register.css';

import React, { Component, PropTypes } from 'react';
import SignupForm from './welcome/SignupForm.jsx';
var debug = require('debug')('Register');

export default class Register extends Component {
  render() {
    return (
      <header className={styles.header}>
        <div className='row'>

          <div className='large-4 large-centered columns'>
            <SignupForm location={this.props.location}/>
          </div>

        </div>
      </header>
    );
  }
}
