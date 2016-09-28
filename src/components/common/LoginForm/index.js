import React, { Component, PropTypes } from 'react';
import { post } from 'superagent';
import Input from 'common/Input';
import Button from 'common/Button';
import css from './styles.css';
import isValidPassword from 'lib/utils/password-validator';

const { string, func, bool } = PropTypes;

// Not sure where this should go yet
const isValidDockerID = (val = '') => val.length;

// NOTE: This component has to be refactored, as it was developed
// in a rush for Docker Store's launch
export default class LoginForm extends Component {
  static propTypes = {
    endpoint: string,
    onSuccess: func,
    onError: func,
    autoFocus: bool,
    csrftoken: string,
  }

  static defaultProps = {
    onSuccess() {},
    onError() {},
  }

  static getInput(value, label) {
    return (
      <Input
        className={css.input}
        value={value}
        id={label}
        inputStyle={{ color: 'white' }}
        underlineFocusStyle={{ borderColor: 'white' }}
        hintText={label}
      />
    );
  }

  state = {
    id: '',
    password: '',
    inProgress: false,
    errors: {
      id: '',
      password: '',

      // This field is disregarded in the form.
      // If we encounter it through the ajax call,
      // it's passed as-is to the onError callback()
      // detail: '',
    },
  }

  getHintStyle(value) {
    return !value ? { color: 'white', opacity: 0.6 } : null;
  }

  handleChange(which, ev) {
    this.setState({ [which]: ev.target.value });
  }

  isValidState() {
    const { id, password } = this.state;
    return isValidDockerID(id) && isValidPassword(password);
  }

  validateDefault(which, { allowEmpty = false } = {}) {
    const { id, password } = this.state;
    let errors = Object.assign({}, this.state.errors);

    const errorId = () => {
      if ((allowEmpty && !id.length) || isValidDockerID(id)) return null;
      return ' ';
    };

    const errorPassword = () => {
      if ((allowEmpty && !password.length) || isValidPassword(password)) {
        return null;
      }
      return ' ';
    };

    switch (which) {
      case 'id': errors.id = errorId(); break;
      case 'password': errors.password = errorPassword(); break;
      default:
        errors = {
          id: errorId(),
          password: errorPassword(),
        };
    }

    this.setState({ errors });
  }

  handleSubmit = (ev) => {
    ev.preventDefault();

    if (!this.isValidState()) {
      this.validateDefault();
      return;
    }

    const { password } = this.state;
    const id = this.state.id && this.state.id.toLowerCase();
    const { endpoint, onSuccess, onError, csrftoken } = this.props;

    this.setState({
      inProgress: true,
      errors: {
        id: null,
        password: null,
      },
    });

    const req = post(endpoint)
      .set('Content-Type', 'application/json')
      .set('Accept', 'application/json');

    if (csrftoken) {
      req.set('X-CSRFToken', csrftoken);
    }

    req.send({ password, username: id }).end((err, res) => {
      this.setState({ inProgress: false });

      if (err) {
        const errors = {};
        let body = {};

        try {
          body = JSON.parse(res.text);
        } catch (e) {
          onError(res.text);
          return;
        }

        if (err.status === 401) {
          errors.id = ' ';
          errors.password = ' ';
        }

        if (body.detail) {
          onError(body.detail);
        }

        if (body.username) {
          errors.id = body.username[0];
        }

        if (body.password) {
          errors.password = body.password[0];
        }

        if (Object.keys(errors).length) {
          this.setState({ errors });
        }

        return;
      }

      onSuccess({ dockerId: id });
    });
  }

  resetFields() {
    this.setState({ id: '', password: '' });
  }

  renderInputID(value) {
    const { autoFocus } = this.props;
    const baseInput = LoginForm.getInput(value, 'Docker ID');
    const errorText = this.state.errors.id;
    const hintStyle = this.getHintStyle(value);
    const onChange = (ev) => this.handleChange('id', ev);
    const onBlur = () => this.validateDefault('id', { allowEmpty: true });
    return React.cloneElement(
      baseInput,
      { onChange, onBlur, hintStyle, errorText, autoFocus }
    );
  }

  renderInputPassword(value) {
    const baseInput = LoginForm.getInput(value, 'Password');
    const errorText = this.state.errors.password;
    const hintStyle = this.getHintStyle(value);
    const onChange = (ev) => this.handleChange('password', ev);
    const onBlur = () => this.validateDefault('password', { allowEmpty: true });
    return React.cloneElement(
      baseInput,
      { onChange, onBlur, hintStyle, errorText, type: 'password' }
    );
  }

  render() {
    const { id, password, inProgress } = this.state;
    const inputID = this.renderInputID(id);
    const inputPassword = this.renderInputPassword(password);

    return (
      <div className={css.main}>
        <form onSubmit={this.handleSubmit}>
          {inputID}
          {inputPassword}
          <Button
            disabled={inProgress}
            className={css.login}
            inverted type="submit"
          >
            Log in
          </Button>
        </form>
      </div>
    );
  }
}
