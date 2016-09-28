import React, { Component, PropTypes } from 'react';
import { reduxForm } from 'redux-form';
import { Button, Card, Checkbox, CheckIcon, Input } from 'common';
import { TOS_LINK, PRIVACY_POLICY_LINK } from 'lib/constants/eusa';
import { XLARGE } from 'lib/constants/sizes';
import { login, signup } from 'actions/account';
import css from './styles.css';
import get from 'lodash/get';
import imageVerify from 'lib/images/verify.png';

const { object, func, bool, string } = PropTypes;

class SignupLoginForm extends Component {
  static propTypes = {
    location: object.isRequired,
    fields: object.isRequired,
    handleSubmit: func.isRequired,
    submitting: bool.isRequired,
    untouch: func.isRequired,
    plan: string.isRequired,
    error: string,
    resetForm: func.isRequired,
  }

  static contextTypes = {
    router: object.isRequired,
  }

  state = {
    verify: false,
  }

  onCheck = () => {
    const { accepted } = this.props.fields;
    const checked = accepted.value === 'checked' ? 'unchecked' : 'checked';
    accepted.onChange(checked);
  }

  onChangeAccount() {
    this.props.resetForm();
  }

  submitForm(values) {
    const { plan } = this.props;
    if (values.account === 'existing') {
      const { username, password } = values;
      return login(username, password).then(() => {
        analytics.track('ddc_login', { plan });
        location.reload();
      });
    }

    const { email, username, password } = values;
    return signup(email, username, password, plan)
      .then(() => {
        analytics.track('ddc_create_account', { plan });
        this.setState({ verify: true });
      });
  }

  renderLogin() {
    const {
      fields: {
        username,
        password,
      },
    } = this.props;
    return (
      <div>
        <Input
          {...username}
          id={'username'}
          placeholder="Username or email"
          errorText={username.touched && username.error}
          style={{ width: '100%', marginBottom: '10px' }}
        />
        <Input
          {...password}
          id={'password'}
          placeholder="Password"
          errorText={password.touched && password.error}
          type="password"
          style={{ width: '100%', marginBottom: '10px' }}
        />
      </div>
    );
  }

  renderSignup() {
    const {
      fields: {
        username,
        email,
        password,
      },
    } = this.props;

    return (
      <div>
        <Input
          {...username}
          id={'username'}
          placeholder="Choose a username"
          errorText={username.touched && username.error}
          style={{ width: '100%', marginBottom: '10px' }}
        />
        <Input
          {...email}
          id={'email'}
          placeholder="Email"
          type="email"
          errorText={email.touched && email.error}
          style={{ width: '100%', marginBottom: '10px' }}
        />
        <Input
          {...password}
          id={'password'}
          placeholder="Password"
          errorText={password.touched && password.error}
          type="password"
          style={{ width: '100%', marginBottom: '10px' }}
        />
      </div>
    );
  }

  renderForm() {
    const {
      fields: {
        accepted,
        account,
      },
      handleSubmit,
      submitting,
      error,
    } = this.props;

    const existingAccount = account.value === 'existing';
    const dockerIDForm = existingAccount ?
      this.renderLogin() : this.renderSignup();
    const tos = existingAccount ? null : (
      <div className={css.fields}>
        <div className={css.euta}>
          <Checkbox
            {...accepted}
            checked={accepted.value === 'checked'}
            onBlur={() => {}}
            onCheck={this.onCheck}
            style={{ width: '' }}
          />
          <div>
            I have read and agree to Docker&#39;s&nbsp;
            <a tabIndex="-1" target="_blank" href={TOS_LINK}>
              Terms of Service
            </a>
            &nbsp;and&nbsp;
            <a tabIndex="-1" target="_blank" href={PRIVACY_POLICY_LINK}>
              Privacy Policy
            </a>.
          </div>
        </div>
      </div>
    );

    const verified = this.props.location.query.verified;
    const accountChoice = verified ? null : (
      <div className={css.already}>
        <label>
          <input
            {...account}
            type="radio"
            value="new"
            checked={account.value === 'new'}
            tabIndex="-1"
            onClick={() => this.onChangeAccount()}
          />&nbsp;Create a new Docker ID
        </label>
        <label>
          <input
            {...account}
            type="radio"
            value="existing"
            checked={account.value === 'existing'}
            tabIndex="-1"
            onClick={() => this.onChangeAccount()}
          />&nbsp;Use my existing Docker ID
        </label>
      </div>
    );

    return (
      <div className={css.contact}>
        <form
          onSubmit={handleSubmit(::this.submitForm)}
          className={css.form}
        >
          <div>
            <div className={css.fields}>
              {accountChoice}
              {dockerIDForm}
            </div>
          </div>
          {tos}
          <div className={css.error}>
            {(accepted.touched && accepted.error)}
            {error}
          </div>
          <div className={css.continue}>
            <Button
              id="submit"
              className={css.submit}
              type="submit"
              disabled={submitting}
              onClick={handleSubmit(::this.submitForm)}
            >
              Next
            </Button>
          </div>
        </form>
      </div>
    );
  }

  renderVerify() {
    return (
      <div className={css.verify}>
        <div className={css.verifyTitle}>
          Verify your email address
        </div>
        <img src={imageVerify} alt="envelope icon" />
        <div className={css.verifySubtitle}>
          We&#39;ve sent you a verification email.<br />
          Click the included link to start your evaluation.
        </div>
      </div>
    );
  }

  render() {
    const form = this.renderForm();
    const verify = this.renderVerify();
    const verified = this.props.location.query.verified;
    const header = verified ? (
      <div className={css.title}>
        <CheckIcon size={XLARGE} className={css.check} />
        <div className={css.titleText}>
          You&#39;ve verified your Docker ID<br />
          <div className={css.titleSubText}>
            Sign in to begin your Docker Datacenter subscription
          </div>
        </div>
      </div>
    ) : null;
    const titleText = this.props.plan === 'free-trial' ?
      'Sign up for your 30 day free evaluation of Docker Datacenter' :
      'Sign up to Purchase a Docker Datacenter Subscription';
    const title = verified ? '' : titleText;
    return (
      <Card
        className={css.card}
        title={title}
      >
        {header}
        {this.state.verify ? verify : form}
      </Card>
    );
  }
}

const mapStateToProps = ({}, { location }) => {
  const verified = get(location, ['query', 'verified']);
  return {
    initialValues: {
      account: verified ? 'existing' : 'new',
      accepted: 'unchecked',
    },
  };
};

const dispatcher = {};

export default reduxForm({
  form: 'signupForm',
  fields: [
    'accepted',
    'account',
    'email',
    'username',
    'password',
  ],
  validate: (values) => {
    const errors = {};
    const {
      accepted,
      account,
      email,
      username,
      password,
    } = values;
    if (account === 'new' && accepted !== 'checked') {
      errors.accepted = 'Please accept the terms of use before continuing';
    }
    if (account === 'new' && !email) {
      errors.email = 'Required';
    }
    if (!username) {
      errors.username = 'Required';
    }
    if (!password) {
      errors.password = 'Required';
    }
    return errors;
  },
},
mapStateToProps,
dispatcher,
)(SignupLoginForm);
