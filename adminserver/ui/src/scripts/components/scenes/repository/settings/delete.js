'use strict';

import React, { Component, PropTypes } from 'react';
const { func, instanceOf, object } = PropTypes;
import Immutable from 'immutable';
import { connect } from 'react-redux';
import DeleteModal from 'components/common/deleteModal';
import * as RepositoriesActions from 'actions/repositories';
import { mapActions, myGlobalAccessLevel } from 'utils';
import { connectModal } from 'components/common/modal';
//import InputLabel from 'components/common/inputLabel';
//import FontAwesome from 'components/common/fontAwesome';
//import CopyToClipboard from 'components/common/copyToClipboard';
import styles from '../repository.css';
import css from 'react-css-modules';
import Button from 'components/common/button';
import { getRawRepoState } from 'selectors/repositories';
import { getIsOrgOwner } from 'selectors/teams';

const mapRepositoryState = (state, props) => ({
  repositories: getRawRepoState,
  isOrgOwner: getIsOrgOwner(state, props)
});

@connectModal()
@connect(mapRepositoryState, mapActions(RepositoriesActions))
@css(styles, {allowMultiple: true})
export default class RepositoryDeleteForm extends Component {
  static propTypes = {
    repositories: instanceOf(Immutable.Map),
    isOrgOwner: object,
    params: object.isRequired,
    actions: object.isRequired,

    hideModal: func,
    showModal: func
  }

  confirmDelete() {
    this.props.showModal(
      (
        <DeleteModal
          resourceType='repository'
          resourceName={ this.props.params.name }
          onDelete={ ::this.deleteRepository }
          hideModal={ this.props.hideModal }/>
      )
    );
  }

  deleteRepository() {
    this.props.actions.deleteRepository({
      namespace: this.props.params.namespace,
      repo: this.props.params.name
    });
    this.props.hideModal();
  }

  onSubmit(data) {
    const { namespace, name } = this.props.params;
    return this.props.actions.updateRepository({namespace, repo: name}, data);
  }

  render() {
    const { namespace } = this.props.params;
    const { isOrgOwner } = this.props;
    const isAdminOrOrgOwner = (myGlobalAccessLevel() === 'admin') || !!isOrgOwner;
    const canDelete = isAdminOrOrgOwner || namespace === window.user.name;

    if (!canDelete) {
      return (<div />);
    }

    return (
      <div styleName='deleteRepository'>
        <div styleName='field-title'>Delete this repository</div>
        <Button variant='alert outline' onClick={ ::this.confirmDelete }>Delete</Button>
      </div>
    );
  }

}
