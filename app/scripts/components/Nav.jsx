'use strict';

import styles from './Nav.css';

import React, { createClass, PropTypes, Component } from 'react';
import _ from 'lodash';
import classnames from 'classnames';
import { Link } from 'react-router';

import CreateButton from './CreateButton.jsx';
import FA from 'common/FontAwesome';
import Logout from 'actions/logout';
import SearchBar from './search/SearchBar.jsx';
import { TopNav, Button } from 'dux';
import { mkAvatarForNamespace } from 'utils/avatar';

var debug = require('debug')('Navbar');

class LoggedInLeftNav extends Component {
  render() {
    return (
      <ul className="leftNav">
        <li><Link to="/">Dashboard</Link></li>
        <li><Link to="/explore/">Explore</Link></li>
        <li><Link to="/organizations/">Organizations</Link></li>
      </ul>
    );
  }
}

class LoggedOutLeftNav extends Component {
  render() {
    return (
      <ul className="rightNav">
        <li><Link to='/explore/'>Explore</Link></li>
        <li><Link to='/help/'>Help</Link></li>
      </ul>
    );
  }
}

class LiDropdown extends Component {
  state = {
    isOpen: false
  }
  _toggleDropdown = (e) => {
    debug('clicked', this.state);
    this.setState({
      isOpen: !this.state.isOpen
    });
  }
  _closeDropdown = (e) => {
    this.setState({
      isOpen: false
    });
  }
  render() {
    debug('LiDropdown', this.state);
    let dropdownClasses = classnames({
      'has-dropdown': true,
      'hover': this.state.isOpen
    });
    return (
      <li className={dropdownClasses} onMouseLeave={_.debounce(this._closeDropdown, 200)}>
        <a onClick={this._toggleDropdown}>{this.props.title}</a>
        <ul className='dropdown'>{this.props.children}</ul>
      </li>
    );
  }
}

var LoggedInRightNav = React.createClass({
  displayName: 'LoggedInRightNav',
  propTypes: {
    user: PropTypes.shape({
      username: PropTypes.string
    }),
    JWT: PropTypes.string,
    history: PropTypes.object
  },
  contextTypes: {
    executeAction: PropTypes.func.isRequired
  },
  getDefaultProps() {
    return {
      user: {
        username: ''
      }
    };
  },
  _logout(e) {
    debug('signing out!');
    e.preventDefault();
    this.context.executeAction(Logout, this.props.JWT);
  },
  render() {
    var username = this.props.user.username;

    const { params } = this.props;
    //Org/User namespace for creating autobuilds
    let currentNamespace = username;
    if (params && params.user) {
      currentNamespace = params.user;
    }

    return (
      <ul className="right">
        <li><SearchBar history={this.props.history}/></li>
        <LiDropdown title='Create'>
          <li><Link to='/add/repository/' query={{namespace: currentNamespace}}><FA icon='fa-book'/> Create Repository</Link></li>
          <li><Link to={`/add/automated-build/${currentNamespace}/`}><FA icon='fa-cogs' /> Create Automated Build</Link></li>
          <li><Link to='/organizations/add/'><FA icon='fa-users'/> Create Organization</Link></li>
        </LiDropdown>
        <LiDropdown title={<span><img className={styles.avatar}
                           src={mkAvatarForNamespace(username)}/><span className={styles.username}>{username}</span></span>}>
          <li><Link to={`/u/${username}/`}>My Profile</Link></li>
          <li className="divider"></li>
          <li><a href='https://docs.docker.com/docker-hub/'>Documentation</a></li>
          <li><Link to='/help/'>Help</Link></li>
          <li className="divider"></li>
          <li><Link to='/account/settings/'>Settings</Link></li>
          <li><a onClick={this._logout}>Log out</a></li>
        </LiDropdown>
      </ul>
    );
  }
});

class LoggedOutRightNav extends Component {
  static propTypes = {
    isHomePage: PropTypes.bool,
    history: PropTypes.object
  };
  onClick = (e) => {
    e.preventDefault();
    this.props.history.pushState(null, '/');
  };
  render() {
    const { isHomePage, history } = this.props;
    let signupButton = (
      <li><Button size="tiny" onClick={this.onClick}>Sign up</Button></li>
    );
    if (isHomePage) {
      signupButton = null;
    }
    return (
      <ul className="right">
        <li><SearchBar history={history}/></li>
        {signupButton}
        <li><Link className="tiny" to='/login/'>Log In</Link></li>
      </ul>
    );
  }
}

//This right nav element is specific to the Registration form in the standalone register page
class LoggedOutRegisterRightNav extends Component {
  render() {
    const { location } = this.props;
    //Default
    var loginLink = <Link to='/login/'>Log In</Link>;

    //Other refs handled here
    var partner_value = location.query.partner_value;
    if (partner_value === 'tutum') {
      debug(process.env.TUTUM_SIGNIN_URL);
      loginLink = <a href={process.env.TUTUM_SIGNIN_URL}>Login to Tutum</a>;
    }

    return (
      <ul className="right">
        <li>{loginLink}</li>
      </ul>
    );
  }
}

export default createClass({
  displayName: 'Nav',
  propTypes: {
    isLoggedIn: PropTypes.bool.isRequired,
    isHomePage: PropTypes.bool,
    user: PropTypes.shape({
      username: PropTypes.string
    }),
    JWT: PropTypes.string,
    history: PropTypes.object,
    location: PropTypes.object,
    params: PropTypes.object
  },
  getDefaultProps: function() {
    return {
      isLoggedIn: false,
      user: {}
    };
  },
  render() {
    const { history, location, JWT, user, params } = this.props;
    var left = (<LoggedOutLeftNav />);
    var right = (
      <LoggedOutRightNav
        isHomePage={this.props.isHomePage}
        history={history}/>
    );
    //If register route is active, show the register specific nav
    if (history.isActive('/register/')) {
      right = (<LoggedOutRegisterRightNav location={location}/>);
    }
    if (this.props.isLoggedIn) {
      left = (<LoggedInLeftNav />);
      right = (
        <LoggedInRightNav JWT={JWT}
                          history={history}
                          params={params}
                          user={user} />
      );
    }
    return (
        <div className="contain-to-grid">
          <TopNav>
            <ul className="title-area">
              <li>
                <Link to='/'>
                  <img src="/public/images/logos/mini-logo.svg" alt='docker logo' className={styles.logo}/>
                </Link>
              </li>
            </ul>
            {left}
            {right}
          </TopNav>
        </div>
    );
  }
});
