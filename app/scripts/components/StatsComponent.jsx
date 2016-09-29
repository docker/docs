'use strict';
import React from'react';

//TODO: add <i className="d-{icon}"></img> to component, currently just use a placeholder docker icon
//<ul className="no-bullet">
//  <li>{this.props.value}</li>
//  <li>{this.props.statsItemName}</li>
//</ul>

var StatsComponent = React.createClass({
  displayName: 'StatsComponent',
  propTypes: {
    value: React.PropTypes.string.isRequired,
    statsItemName: React.PropTypes.string.isRequired
  },
  render: function() {
    return (
    <div>
      <p>{this.props.value}</p>
      <p>{this.props.statsItemName}</p>
    </div>
    );
  }
});

module.exports = StatsComponent;
