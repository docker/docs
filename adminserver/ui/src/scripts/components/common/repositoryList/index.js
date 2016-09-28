'use strict';

import React, { Component, PropTypes } from 'react';
const { bool, oneOf, func, array, object } = PropTypes;
import Pagination from 'components/common/pagination';
import RepositoryListItem from './repositoryListItem';
import styles from './repositoryList.css';

export default class RepositoryList extends Component {
  static propTypes = {
    canEditPermissions: bool,
    canRemovePermissions: bool,

    // Whether we're rendering within context of an organization or user. This
    // will change the way our RepositoryListItem shows visibility statuses.
    context: oneOf(['team', 'user']),

    onEditPermissions: func,
    onRemovePermissions: func,

    repositories: array.isRequired,
    location: object
  }

  static defaultProps = {
    // By default render with a 'User' context
    context: 'user',
    canEditPermissions: true
  }

  render() {
    const {
      canEditPermissions,
      canRemovePermissions,
      context,
      onEditPermissions,
      onRemovePermissions,
      repositories
    } = this.props;

    const repos = repositories.map((repo) => {
      return (
        <RepositoryListItem
          canEditPermissions={ canEditPermissions }
          canRemovePermissions={ canRemovePermissions }
          context={ context }
          key={ repo.id }
          onEditPermissions={ onEditPermissions }
          onRemovePermissions={ onRemovePermissions }
          repo={ repo } />
      );
    });

    return (
      <Pagination
        location={ location }
        className={ styles.repositoryList }
        pageSize={ 10 }>
        { repos }
      </Pagination>
    );
  }
}
