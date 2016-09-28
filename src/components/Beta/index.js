import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import isStaging from 'lib/utils/isStaging';
import isDev from 'lib/utils/isDevelopment';
import css from './styles.css';
import LoginForm from 'common/LoginForm';
import BetaForm from './BetaForm';
import {
  accountFetchCurrentUser,
  accountFetchUserEmails,
} from 'actions/account';
import {
  whitelistAmIWaiting,
  whitelistFetchAuthorization,
} from 'actions/whitelist';
import routes from 'lib/constants/routes';
import { readCookie } from 'lib/utils/cookie-handler';
import DDCBanner from 'components/Home/DDCBanner';
import HelpArticlesCards from 'components/Home/HelpArticlesCards';
const { array, bool, func, object, shape } = PropTypes;

const mapStateToProps = ({ account }) => {
  const { currentUser, isCurrentUserBetalisted, userEmails } = account;
  return {
    currentUser: currentUser || {},
    isCurrentUserBetalisted,
    userEmails,
  };
};

const dispatcher = {
  accountFetchCurrentUser,
  accountFetchUserEmails,
  whitelistAmIWaiting,
  whitelistFetchAuthorization,
};

@connect(mapStateToProps, dispatcher)
export default class Beta extends Component {
  static propTypes = {
    currentUser: object.isRequired,
    isCurrentUserBetalisted: bool.isRequired,
    location: object.isRequired,
    userEmails: shape({
      results: array,
    }),
    accountFetchCurrentUser: func.isRequired,
    accountFetchUserEmails: func.isRequired,
    whitelistAmIWaiting: func.isRequired,
    whitelistFetchAuthorization: func.isRequired,
  }

  static contextTypes = {
    router: shape({
      push: func.isRequired,
    }).isRequired,
  }

  state = { error: '' }

  onError = (error) => {
    this.setState({ error });
  }

  maybeRenderError() {
    const { error } = this.state;
    return error ? <div className={css.error}>{error}</div> : null;
  }

  betaSuccess = (values) => {
    // onSuccess of beta signup submission
    // We need to fetch whether a user has access to the whitelist
    const { currentUser } = this.props;
    this.props.whitelistAmIWaiting();
    analytics.identify(currentUser.id, {
      firstName: values.firstName,
      lastName: values.lastName,
      company: values.company,
    }, () => {
      analytics.track('Signed Up for Private Beta');
    });
  };

  loginSuccess = () => {
    analytics.track('Logged In on Beta Page');

    // If user is on the whitelist then redirect them to the landing page
    this.props.whitelistFetchAuthorization().then(() => {
      this.context.router.replace({ pathname: routes.home() });
    });

    this.props.accountFetchCurrentUser().then(user => {
      const username = user.username;
      this.props.whitelistAmIWaiting();
      this.props.accountFetchUserEmails({ user: username }).then(emails => {
        // If user is already authorized, send them to the landing page

        const primaryEmails = emails &&
          emails.results &&
          emails.results.filter(r => r.primary);
        const primaryEmail = primaryEmails &&
          primaryEmails.length &&
          primaryEmails[0].email;

        if (user && primaryEmail) {
          analytics.identify(user.id, {
            Docker_Hub_User_Name__c: user.username,
            username: user.username,
            dockerUUID: user.id,
            point_of_entry: 'docker_store',
            email: primaryEmail,
          });
        }
      });
    });
  }

  renderLogin() {
    const endpoint = '/v2/users/login/';
    let registerUrl = 'https://cloud.docker.com/';
    let forgotPasswordUrl = 'https://cloud.docker.com/reset-password';
    if (isStaging() || isDev()) {
      registerUrl = 'https://cloud-stage.docker.com/';
      forgotPasswordUrl = 'https://cloud-stage.docker.com/reset-password';
    }
    return (
      <div className={css.login}>
        <div className={css.subText}>
          Login with your Docker ID to access the Beta
        </div>
        {this.maybeRenderError()}
        <div className={css.form} key="login-sign-in">
          <LoginForm
            autoFocus
            csrftoken={readCookie('csrftoken')}
            endpoint={endpoint}
            onSuccess={this.loginSuccess}
            onError={this.onError}
          />
          <div className={css.signupFlow}>
            <a href={forgotPasswordUrl}>Forgot Password?</a> |
            <a href={registerUrl}>Create an account</a>
          </div>
        </div>
      </div>
    );
  }

  renderBetaSignup() {
    const { userEmails } = this.props;
    return (
      <div className={css.beta}>
        <div className={css.subText}>
          Currently Beta is invite-only.
          Fill in the form below to request your invite.
        </div>
        <div className={css.form} key="beta-signup">
          <BetaForm
            onSuccess={this.betaSuccess}
            emails={userEmails.results}
          />
        </div>
      </div>
    );
  }

  renderHero() {
    const {
      currentUser,
      isCurrentUserBetalisted,
      userEmails,
    } = this.props;
    let form;
    const isLoggedIn = currentUser && currentUser.id || false;
    if (isLoggedIn && !!userEmails.results && userEmails.results.length > 0) {
      if (isCurrentUserBetalisted) {
        form = (
          <div className={css.betaThanks}>
            Thank you! We've received your request and we'll email you
            when you are accepted to the Beta.
          </div>
        );
      } else {
        form = this.renderBetaSignup();
      }
    } else {
      form = this.renderLogin();
    }
    return (
      <div className={css.heroWrapper}>
        <div className="wrapped">
          <div className={css.heroContent}>
            <div className={css.title}>Docker Store Beta</div>
            {form}
          </div>
        </div>
      </div>
    );
  }

  render() {
    return (
      <div className={css.home}>
        {this.renderHero()}
        <DDCBanner isBetaPage />
        <HelpArticlesCards />
      </div>
    );
  }
}
