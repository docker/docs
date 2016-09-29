'use strict';
/**
TODO: UNUSED COMPONENTS. SHOULD REMOVE
*/
import React from 'react';
import { Link } from 'react-router';

var UserPage = React.createClass({
  getDefaultProps: function() {
    return {
      user: {},
      JWT: ''
    };
  },
  render: function() {
    return (
      <div>
          This will be the base wrapper of the 'Users' page where either your or another users profile will appear <br/>
          This will let you see your public facing page at /u/username/ too <br/>
          'Your' homepage/dashboard will live at /home/<br/>
          <RouteHandler JWT={this.props.JWT} user={this.props.user}/>
      </div>
    );
  }
});

var RootUser = React.createClass({
  render: function() {
    return (
      <div>
        <p/>
        This is root user page.<br/>
        When not looking at a specific user or an owned image<br/>
        This will show a list of repos/images owned by the root user <br/>
      <Link to={`/r/testing/1234/`}>This could be a image box of some sort</Link>
      </div>
    );
  }
});

var User = React.createClass({
    // This page should either be the users home page or the view page for other users
  render: function() {
    return (
      <div>
        <p/>
        This is the UID: {this.props.params.uid}<br/>
        This is main user page.<br/>
        This will show a list of repos/images owned by the user <br/>
        <RouteHandler />
      </div>
      );
  }
});

module.exports = {
  userpage: UserPage,
  rootuser: RootUser,
  user: User
};
