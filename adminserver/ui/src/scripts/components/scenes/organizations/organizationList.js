'use strict';

import React, { Component, PropTypes } from 'react';
const { instanceOf, array, bool, string } = PropTypes;
import { Map } from 'immutable';
import { connect } from 'react-redux';
// Actions
import autoaction from 'autoaction';
import { listRepositories } from 'actions/repositories';
import { checkOrganizationMembership } from 'actions/organizations';
// Components
import { Link } from 'react-router';
import FontAwesome from 'components/common/fontAwesome';
import TagLabel from 'components/common/tagLabel';
import TagLabelList from 'components/common/labelList';
import Pagination from 'components/common/pagination';
// Selectors
import { currentUserSelector } from 'selectors/users';
import { createStructuredSelector } from 'reselect';
// Misc
import styles from './list.css';


/**
 * This expects to receive an array of organization names:
 * ['one', 'two']
 */
export default class OrganizationList extends Component {
  static propTypes = {
    organizations: array.isRequired,
    reposByOrg: instanceOf(Map),
    location: PropTypes.object
  }

  createRow(org) {
    const repos = this.props.reposByOrg.getIn(
      [org, 'entities', 'repo'],
      new Map()
    ).toArray();

    return (
      <OrganizationListItem
        key={ org }
        repos={ repos }
        orgName={ org } />
    );
  }

  render() {
    const { organizations } = this.props;
    return (
      <Pagination
        location={ location }
        className={ styles.orgList } pageSize={ 10 }>
        { organizations.map(::this.createRow) }
      </Pagination>
    );
  }
}

const mapState = createStructuredSelector({
  // TODO: Repalce these
  // isOrgOwner: isOrgOwnerSelector,
  // isUserMember: isUserMemberSelector,
  currentUser: currentUserSelector
});

@connect(mapState)
@autoaction({
  listRepositories: (props) => ({ namespace: props.orgName, limit: 10 }),
  checkOrganizationMembership: (props) => ({
    name: props.orgName,
    member: props.currentUser.name
  })
}, {
  listRepositories,
  checkOrganizationMembership
})
class OrganizationListItem extends Component {

  static propTypes = {
    isOrgOwner: bool,
    isUserMember: bool,
    orgName: string.isRequired,
    repos: array
  }

  renderLabels() {
    if (this.props.isOrgOwner) {
      return <TagLabel variant='role' tooltip='You are an owner of this organization'>Owner</TagLabel>;
    } else if (this.props.isUserMember) {
      return <TagLabel variant='role' tooltip='You are a member of this organization'>Member</TagLabel>;
    }
  }

  getRepos() {
    const { repos, orgName } = this.props;

    return repos.map(repo => {
      const { name } = repo;
      return <Link to={ `/repositories/${orgName}/${name}` }><FontAwesome icon='fa-book' />{ name }</Link>;
    });
  }

  render() {
    const { orgName } = this.props;
    return (
      <div key={ orgName } className={ styles.row }>
        <div className={ styles.orgGrid }>
          <FontAwesome icon='fa-users' className={ styles.orgAvatar } />
          <h2><Link to={ `/orgs/${orgName}` }>{ orgName }</Link></h2>
          { this.renderLabels() }
        </div>
        <div className={ styles.labelsGrid }>
          <TagLabelList variant='repository' labels={ this.getRepos() } />
        </div>
      </div>
    );
  }

}
