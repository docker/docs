'use strict';
import React from 'react';
import _ from 'lodash';
import { FluxibleMixin } from 'fluxible';

var CreateButton = React.createClass({
  propTypes: {
    what: React.PropTypes.string.isRequired
  },
  getDefaultProps: function() {
    return {
      what: ''
    };
  },
  render: function() {
    return (
      <button className="create-object-btn">
        Create {this.props.what}<i className="d-add-3 link-icon"></i>
      </button>
    );
  }
});

module.exports = CreateButton;
