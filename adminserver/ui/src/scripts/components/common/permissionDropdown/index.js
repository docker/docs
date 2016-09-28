'use strict';

import React, { Component, PropTypes } from 'react';
const { func, instanceOf, string } = PropTypes;
import { RepositoryRecord } from 'records';
import Dropdown from 'components/common/dropdown';

/**
 * Renders a permission dropdown for a repository.
 *
 * When the permissions are changed the `props.onChange` function is called
 * with:
 *
 * - The first argument as the repo
 * - The second argument as the new accessLevel
 */
export default class PermissionDropdown extends Component {

  static propTypes = {
    repo: instanceOf(RepositoryRecord),
    accessLevel: string.isRequired,
    // Called when the access level changes
    onChange: func.isRequired
  }

  render() {
    const { accessLevel, ...props } = this.props;
    let dropdownProps = {
      ...props,
      defaultValue: accessLevel
    };

    return (
      <Dropdown values={ dropdownProps }>
        <option value='admin'>Admin</option>
        <option value='read-write'>Read & Write</option>
        <option value='read-only'>Read only</option>
      </Dropdown>
    );
  }
}
