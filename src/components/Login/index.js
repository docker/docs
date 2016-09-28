import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import staticBox from 'lib/decorators/StaticBox';
import { readCookie } from 'lib/utils/cookie-handler';
import isStaging from 'lib/utils/isStaging';
import isDev from 'lib/utils/isDevelopment';
import {
  accountFetchCurrentUser,
  accountFetchUserEmails,
} from 'actions/account';
import { DockerFlatIcon as DockerFlat } from 'common/Icon';
import LoginForm from './LoginForm';
import css from './styles.css';
import qs from 'qs';
import {
  isValidUrl,
  isQsPathSecure,
} from 'lib/utils/url-utils';
import routes from 'lib/constants/routes';

const { func, shape } = PropTypes;
const dispatcher = { accountFetchCurrentUser, accountFetchUserEmails };

@staticBox
@connect(null, dispatcher)
export default class Login extends Component {
  static propTypes = {
    accountFetchCurrentUser: func.isRequired,
    accountFetchUserEmails: func.isRequired,
  }

  static contextTypes = {
    router: shape({
      replace: func.isRequired,
    }).isRequired,
  }

  state = { error: '' }

  onError = (error) => {
    this.setState({ error });
  }

  navigate = () => {
    let query = window && window.location.search;
    // Scrape off first '?' from search
    if (query[0] === '?') {
      query = query.substring(1, query.length);
    }
    const { next } = qs.parse(query);
    const nextUrl = `${window.location.origin}${next}`;

    if (isValidUrl(nextUrl) && isQsPathSecure(next)) {
      window.location = nextUrl;
    } else {
      window.location = routes.home();
    }
  }

  maybeRenderError() {
    const { error } = this.state;
    return error ? <div className={css.error}>{error}</div> : null;
  }

  render() {
    const endpoint = '/v2/users/login/';
    let registerUrl = 'https://cloud.docker.com/';
    let forgotPasswordUrl = 'https://cloud.docker.com/reset-password';
    if (isStaging() || isDev()) {
      registerUrl = 'https://cloud-stage.docker.com/';
      forgotPasswordUrl = 'https://cloud-stage.docker.com/reset-password';
    }
    return (
      <div>
        {this.maybeRenderError()}
        <div className={css.banner}>
          <Link to="/">
            <DockerFlat />
          </Link>
          <h1>Welcome to the Docker Store</h1>
          <p>Login with your <strong>Docker ID</strong></p>
        </div>
        <div className={css.form}>
          <LoginForm
            autoFocus
            csrftoken={readCookie('csrftoken')}
            endpoint={endpoint}
            onSuccess={this.navigate}
            onError={this.onError}
          />
        </div>
        <div className={css.more}>
          <a href={forgotPasswordUrl}>Forgot Password?</a>
          <a href={registerUrl}>Create Account</a>
        </div>
      </div>
    );
  }
}
