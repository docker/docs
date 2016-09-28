import React, { Component, PropTypes } from 'react';
import { reduxForm } from 'redux-form';
import {
  Button,
  Input,
} from 'common';
import { whitelistSubscribeToBeta } from 'actions/whitelist';
import EmailSelect from './emailSelect';
import validate from './validations';
import css from './styles.css';

const { array, func, object } = PropTypes;

const dispatcher = {
  whitelistSubscribeToBeta,
};

const fields = [
  'firstName',
  'lastName',
  'company',
  'email',
];

const mapStateToProps = ({ account }, props) => {
  const initialValues = {
    email: props && props.emails && props.emails[0].email,
  };
  return {
    user: account && account.currentUser,
    initialValues,
  };
};

class BetaForm extends Component {

  static propTypes = {
    user: object,
    emails: array.isRequired,
    fields: object.isRequired,
    onSuccess: func.isRequired,
    handleSubmit: func.isRequired,
    whitelistSubscribeToBeta: func.isRequired,
  }

  state = {
    inProgress: false,
  }

  onSelectChange = (field) => (data) => {
    const fieldObject = this.props.fields[field];
    fieldObject.onChange(data.value);
  }

  onSubmit = (values) => {
    this.setState({ inProgress: true });
    // submit form data to api backend
    const {
      whitelistSubscribeToBeta: subscribeToBeta,
      onSuccess,
    } = this.props;
    const {
      firstName: first_name,
      lastName: last_name,
      company,
      email,
    } = values;
    subscribeToBeta({
      first_name,
      last_name,
      company,
      email,
    }).then(() => {
      onSuccess(values);
      this.setState({ inProgress: false });
    });
  }

  render() {
    const {
      emails,
      fields: propFields,
      handleSubmit,
    } = this.props;
    const {
      inProgress,
    } = this.state;
    const {
      firstName,
      lastName,
      company,
    } = propFields;
    const submitText = inProgress ? 'Submitting...' : 'Request Beta Access';
    const emailsArray = emails && emails.map(email => email.email);
    return (
      <form
        key="beta-form"
        className={css.main}
        onSubmit={handleSubmit(this.onSubmit)}
      >
        <Input
          {...firstName}
          className={css.input}
          id={'first'}
          placeholder="First Name"
          inputStyle={{ color: 'white', width: '100%' }}
          underlineFocusStyle={{ borderColor: 'white' }}
          errorText={firstName.touched && firstName.error}
        />
        <Input
          {...lastName}
          className={css.input}
          id={'last'}
          placeholder="Last Name"
          inputStyle={{ color: 'white', width: '100%' }}
          underlineFocusStyle={{ borderColor: 'white' }}
          errorText={lastName.touched && lastName.error}
        />
        <Input
          {...company}
          className={css.input}
          id={'company'}
          placeholder="Company"
          inputStyle={{ color: 'white', width: '100%' }}
          underlineFocusStyle={{ borderColor: 'white' }}
          errorText={company.touched && company.error}
        />
        <EmailSelect
          accountEmails={emailsArray}
          fields={propFields}
          onSelectChange={this.onSelectChange('email')}
        />
        <Button
          disabled={inProgress}
          className={css.submit}
          inverted type="submit"
        >
          {submitText}
        </Button>
      </form>
    );
  }
}

export default reduxForm({
  form: 'betaForm',
  fields,
  validate,
},
mapStateToProps,
dispatcher,
)(BetaForm);
