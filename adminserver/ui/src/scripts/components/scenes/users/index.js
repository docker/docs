'use strict';

import React, {
    Component,
    PropTypes
} from 'react';
const {
    object,
    func,
    string
} = PropTypes;
import FontAwesome from 'components/common/fontAwesome';
import ui from 'redux-ui';
import { VelocityComponent } from 'velocity-react';
// Components
import AddUserForm from './addUserForm';
import PaginatedTable from 'components/common/paginatedTable';
import Button from 'components/common/button';
import OrgList from './orgList.js';
// Misc
import styles from './styles.css';
import css from 'react-css-modules';
import { usersSelector } from 'selectors/users';
import { createStructuredSelector } from 'reselect';
import { connect } from 'react-redux';
import Spinner from 'components/common/spinner';
import { listUsers, createUser } from 'actions/users';
import consts from 'consts';
// actions
import autoaction from 'autoaction';
import { reduxForm, reset } from 'redux-form';
import { createValidator, required } from 'validation';
import { mapActions } from 'utils';
import { getAuthMethod } from 'selectors/settings';
import { saveAuthSettings } from 'actions/settings';
import { push } from 'react-router-redux';

import { currentUserSelector } from 'selectors/users';

const mapState = createStructuredSelector({
    authMethod: getAuthMethod,
    users: usersSelector,
    currentUser: currentUserSelector
});

@connect(
    mapState,
    mapActions({
        createUser,
        saveAuthSettings,
        reset,
        push
    })
)
@autoaction({
    listUsers: []
}, {
    listUsers
})
@ui({
    state: {
        // should we show the create user form?
        userFormVisible: false,
        // after creating a user, should we keep showing the create user form?
        createAnotherUser: false
    }
})
@reduxForm({
    form: 'newUser',
    fields: ['username', 'password', 'isAdmin', 'isActive'],
    validate: createValidator({
      username: [required],
      password: [required]
    }),
    initialValues: {
        isAdmin: false,
        isActive: true
    }
})
@css(styles, { allowMultiple: true })
export default class Users extends Component {

    static propTypes = {
        ui: object,
        updateUI: func,
        users: object,
        actions: object,
        fields: object,
        handleSubmit: func,
        authMethod: string,
        location: object,
        currentUser: object
    }

    static childContextTypes = {
        location: object
    }

    getChildContext() {
        return {
            location: this.props.location
        };
    }


    showUserEditForm = (evt) => {
        evt.preventDefault();
        this.props.updateUI({
            userFormVisible: !this.props.ui.userFormVisible
        });
    }

    cancelAddUser(evt) {
        evt.preventDefault();
        this.props.updateUI({
            userFormVisible: false
        });
    }

    // called when Save & Add Another is clicked
    shouldCreateAnotherUser = () => {
        this.props.updateUI({ createAnotherUser: true });
    }

    onSubmit = (data) => {
        data.name = data.username;
        data.type = 'user';
        data.isActive = true;
        this.props.actions.createUser(data);
        this.props.actions.reset('newUser');

        if (this.props.ui.createAnotherUser) {
            this.props.updateUI({
                createAnotherUser: false
            });
        } else {
            this.props.updateUI({
                userFormVisible: false
            });
        }
    }

    userDetail = (username) => () => {
        this.props.actions.push(`/users/${username}`);
    }

    render() {
      const status = [
        [consts.users.LIST_USERS]
      ];

      const {
        users,
        currentUser,
        ui: {
          userFormVisible
        }
      } = this.props;

      return (
        <Spinner loadingStatus={ status }>
            {
                currentUser.isAdmin &&
                <div styleName='textright'>
                    <Button id='new-user-button' variant='secondary' onClick={ ::this.showUserEditForm } disabled={ userFormVisible }>New user</Button>
                </div>
            }
            <VelocityComponent animation={ userFormVisible ? 'slideDown' : 'slideUp' } duration={ 250 }>
                <AddUserForm
                  submitHandler={ ::this.onSubmit }
                  addAnotherUserHandler={ ::this.shouldCreateAnotherUser }
                  cancel={ ::this.cancelAddUser }
                />
            </VelocityComponent>
            <PaginatedTable
              perPage={ 10 }
              headers={ ['Username', 'Full Name', 'Organizations', ''] }
              rows={ Object.keys(users).sort().map((name, i) => (
              <tr key={ i }>
                <td>{ name }</td>
                <td>{ users[name].fullName ? users[name].fullName : 'No name' }</td>
                <td><OrgList username={ name } /></td>
                <td className='icon' styleName='icon'><FontAwesome icon='fa-pencil' onClick={ ::this.userDetail(name) } /></td>
              </tr>
            )) } />
        </Spinner>
      );
    }
}
