'use strict';

import React, { Component, PropTypes } from 'react';
const { object, func, bool } = PropTypes;
// Components
import UserRow from './userRow';
import PaginatedTable from 'components/common/paginatedTable';
// CSS
import styles from './list.css';
import css from 'react-css-modules';
import { currentUserSelector } from 'selectors/users';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';

const mapState = createStructuredSelector({
  currentUser: currentUserSelector
});

/**
 * OrgUserList shows a list of users for the currently selected organization or
 * team.
 *UserRecords
 * @param members  array of  to show
 * @param allowDelete bool whether to show the corss for deletion
 * @param onDelete function which receives a UserRecord for deleting from the
 *                 parent resource (team/org/repo etc.)
 */
@connect(mapState)
@css(styles)
export default class OrgUserList extends Component {
  static propTypes = {
    members: object,
    onDelete: func,
    params: object,
    currentUser: object,
    canEdit: bool
  }

  makeUserRows = () => {

    const {
      members,
      onDelete,
      params,
      currentUser,
      canEdit
    } = this.props;

    return Object.keys(members).map((member) => {
      return (
        <UserRow
          key={ member }
          member={ members[member] }
          onDelete={ onDelete }
          userIsAdmin={ members[member].isAdmin || currentUser.isAdmin }
          canEdit={ canEdit }
          params={ params }
        />
      );
    });
  }

  render() {

    return (

    <PaginatedTable
      perPage={ 5 }
      headers={ ['Username', 'Full Name', ''] }
      rows={ this.makeUserRows() } />
    );
  }
}
