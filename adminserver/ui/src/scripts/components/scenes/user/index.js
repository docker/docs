'use strict';

import React, { Component, PropTypes } from 'react';
import FA from 'components/common/fontAwesome';
import { Tabs, Tab } from 'components/common/tabs';
import { Link } from 'react-router';
import css from 'react-css-modules';
import styles from './styles.css';
import { createStructuredSelector } from 'reselect';
import { selectUser } from 'selectors/users';
import { connect } from 'react-redux';
import autoaction from 'autoaction';
import { getUser } from 'actions/users';
import Spinner from 'components/common/spinner';
import consts from 'consts';

const mapState = createStructuredSelector({
  user: selectUser
});

@autoaction({
  getUser: (props) => props.params.username
}, { getUser })
@connect(mapState)
@css(styles)
export default class User extends Component {

  static propTypes = {
    params: PropTypes.object,
    user: PropTypes.object,
    children: PropTypes.node
  }

  render() {

    const {
      params,
      user,
      children
    } = this.props;

    const status = [
      [consts.users.GET_USER]
    ];

    return (
      <Spinner
        loadingStatus={ status }>
        <div styleName='wrapper'>
          <div styleName='left'>
            <div styleName='usersIcon'>
              <FA icon='fa-user' />
            </div>
            <h2 styleName='orgName' id='user-detail-username'>{ user.name }</h2>
          </div>
          <div styleName='content'>
            <UserTabsHeader userName={ params.username } />
            { children }
          </div>
        </div>
      </Spinner>
    );
  }
}

class UserTabsHeader extends Component {
  static propTypes = {
    userName: PropTypes.string
  }

  render() {
    const { userName } = this.props;

    return (
    <Tabs header>
      <Tab>
        <Link to={ `/users/${ userName }/repos` }>Repositories</Link>
      </Tab>
      <Tab>
        <Link to={ `/users/${ userName }/teams` }>Teams</Link>
      </Tab>
      <Tab id='user-settings-tab'>
        <Link to={ `/users/${ userName }/settings` }>Settings</Link>
      </Tab>
    </Tabs>
    );
  }
}
