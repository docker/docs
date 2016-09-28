'use strict';

import React from 'react';
import { PageHeader, Module } from 'dux';
import FA from 'common/FontAwesome';
import styles from './Help.css';
import { HelpContents } from './HelpContents';

var debug = require('debug')('Help');

// render single row within helpItem
const renderRow = ({href, text}, idx) => {
  if (href === undefined) {
    return <div className={styles.helpItemText} key={idx}><span>{text}</span></div>;
  } else {
    return <div className={styles.helpItemLink} key={idx}><a href={href} target='_blank'>{text}</a></div>;
  }
};

// render subtext btn if required
const renderSubtext = (arg = {}) => {
  const {text, btnStyle, link} = arg;

  // open a new tab given a link/href
  const goToLink = (l) => () => window.open(l, '_blank');

  if (text) {
    return (
      <div className={styles.helpItemSubtext}>
        <button className={btnStyle} onClick={goToLink(link)}>
          {text}
        </button>
      </div>
    );
  }
};

// render helper for each helpItem column
const renderColumn = ({ title, icon, body, subtext }, index) => {
  const iconClass = icon + ' ' + styles.helpItemIcon;
  return (
    <div key={index} className={styles.helpItem}>
      <div className={styles.helpItemContent}>
        <div className={styles.helpItemIconWrapper}>
          <FA icon={iconClass}/>
        </div>
        <h5 className={styles.helpItemTitle}>{title}</h5>
        <div className={styles.helpItemBody}>
          {body.map(renderRow)}
        </div>
      </div>
      {renderSubtext(subtext)}
    </div>
  );
};

const Home = React.createClass({
  displayName: 'Help',
  render() {
    return (
      <div>
        <PageHeader title='Help' />
        <div className='row'>
          <div className='large-12 large-centered columns'>
            <div className={styles.helpDescription}>
              <h5>Docker Hub Resources</h5>
              <span>Learn more about Docker, find & share solutions, and access support</span>
            </div>
            <Module>
              <div className={styles.helpItemList}>
                {HelpContents.map(renderColumn)}
              </div>
            </Module>
          </div>
        </div>
      </div>
    );
  }
});

module.exports = Home;

