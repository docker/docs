'use strict';

import React, { Component, PropTypes } from 'react';
const { bool, func, oneOf, instanceOf } = PropTypes;
import { RepositoryRecord } from 'records';
import { Link } from 'react-router';
import Tooltip from 'rc-tooltip';
import FontAwesome from 'components/common/fontAwesome';
import { humanPermissions } from 'utils';
import RepositoryName from 'components/common/repositoryName';
import TagLabel from 'components/common/tagLabel';
import PermissionDropdown from 'components/common/permissionDropdown';
import { permissionsByAccessLevel } from 'components/common/docSnippets';
import styles from './repositoryList.css';

export default class RepositoryListItem extends Component {

  static propTypes = {
    context: oneOf(['team', 'user']),

    // These are used to change a team's access level
    canEditPermissions: bool,
    onEditPermissions: func,
    canRemovePermissions: bool,
    onRemovePermissions: func,

    hidden: bool,
    repo: instanceOf(RepositoryRecord)
  }

  renderVisibility() {
    const { repo: { visibility } } = this.props;
    if (visibility === 'private') {
      return <span className={ styles.visibility }>{ visibility }</span>;
    }
  }

  /**
   * By default this renders a non-editable permission label.
   *
   * However, if:
   *  - The user is an admin or in the org 'owners' team
   *  - And this is a repository list within the org team
   * we should show a dropdown box which fires a callback for editing the team's
   * access level.
   *
   */
  renderPermission() {
    const {
      context,
      canEditPermissions,
      onEditPermissions,
      canRemovePermissions,
      onRemovePermissions,
      repo
    } = this.props;

    if (context === 'user') {
      // TODO need to show permissions on repos page too
      return null;
    }

    // If we don't have the current access level for the repository do not show
    // the edit dropdown
    if (!repo.accessLevel) {
      return null;
    }

    if (canEditPermissions === true) {
      return (
        <div className={ styles.permissionWrapper + ' ' + styles.editable }>
          <PermissionDropdown
            repo={ repo }
            accessLevel= { repo.accessLevel }
            onChange={ (evt) => onEditPermissions(repo, evt.target.value) }
            className={ styles.permission }
          />
          <div className={ styles.remove }>
            { canRemovePermissions && repo.accessLevel &&
              <Tooltip
                overlay={
                  <div style={ { maxWidth: 400 } }>
                    Revoke this team's <code>{ repo.accessLevel }</code> access to <code>{ repo.namespace }/{ repo.name }</code>
                  </div>
                }
                placement='right'
                align={ { overflow: { adjustY: 0 } } }
                trigger={ ['hover'] }>
                <FontAwesome
                  icon='fa-times-circle'
                  onClick={ () => onRemovePermissions(repo) }
                />
              </Tooltip>
            }
          </div>
        </div>
      );
    }

    return (
      <div className={ styles.permissionWrapper }>
        <span className={ styles.permission }>
          <TagLabel
            variant='accessLevel'
            tooltip={
              <div>
                This team has the following permissions on <code>{ repo.namespace }/{ repo.name }</code>:
                <p>{ permissionsByAccessLevel[repo.accessLevel] }</p>
              </div>
            } >
            { humanPermissions(repo.accessLevel) }
          </TagLabel>
        </span>
      </div>
    );
  }

  renderDescription() {
    let desc = this.props.repo.shortDescription;
    if (desc) {
      return (
        <div className={ styles.description }>
          { desc }
        </div>
      );
    }
  }

  render() {
    const {
      repo: {
        name,
        namespace,
        namespaceType
      }
    } = this.props;

    const permission = this.renderPermission();
    // Note z-index depends on list being sorted by repo id descending
    return (
      <div className={ styles.repositoryListItem }>
        <div className={ styles.mainContent }>
          <h2>
            <RepositoryName inlined={ true }>
              { namespaceType === 'user'
                  ? <span className={ styles.namespace }>{ namespace }</span>
                  : <Link className={ styles.namespace } to={ `/orgs/${namespace}` }>{ namespace }</Link> }
              <Link className={ styles.name } to={ `/repositories/${namespace}/${name}` }>{ name }</Link>
            </RepositoryName>
            { this.renderVisibility() }
          </h2>
          { this.renderDescription() }
        </div>
        { permission &&
          <div className={ styles.secondaryContent }>
            { permission }
          </div>
        }
      </div>
    );
  }
}
