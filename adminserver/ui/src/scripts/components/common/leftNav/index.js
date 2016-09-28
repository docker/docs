'use strict';

import React, { Component, PropTypes } from 'react';
const {
  bool,
  func
} = PropTypes;
import ui from 'redux-ui';
import classnames from 'classnames';
import FA from 'components/common/fontAwesome';
import styles from './leftNav.css';
import LeftNavItem from './leftNavItem';

@ui() // injects updateUI
export default class LeftNav extends Component {

  static propTypes = {
    isAdmin: bool,
    navExpanded: bool,
    updateUI: func
  }

  toggle = (evt) => {
    evt.preventDefault();
    this.props.updateUI({
      navExpanded: !this.props.navExpanded
    });
  }

  render() {
    const {
      isAdmin,
      navExpanded
    } = this.props;
    const navContainer = classnames({
      [styles.expanded]: navExpanded,
      [styles.container]: true
    });

    return (
      <div className={ navContainer }>
        <button
          onClick={ this.toggle }>
          { navExpanded ?
            <FA icon='fa-arrow-left '/>
            : <FA icon='fa-bars' /> }
        </button>
        <LeftNavItem
          page='/repositories'
          pageName='Repositories'
          icon='fa-book'
          navExpanded={ navExpanded } />
        <LeftNavItem
          page='/orgs'
          pageName='Organizations'
          icon='fa-users'
          navExpanded={ navExpanded } />
        <LeftNavItem
          page='/users'
          pageName='Users'
          icon='fa-user'
          navExpanded={ navExpanded } />
        { isAdmin ?
          <LeftNavItem
            page='/admin/settings/'
            pageName='Settings'
            icon='fa-cog'
            navExpanded={ navExpanded } />
          : undefined }
      </div>
    );
  }

}
