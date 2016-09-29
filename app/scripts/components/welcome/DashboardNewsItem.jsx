'use strict';
const React = require('react');

//This is a placeholder for a news item that is shown in the dashboard. Relevant to anything that needs attention.

var DashboardNewsItem = React.createClass({
  render: function() {
    return (
      <div className="large-3 columns">
        <h6>Run Docker in your network.</h6>
        <div>
          <a href="#">Get a free 30 day trial!</a><i className="icon-right-arrow" />
        </div>
      </div>
    );
  }
});

module.exports = DashboardNewsItem;
