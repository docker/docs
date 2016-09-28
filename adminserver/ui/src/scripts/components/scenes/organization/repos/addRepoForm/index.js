'use strict';

import React, { Component, PropTypes } from 'react';
const { bool, instanceOf, func, shape, object } = PropTypes;
// components
import Box from 'components/common/box';
import RadioGroup from 'components/common/radioGroup';
import AddExistingRepoForm from './addExistingRepoForm.js';
import AddNewRepoForm from './addNewRepoForm.js';
// actions
import autoaction from 'autoaction';
import {
  listOrganizations
} from 'actions/organizations';
// selectors
import {
  getPageOrg
} from 'selectors/organizations';
import {
  createStructuredSelector
} from 'reselect';
import { connect } from 'react-redux';
// misc
import ui from 'redux-ui';
import { OrganizationRecord } from 'records';
import styles from './addRepoForm.css';
import css from 'react-css-modules';

const mapState = createStructuredSelector({
  org: getPageOrg
});

/**
 * AddRepoForm is a parent container for a new or existing repo form.  This
 * calls all actions for data shared between both forms and passes down data.
 * We do this to prevent spamming the API every time someone switches between
 * a new and existing repo.
 *
 */
@connect(mapState)
@autoaction({
  listOrganizations: []
}, { listOrganizations })
@ui({
  state: {
    isExisting: true
  }
})
@css(styles)
export default class AddRepoForm extends Component {

  static propTypes = {
    ui: shape({ isExisting: bool }),
    updateUI: func,
    onHide: func,

    org: instanceOf(OrganizationRecord),
    params: object
  }

  updateChoice(choice) {
    this.props.updateUI('isExisting', choice === 'e');
  }

  render() {
    const {
      ui: { isExisting },
      org,
      onHide
    } = this.props;

    return (
      <Box>
        <h2>Add { isExisting ? 'an' : 'a' }
          <div styleName='choices'>
            <RadioGroup
              onChange={ ::this.updateChoice }
              initialChoice='e'
              choices={ [
                { label: 'Existing', value: 'e' },
                { label: 'New', value: 'n' }
              ] } />
          </div>
        repository</h2>
        { isExisting
            ? <AddExistingRepoForm
                params={ this.props.params }
                onHide={ onHide }
                org={ org } />
            : <AddNewRepoForm
                params={ this.props.params }
                onHide={ onHide }
                org={ org } /> }
      </Box>
    );
  }
}
