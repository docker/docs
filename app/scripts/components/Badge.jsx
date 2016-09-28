'use strict';

import React from 'react';

var Badge = React.createClass({
    render: function() {
      var name = this.props.name;
      var category = this.props.category;
      return (
        <div className={category}>
          <div>{name}</div>
        </div>
      );
    }
});

module.exports = Badge;
