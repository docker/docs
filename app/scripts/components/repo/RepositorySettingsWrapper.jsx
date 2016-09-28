'use strict';

import React, {
  PropTypes,
  createClass,
  cloneElement
} from 'react';

import RepoSecondaryNav from './RepoSecondaryNav';
import RouteNotFound404Page from '../common/RouteNotFound404Page';
import omit from 'lodash/object/omit';
var debug = require('debug')('RepoSettingsWrapper');

export default createClass({
  displayName: 'RepositorySettings',
  contextTypes: {
    getStore: PropTypes.func.isRequired
  },
  propTypes: {
    description: PropTypes.string.isRequired,
    fullDescription: PropTypes.string.isRequired,
    isPrivate: PropTypes.bool.isRequired,
    isAutomated: PropTypes.bool.isRequired,
    name: PropTypes.string.isRequired,
    namespace: PropTypes.string.isRequired,
    status: PropTypes.number.isRequired,
    canEdit: PropTypes.bool.isRequired
  },
  render() {
    if (this.props.canEdit) {
      const { namespace, name, canEdit, isAutomated } = this.props;
      return (
        <div>
          <RepoSecondaryNav user={namespace}
                            splat={name}
                            canEdit={canEdit}
                            isAutomated={isAutomated}/>
          {this.props.children && cloneElement(this.props.children, omit(this.props, 'children'))}
        </div>
      );
    } else {
      return (
        <RouteNotFound404Page />
      );
    }
  }
});
