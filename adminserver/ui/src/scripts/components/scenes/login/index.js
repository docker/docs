'use strict';

import React from 'react';
import { connect } from 'react-redux';
import cn from 'classnames';

import { mapActions } from 'utils';
import Logo from 'components/common/logo';
import * as AuthActions from 'actions/auth';
import Button from 'components/common/button';
import consts from 'consts';
import styles from './styles.css';

const mapState = (state) => ({ auth: state.auth });

@connect(mapState, mapActions(AuthActions))
export default class Login extends React.Component {
  static propTypes = {
    actions: React.PropTypes.object,
    auth: React.PropTypes.object
  }

  componentWillMount() {
    // When showing the auth form trigger something to show that we're not
    // signed in.
    this.props.actions.showingAuthForm();
  }

  submit = (evt) => {
    evt.preventDefault();
    this.props.actions.logIn({
      username: this.refs.username.value,
      password: this.refs.password.value
    });
  }

  renderMessage() {
    if (this.props.auth.status === consts.loading.FAILURE) {
      return (
        <div>
          <span className={ styles.errorSpan }>Invalid username and/or password. Please try again.</span>
        </div>
      );
    }
  }

  render() {
    let formClass = cn({
      [styles.loginForm]: true
    });

    return (
      <div className={ styles.container }>
        <form className={ formClass } onSubmit={ (evt) => { this.submit(evt); } }>
          <Logo scale={ 5 } className={ styles.logo } pathClassName={ styles.logoFill } />
          <h1 className={ styles.heading }>Docker Trusted Registry</h1>

          { this.renderMessage() }

          <input type='text' ref='username' name='username' placeholder='Name' />
          <input type='password' ref='password' name='password' placeholder='Password' />

          <div className={ styles.submitContainer }>
            <Button type='submit'>Log in</Button>
          </div>
        </form>
      </div>
    );
  }

}
