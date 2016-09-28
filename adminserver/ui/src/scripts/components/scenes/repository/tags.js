'use strict';

import React, { Component, PropTypes } from 'react';
const { func, object, string } = PropTypes;
import { connect } from 'react-redux';
import { getRepositoryTags, deleteRepoManifests } from 'actions/repositories';
import DeleteModal from 'components/common/deleteModal';
import FontAwesome from 'components/common/fontAwesome';
import { connectModal } from 'components/common/modal';
import Checkbox from 'components/common/checkbox';
import styles from './repository.css';
import consts from 'consts';
import Spinner from 'components/common/spinner';
import {
  getTagsForRepo,
  getRawRepoState,
  getAccessLevel
} from 'selectors/repositories';
import css from 'react-css-modules';
import Table from 'components/common/table';
import autoaction from 'autoaction';
import ui from 'redux-ui';

import { mapActions } from 'utils';

const mapRepositoryState = (state, props) => {
  return {
    repositories: getRawRepoState,
    tags: getTagsForRepo(state, props),
    accessLevel: getAccessLevel(state, props)
  };
};

@connectModal()
@connect(
  mapRepositoryState,
  mapActions({
    deleteRepoManifests
  })
)
@autoaction({
  getRepositoryTags: (props) => {
    return {
      namespace: props.params.namespace,
      repo: props.params.name
    };
  }
}, {
  getRepositoryTags
})
@ui({
  state: {
    selectedRows: []
  }
})
@css(styles)
export class RepositoryTagsTab extends Component {

  static propTypes = {
    params: object,
    repositories: object,
    actions: object.isRequired,

    hideModal: func,
    showModal: func,

    tags: object,
    accessLevel: string,
    ui: object,
    updateUI: func
  }

  onDeleteTags = (evt) => {
    evt.preventDefault();
    const { namespace, name: repo } = this.props.params;

    const {
      ui: {
        selectedRows
      }
    } = this.props;


    let deleteTags = () => {
      this.props.actions.deleteRepoManifests({
        namespace,
        repo,
        references: selectedRows.slice()
      });
      this.props.updateUI({
        selectedRows: []
      });
      this.props.hideModal();
    };

    this.props.showModal(
      (
        <DeleteModal
          resourceType='tag'
          resourceName={ (selectedRows.length === 1) ? 'This tag' : 'Multiple tags' }
          onDelete={ deleteTags }
          hideModal={ this.props.hideModal }/>
      )
    );
  };

  selectRow = (tag) => () => {
    const {
      selectedRows
    } = this.props.ui;
    // deselect
    if (selectedRows.includes(tag)) {
      return this.props.updateUI({
        selectedRows: selectedRows.filter((row) => {
          return row !== tag;
        })
      });
    }
    // select
    this.props.updateUI({
      selectedRows: selectedRows.concat([tag])
    });
  };

  selectAll = (tags) => () => {
    const {
      selectedRows
    } = this.props.ui;
    // deselect
    if (selectedRows.size === tags.size) {
      return this.props.updateUI({
        selectedRows: []
      });
    }
    // select
    this.props.updateUI({
      selectedRows: tags.toArray()
    });
  };

  makeDisplayTags = (tags) => {
    const displayTags = tags.valueSeq().toArray();

    const {
      selectedRows
    } = this.props.ui;

    return displayTags
      .map((tag, i) => {
        const inNotary = tag.get('inNotary');
        return (
          <tr key={ `${i}` } styleName='tagRow'>
            <td width={ 55 }>
              <Checkbox
                styleName='check'
                onChange={ ::this.selectRow(tag) }
                isChecked={ selectedRows.includes(tag) ? true : false }
              />
            </td>
            <td>{ tag.get('name') }</td>
            <td>
              { inNotary && tag.get('hashMismatch') && <IsOutOfDate /> }
              { inNotary && !tag.get('hashMismatch') && <IsSigned /> }
            </td>
          </tr>
        );
      });
  };

  render() {
    const {
      tags,
      accessLevel,
      ui: {
        selectedRows
      },
      params: {
        namespace,
        repo
      }
    } = this.props;

    const status = [[
      consts.repositories.GET_REPOSITORY_TAGS,
      namespace,
      repo
    ]];

    return (
      <Spinner loadingStatus={ status } styleName='tagsTable'>
        { tags.size > 0 ?
        <span>
        <div styleName='actionrow'>
          <div styleName='check'>
            <Checkbox
              ref='selectall'
              styleName='check'
              onChange={ ::this.selectAll(tags) }
            />
          </div>
          <div styleName='selectText'>
            Select all
            <hr />
            { /* <FontAwesome icon='fa-search'/> */ }
          </div>
          <div styleName='delete'>
            { (accessLevel === 'admin' || accessLevel === 'read-write') && selectedRows.length > 0 && (
              <button onClick={ ::this.onDeleteTags } styleName='deleteRepoButton'><FontAwesome icon='fa-trash-o'/> Delete</button>
            ) }
          </div>
        </div>
        <Table
          styleName='tagsTable'
          headers={ ['', 'Tags', ''] }>
          { (() => {
            return this.makeDisplayTags(tags);
          })() }
        </Table>
          </span>
        : <p>This repository has no tags.</p> }
      </Spinner>
    );
  }
}

@css(styles)
class IsSigned extends Component {
  render() {
    return (
      <span styleName='signed'>
          <FontAwesome icon='fa-check'/> signed
      </span>
    );
  }
}

@css(styles)
class IsOutOfDate extends Component {
  render() {
    return (
      <span styleName='outOfDate'>
          <FontAwesome icon='fa-exclamation'/> outdated
      </span>
    );
  }
}
