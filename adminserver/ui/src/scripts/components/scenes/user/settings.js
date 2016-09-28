'use strict';

import React, { Component, PropTypes } from 'react';
const { object, func } = PropTypes;
import styles from 'components/scenes/settings/formstyle.css';
import Input from 'components/common/input';
import InputLabel from 'components/common/inputLabel';
import ToggleWithLabel from 'components/common/toggleSwitch/toggleWithLabel';
import Button from 'components/common/button';
import css from 'react-css-modules';
import { reduxForm } from 'redux-form';
import { selectUser } from 'selectors/users';
import { createStructuredSelector } from 'reselect';
import { connect } from 'react-redux';
import autoaction from 'autoaction';
import { getUser, updateUser, changeUserPassword, deleteUser } from 'actions/users';
import { mapActions } from 'utils';
import { currentUserSelector } from 'selectors/users';
import { connectModal } from 'components/common/modal';
import DeleteModal from 'components/common/deleteModal';

const mapState = createStructuredSelector({
  user: selectUser,
  currentUser: currentUserSelector
});

@connectModal()
@connect(mapState, mapActions({
  changeUserPassword,
  updateUser,
  deleteUser
}))
@autoaction({
  getUser: (props) => props.params.username
}, { getUser })
@reduxForm({
    form: 'newUser',
    fields: ['isActive', 'username', 'newPassword', 'oldPassword', 'isAdmin', 'fullName']
  }, (state, props) => {
    return {
      initialValues: {
        username: props.user.name,
        fullName: props.user.fullName,
        isAdmin: props.user.isAdmin
      }
    };
  }
)
@css(styles, { allowMultiple: true })
export default class UserSettings extends Component {

  static propTypes = {
    fields: object,
    handleSubmit: func,
    updateUI: func,
    user: object,
    params: object,
    actions: object,
    currentUser: object,
    hideModal: func,
    showModal: func
  }

  onSubmit = (data) => {
    if (data.newPassword && data.oldPassword) {
      this.props.actions.changeUserPassword({
        name: data.username,
        oldPassword: data.oldPassword,
        newPassword: data.newPassword
      });
    }

    if (data.fullName && this.props.fields.fullName.dirty) {
      this.props.actions.updateUser(
        data.username,
        {
          fullName: data.fullName
        });
    }

    if (this.props.fields.isAdmin.dirty) {
      this.props.actions.updateUser(
        data.username,
        {
          isAdmin: data.isAdmin
        });
    }
  }

  onDelete = () => {
    this.props.showModal(
      (
        <DeleteModal
          resourceType='user'
          resourceName={ this.props.user.name }
          onDelete={ ::this.deleteUser }
          hideModal={ this.props.hideModal }/>
      )
    );
  }

  deleteUser = () => {
    this.props.actions.deleteUser(this.props.user.name);
  }

  render () {

    const {
      fields: {
        isActive,
        username,
        newPassword,
        oldPassword,
        isAdmin,
        fullName
      },
      user,
      handleSubmit
    } = this.props;

    const canEdit = (this.props.currentUser.isAdmin || this.props.currentUser.name === user.name);

    return (
      <span>
      <form method='POST' onSubmit={ handleSubmit(::this.onSubmit) } styleName='pushed'>
        <input type='hidden' {...isActive} />
        <div styleName='formbox'>
          <h2>User info</h2>
          <div styleName='row'>
            <div styleName='halfColumn'>
              <InputLabel>Username</InputLabel>
              <Input
                disabled={ true }
                formfield={ username }
                type='text' />
            </div>
            <div styleName='halfColumn'>
              <InputLabel
                isOptional={ true }>Full Name</InputLabel>
              <Input
                disabled={ !canEdit }
                formfield={ fullName }
                type='text' />
            </div>
          </div>
          {
            canEdit &&
            <div styleName='row'>
              <div styleName='halfColumn'>
                <InputLabel
                  tip='To save a new password enter your current password.'
                >Current Password</InputLabel>
                <Input
                  formfield={ oldPassword }
                  type='password'/>
              </div>
              <div styleName='halfColumn'>
                <InputLabel>New Password</InputLabel>
                <Input
                  formfield={ newPassword }
                  type='password'/>
              </div>
            </div>
          }
          {
            this.props.currentUser.isAdmin &&
            <div styleName='row'>
              <div styleName='halfColumn'>
                <ToggleWithLabel
                  labelOptions={ {
                    inline: true,
                    labelText: 'Trusted Registry admin'
                  } }
                  initial={ user.isAdmin !== undefined ? user.isAdmin : false }
                  formField={ isAdmin }
                />
              </div>
            </div>
          }
          {
            canEdit &&
            <div styleName='textright'>
              <Button variant='primary' type='submit'>Save</Button>
            </div>
          }
        </div>
      </form>
      {
        this.props.currentUser.isAdmin &&
        <div styleName='row'>
          <div styleName='halfColumn formbox'>
            <h3>Delete</h3>
            <div id='delete-user'>
              Delete user (this cannot be undone)
              <div style={ {float: 'right'/*TODO*/} }>
                <Button variant='alert' onClick={ ::this.onDelete }>Delete</Button>
              </div>
            </div>
          </div>
        </div>
      }
      </span>
    );
  }
}

