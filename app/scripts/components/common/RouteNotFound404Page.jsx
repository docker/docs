'use strict';

import React from 'react';
import { Module } from 'dux';
import styles from './RouteNotFound404Page.css';

var RouteNotFound404Page = React.createClass({
  displayName: 'RouteNotFound404Page',
  render: function() {
    return (
      <div className={styles.wrap}>
        <div className="row">
          <div className="large-8 large-centered columns">
            <div className={styles.messageModule}>
              <h1 className={styles.heading}>404</h1>
              <h2 className={styles.subheading}> Page Not Found</h2>
              <p className={styles.message}>Sorry, but the page you were trying to view does not exist.</p>
            </div>
          </div>
        </div>
      </div>
    );
  }
});

module.exports = RouteNotFound404Page;
