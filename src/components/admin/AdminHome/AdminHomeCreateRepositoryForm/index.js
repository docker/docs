import React, { Component, PropTypes } from 'react';
import { reduxForm } from 'redux-form';
import { Input, Button } from 'components/common';
import css from './styles.css';

class AdminHomeCreateRepositoryForm extends Component {
  static propTypes = {
    fields: PropTypes.object.isRequired,
    submitting: PropTypes.bool.isRequired,
    handleSubmit: PropTypes.func.isRequired,
  }

  render() {
    const {
      fields: {
        display_name,
        namespace,
        reponame,
        publisher,
      },
      submitting,
      handleSubmit,
    } = this.props;

    const repoInputStyle = { width: '200px' };
    const inputStyle = { width: '256px' };

    return (
      <form onSubmit={handleSubmit}>
        <div>
          <Input
            id={'display_name_input'}
            placeholder="Display Name"
            style={inputStyle}
            errorText={display_name.touched && display_name.error}
            {...display_name}
          />
        </div>
        <div>
          store&nbsp;/&nbsp;<Input
            id={'namespace_input'}
            errorText={namespace.touched && namespace.error}
            placeholder="namespace"
            style={repoInputStyle}
            {...namespace}
          />&nbsp;/&nbsp;<Input
            id={'reponame_input'}
            errorText={reponame.touched && reponame.error}
            placeholder="reponame"
            style={repoInputStyle}
            {...reponame}
          />
        </div>
        <div>
          <Input
            id={'publisher_id_input'}
            placeholder="Publisher Docker ID"
            style={inputStyle}
            errorText={publisher.id.touched && publisher.id.error}
            {...publisher.id}
          />
        </div>
        <div>
          <Input
            id={'publisher_name_input'}
            placeholder="Publisher Name"
            style={inputStyle}
            errorText={
              publisher.name.touched && publisher.name.error
            }
            {...publisher.name}
          />
        </div>
        <div className={css.buttons}>
          <Button type="submit" disabled={submitting}>Add Repository</Button>
        </div>
      </form>
    );
  }
}

export default reduxForm({
  form: 'adminHomeCreateRepository',
  fields: [
    'display_name',
    'namespace',
    'reponame',
    'publisher.id',
    'publisher.name',
  ],
  validate: values => {
    const errors = {};
    if (!values.display_name) {
      errors.display_name = 'Required';
    }

    if (!values.namespace) {
      errors.namespace = 'Required';
    } else if (!/^[a-z0-9_]+$/.test(values.namespace)) {
      errors.namespace = 'Invalid namespace';
    }

    if (!values.reponame) {
      errors.reponame = 'Required';
    } else if (!/^[a-zA-Z0-9-_.]+$/.test(values.reponame)) {
      errors.reponame = 'Invalid reponame';
    }

    errors.publisher = {};
    if (!values.publisher.id) {
      errors.publisher.id = 'Required';
    } else if (!/^[a-z0-9_]+$/.test(values.publisher.id)) {
      errors.publisher.id = 'Invalid Publisher ID';
    }

    if (!values.publisher.name) {
      errors.publisher.name = 'Required';
    }

    return errors;
  },
})(AdminHomeCreateRepositoryForm);
