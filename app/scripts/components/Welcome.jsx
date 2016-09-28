'use strict';

import styles from './Welcome.css';

import React, { createClass, PropTypes } from 'react';
import LoginForm from './welcome/LoginForm.jsx';
import SignupForm from './welcome/SignupForm.jsx';
import { Link } from 'react-router';
import classnames from 'classnames';
import Button from '@dux/element-button';
var debug = require('debug')('Welcome');

let LowerSection = createClass({
  displayName: 'LowerSection',
  propTypes: {
    companyList: PropTypes.arrayOf(PropTypes.object)
  },
  getDefaultProps() {
    return {
      companyList: [
        {logo: 'harvard'},
        {logo: 'sony'},
        {logo: 'nordstrom'},
        {logo: 'oracle'},
        {logo: 'zendesk'},
        {logo: 'gopro'},
        {logo: 'autodesk'},
        {logo: 'cisco'},
        {logo: 'zenefits'},
        {logo: 'dollar-shave-club'},
        {logo: 'zipcar'},
        {logo: 'gilt'},
        {logo: 'adp'},
        {logo: 'makerbot'},
        {logo: 'oculus'},
        {logo: 'adobe'}
      ]
    };
  },
  mkCompanyLi({name, url, logo}) {
    let img = `/public/images/customers/${logo}.png`;
    return (
      <li key={logo} className={styles.companyLinkLi}>
        <img src={img} />
      </li>
    );
  },
  render() {

    let companyClasses = classnames({
      'small-block-grid-4': true,
      [styles.companyRow]: true
    });

    return (
      <section className="callout-section">
        <h3 className={styles.lowerHeading}>The most innovative companies use Docker</h3>
        <div className="row">
          <div className="large-12 columns">
            <ul className={companyClasses}>
              {this.props.companyList.map(str => { return this.mkCompanyLi(str); })}
            </ul>
          </div>
        </div>
      </section>
    );
  }
});

var Welcome = React.createClass({
  displayName: 'WelcomePage',
  transitionExplore: function(e){
    e.preventDefault();
    this.props.history.pushState(null, '/explore/');
  },
  render() {

    let buildShipRun = classnames({
      'large-8 columns': true,
      [styles.buildShipRun]: true
    });

    return (
      <div className={styles.flex}>
        <header className={styles.header}>
          <div className={'row ' + styles.top}>

            <div className={buildShipRun}>
              <h1 className={styles.headingHero}>Build, Ship, & Run</h1>
              <h1 className={styles.headingHeroAlt}>Any App, Anywhere</h1>
              <p className={styles.subtext}>Dev-test pipeline automation, 100,000+ free apps, public and private registries</p>
            </div>

            <div className='large-4 columns'>
              <SignupForm location={this.props.location}/>
            </div>

          </div>
        </header>
        <div className={'row ' + styles.browse}>
          <Link to='/explore/'>Browse</Link> Thousands of the most popular software tools in the Docker Image Library
        </div>
        <div className={styles.white}>
          <div className={styles.footer}>
            <p className={styles.footerCopy}>&copy; 2016 Docker Inc.</p>
          </div>
        </div>
      </div>
    );
  }
});

module.exports = Welcome;
