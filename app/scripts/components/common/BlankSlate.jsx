'use strict';

import styles from './BlankSlate.css';
import React, { PropTypes, Component } from 'react';
import { Link } from 'react-router';
import FA from './FontAwesome';
import classnames from 'classnames';

const { string, object } = PropTypes;

/*
 * Blankslate
 *
 * Usage:
 * <BlankSlates title='Welcome to Docker Hub' sub='Here are a few things to get your started.'>
 *   <BlankSlate link='addRepo' icon='fa-book' title='Create Repository' query={{namespace: ns}} />
 *   <BlankSlate link='addOrg' icon='fa-users' title='Create Organization' />
 *   <BlankSlate link='explore' icon='fa-compass' title='Explore Repositories' />
 * </BlankSlates>
 */

export class BlankSlates extends Component {

  static propTypes = {
    title: string
  }

  render() {
    return (
      <div className={styles.blankSlates}>
        <div className="row">
          <div className="large-5 columns">
            <h1>{this.props.title}</h1>
            <p>{this.props.subtext}</p>
          </div>
          <div className="large-7 columns">
            {this.props.children}
          </div>
        </div>
      </div>
    );
  }
}

export class BlankSlate extends Component {
  static defaultProps = {
    query: {}
  }

  render() {
    const linkClasses = classnames({
      'button': true,
      [styles.link]: true
    });

    const { link, query, icon, title } = this.props;
    return (
      <Link to={link}
            query={query}
            className={linkClasses}>
        <FA icon={icon}/><br/>{title}
      </Link>
    );
  }
}

BlankSlate.propTypes = {
  link: string.isRequired,
  icon: string.isRequired,
  title: string,
  subtext: string,
  query: object
};
