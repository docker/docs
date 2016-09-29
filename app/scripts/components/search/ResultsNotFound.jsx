'use strict';

import React, { Component } from 'react';
import styles from './ResultsNotFound.css';

/*
 * A component similar to the RouteNotFound404Page meant for no results
 * from a search or filter. Any of the props can be customized beyond text, ex:
 *  <ResultsNotFound heading={  <FA icon='fa-search' /> } />
 * to include a search icon.
 */
export default class ResultsNotFound extends Component {

  static defaultProps = {
    heading: 'Sorry!',
    subheading: 'We couldn\'t find any results for this search.',
    message: 'Please double check your input and try again.'
  }

  render() {
    return (
      <div className={styles.wrap}>
        <div className="row">
          <div className="large-12 columns">
            <div className={styles.messageModule}>
              <h1 className={styles.heading}>{this.props.heading}</h1>
              <h2 className={styles.subheading}>{this.props.subheading} </h2>
              <p className={styles.message}>{this.props.message}</p>
            </div>
          </div>
        </div>
      </div>
    );
  }
}
