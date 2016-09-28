'use strict';

import React, { Component, PropTypes } from 'react';
import FA from 'components/common/fontAwesome';
import Dropdown from 'components/common/dropdown';
import css from 'react-css-modules';
import listStyles from './list.css';
import styles from './userRow.css';
import { updateOrganizationMember } from 'actions/organizations';
import { mapActions } from 'utils';
import { connect } from 'react-redux';

// TODO: Tooltip
/**
 * DeleteIcon is a helper component which shows the cross for removing the user
 * from the parent resource.
 */
let DeleteIcon = ({ userRecord, onDelete }) => (
  <span
    onClick={ () => onDelete(userRecord) }
    styleName='remove'>
    <FA icon='fa-times'/>
  </span>
);
DeleteIcon = css(DeleteIcon, listStyles);

@connect(
  (() => {})(),
  mapActions({
    updateOrganizationMember
  })
)
@css(styles)
export default class UserRow extends Component {

  static propTypes = {
    member: PropTypes.object,
    onDelete: PropTypes.func,
    isAdmin: PropTypes.bool,
    actions: PropTypes.object,
    params: PropTypes.object,
    canEdit: PropTypes.bool
  }

  updateMember = (evt) => {
    const {
      member
    } = this.props;

    // Quick fix by @camacho to resolve
    // https://github.com/docker/dhe-deploy/issues/2547
    const memberAttributes = this.getMemberAttributes(member);

    this.props.actions.updateOrganizationMember({
      name: this.props.params.org,
      member: memberAttributes.name,
      isAdmin: evt.target.value === 'admin' ? true : false,
      isPublic: true
    });
  }

  getMemberAttributes(member) {
    return member.member ? member.member : member;
  }

  render() {

    const {
      member,
      onDelete,
      canEdit
    } = this.props;

    const permissions = {
      defaultValue: member.isAdmin ? 'admin' : 'member',
      onChange: ::this.updateMember
    };

    const memberAttributes = this.getMemberAttributes(member);

    return (
      <tr key={ memberAttributes.name }>
        <td>
          <FA icon='fa-user'/> { memberAttributes.name }
        </td>
        <td>
          { memberAttributes.fullName }
        </td>
        <td>
            <span>
              <span styleName='select'>
                <Dropdown
                  disabled={ !canEdit }
                  values={ permissions }
                >
                  <option value='admin'>Admin</option>
                  <option value='member'>Member</option>
                </Dropdown>
              </span>
              {
                canEdit &&
                <DeleteIcon userRecord={ member } onDelete={ onDelete }/>
              }
            </span>
        </td>
      </tr>
    );
  }
}
