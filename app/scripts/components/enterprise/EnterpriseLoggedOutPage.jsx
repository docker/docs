'use strict';

import React, {
  Component,
  PropTypes
  } from 'react';
import { Link } from 'react-router';
import Card, { Block } from '@dux/element-card';
import styles from './EnterpriseLoggedOutPage.css';
import LoginForm from '../welcome/LoginForm.jsx';
import enterpriseAttemptLogin from '../../actions/enterpriseAttemptLogin';
import plansAndSubscriptions from './plansAndSubscriptions.js';
const { oneOf } = PropTypes;
import has from 'lodash/object/has';
const debug = require('debug')('EnterpriseLoggedOutPage');
const plans = ['server', 'cloud', 'trial', 'micro', 'small', 'medium', 'large', 'xlarge', 'xxlarge'];

export default class EnterpriseLoggedOutPage extends Component {
  static propTypes = {
    type: oneOf(plans).isRequired
  }

  render() {
    let planType = this.props.type;
    if (!has(plansAndSubscriptions, planType)) {
      //handle bad plan link with default
      planType = 'micro';
    }
    const {
      type,
      name,
      includes,
      redirect_value,
      notes
    } = plansAndSubscriptions[planType];

    const contactSales = <a href='http://goto.docker.com/sales-inquiry.html' target="_blank">contact sales</a>;
    const signUp = (
      <Link to='/register/' query={{ redirect_value }}>
        Sign up
      </Link>
    );
    const signUpFree = (
      <Link to='/register/' query={{ redirect_value }}>
        {`Sign up free today.`}
      </Link>
    );
    const includesArea = (
      <ul>
        { includes.map((line) => <li key={line}>{line}</li>) }
      </ul>
    );
    const complete = planType === 'trial' ? 'start your trial' : 'complete the online purchase';
    return (
      <Card>
        <Block>
          <div className='row'>
            <div className='large-7 columns'>
              <div className={styles.leftCol}>
                <div className={styles.leftHeader}>Thank you for your interest in {type}</div>
                <p>The {name}:</p>
                {includesArea}
                <p>{signUp} or log in with your Docker ID to {complete}.</p>
                <p>For questions, {contactSales}</p>
                <p>{notes}</p>
              </div>
            </div>
            <div className='large-5 columns'>
              <div className={styles.rightCol}>
                <div className={styles.rightHeader}>Please enter your Docker ID to continue</div>
                <p>Don’t have a Docker ID? {signUpFree}</p>
                <LoginForm loginAction={enterpriseAttemptLogin} includeHelp/>
              </div>
            </div>
          </div>
        </Block>
      </Card>
    );
  }
}
