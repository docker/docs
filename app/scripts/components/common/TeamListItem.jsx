'use strict';

import React, { PropTypes, createClass } from 'react';
import { Link } from 'react-router';
import classnames from 'classnames';
var debug = require('debug')('TeamListItem');

let TeamShape = {
  description: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  id: PropTypes.number.isRequired
};

export default createClass({
  displayName: 'TeamListItem',
  propTypes: TeamShape,
  render() {

    return (
      <li key={this.props.name}
          className="team-item">
        <a href='#'>
          {this.props.name}
        </a>
      </li>
    );
  }
});
