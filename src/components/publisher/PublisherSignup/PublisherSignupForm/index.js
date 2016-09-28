import React, { Component, PropTypes } from 'react';
import { reduxForm } from 'redux-form';
import { Input, Button } from 'components/common';
import css from './styles.css';

class PublisherSignupForm extends Component {
  static propTypes = {
    fields: PropTypes.object.isRequired,
    submitting: PropTypes.bool.isRequired,
    handleSubmit: PropTypes.func.isRequired,
  }

  render() {
    const {
      fields: {
        first_name,
        last_name,
        company,
        phone_number,
      },
      submitting,
      handleSubmit,
    } = this.props;

    return (
      <form className={css.form} onSubmit={handleSubmit}>
        <div>
          <Input
            id={'first_name_input'}
            placeholder="First Name"
            style={{ width: '100%' }}
            errorText={first_name.touched && first_name.error}
            { ...first_name}
          />
        </div>
        <div>
          <Input
            id={'last_name_input'}
            placeholder="Last Name"
            style={{ width: '100%' }}
            errorText={last_name.touched && last_name.error}
            { ...last_name}
          />
        </div>
        <div>
          <Input
            id={'company_input'}
            placeholder="Company"
            style={{ width: '100%' }}
            errorText={company.touched && company.error}
            { ...company}
          />
        </div>
        <div>
          <Input
            id={'phone_number_input'}
            placeholder="Phone"
            style={{ width: '100%' }}
            errorText={phone_number.touched && phone_number.error}
            { ...phone_number}
          />
        </div>
        <div className={css.center}>
          <Button
            className={css.button}
            style={{ width: '100%' }}
            disabled={submitting}
          >
            Sign Up to Become a Publisher
          </Button>
        </div>
      </form>
    );
  }
}

export default reduxForm({
  form: 'publisherEnrollForm',
  fields: [
    'first_name',
    'last_name',
    'company',
    'phone_number',
  ],
  validate: values => {
    const errors = {};

    if (!values.first_name) {
      errors.first_name = 'Required';
    }

    if (!values.last_name) {
      errors.last_name = 'Required';
    }

    if (!values.company) {
      errors.company = 'Required';
    }

    return errors;
  },
})(PublisherSignupForm);
