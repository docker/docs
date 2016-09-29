'use strict';
import React from 'react';
import moment from 'moment';

//TODO: remove this after a better solution is found
var _getNotificationNameForID = function(user, notificationType) {
  //TODO: maybe create a link to the user's public profile page if any
    if(notificationType === 'trusted_build_fail') {
        return 'Automated build failure.';
    } else if (notificationType === 'new_repo_star') {
        return user + ' starred a repository';
    } else if (notificationType === 'new_repo_comment') {
        return user + ' added a new comment.';
    } else {
        return '';
    }
};

var ActivityFeedItem = React.createClass({
  displayName: 'ActivityFeed',
  getDefaultProps: function() {
    return {
      notification: 0,
      /*eslint-disable camelcase */
      last_occurence: '',
      /*eslint-enable camelcase */
      user: ''
    };
  },
  render: function() {
    return (
      <li>
        <div className="activity-feed-item clearfix">
          <p>{_getNotificationNameForID(this.props.user, this.props.notification)}</p>
          <small className="right">{moment(this.props.last_occurence).fromNow()}</small>
        </div>
      </li>);
  }
});

module.exports = ActivityFeedItem;
