'use strict';

import React, { Component, PropTypes } from 'react';
import Input from 'components/common/input';
import InputLabel from 'components/common/inputLabel';
import ToggleWithLabel from 'components/common/toggleSwitch/toggleWithLabel';
import Button from 'components/common/button';
import { reduxForm } from 'redux-form';
import { checkRequiredFields, reportMissingFields } from 'validation';
import css from 'react-css-modules';
import styles from 'components/scenes/settings/formstyle.css';

@reduxForm({
  form: 'newUser',
  fields: ['isActive', 'username', 'password', 'isAdmin', 'fullName'],
  validate: (values) => {
    return reportMissingFields(checkRequiredFields(['username', 'password'], values), {});
  },
  initialValues: {
    isAdmin: false,
    isActive: true
  }
})
@css(styles)
export default class AddUserForm extends Component {

  static propTypes = {
    fields: PropTypes.object,
    submitHandler: PropTypes.func,
    handleSubmit: PropTypes.func,
    updateUI: PropTypes.func,
    cancel: PropTypes.func,
    addAnotherUserHandler: PropTypes.func
  }

  render () {
    const {
      fields: {
        isActive,
        username,
        password,
        isAdmin,
        fullName
      },
      handleSubmit,
      cancel,
      addAnotherUserHandler
    } = this.props;
    return (
       <form method='POST' onSubmit={ handleSubmit(this.props.submitHandler) } styleName='pushed' id='new-user-form'>
         <input type='hidden' {...isActive} />
        <div styleName='formbox'>
          <h2>New user</h2>
          <div styleName='row'>
            <div styleName='halfColumn'>
              <InputLabel>Username</InputLabel>
              <Input
                formfield={ username }
                type='text' />
            </div>
            <div styleName='halfColumn'>
              <InputLabel>Password</InputLabel>
              <Input
                formfield={ password }
                type='password' />
            </div>
          </div>
          <div styleName='row'>
            <div styleName='halfColumn'>
              <InputLabel
                isOptional={ true }>Full Name</InputLabel>
              <Input
                formfield={ fullName }
                type='text' />
            </div>
            <div styleName='halfColumn'>
              <ToggleWithLabel
                labelOptions={ {
                  inline: true,
                  labelText: 'Trusted Registry admin'
                } }
                initial={ false }
                formField={ isAdmin }
              />
            </div>
          </div>
          <div styleName='textright'>
            <Button variant='primary simple' onClick={ cancel }>Cancel</Button>
            <Button variant='primary outline' onClick={ addAnotherUserHandler }>Save &amp; create another</Button>
            <Button variant='primary' type='submit'>Save</Button>
          </div>
        </div>
      </form>
    );
  }
}
