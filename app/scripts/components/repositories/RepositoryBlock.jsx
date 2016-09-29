'use strict';

import React, { PropTypes } from 'react';
import { Link } from 'react-router';

export default React.createClass({
  displayName: 'RepositoryBlock',
  propTypes: {
    namespace: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired,
    status: PropTypes.number,
    description: PropTypes.string,
    fullDescription: PropTypes.string,
    isPrivate: PropTypes.bool,
    // trusted means Automated Build
    isTrusted: PropTypes.bool,
    isOfficial: PropTypes.bool,
    starCount: PropTypes.number,
    pullCount: PropTypes.number
  },
  getDefaultProps: function() {
    return {
      status: 0,
      description: '',
      fullDescription: '',
      isPrivate: true,
      isTrusted: false,
      isOfficial: false,
      starCount: 0,
      pullCount: 0
    };
  },
  render: function() {
    return (<li key={`${this.props.namespace}-${this.props.name}`} className="repository repo-border">
              <div className="repo-wrapper">
                <div className="row repo-header">
                  <div className="logo small-2 columns">
                    <i className="d-logo"></i>
                  </div>
                  <div className="header small-10 columns">
                    <h6 className="title">
                      <small><Link to={`/u/${this.props.namespace}`}>{this.props.namespace}</Link></small>
                      <br/>
                      <Link to={`/r/${this.props.namespace}/${this.props.name}/`}}>{this.props.name}</Link>
                      <br/>
                      <span>{this.props.isPrivate ? 'Private' : 'Public'}</span>
                      <span>{this.props.isTrusted ? ' | Automated Build' : ''}</span>
                    </h6>
                  </div>
                </div>
                <div className="row repo-content">
                  <p className="repo-short-description">{this.props.description}</p>
                </div>
                <div className="row repo-stats">
                  <p className="repo-stars small-4 columns">{this.props.starCount} STARS</p>
                  <p className="repo-downloads small-8 columns">{this.props.pullCount} DOWNLOADS</p>
                </div>
              </div>
            </li>
    );
  }
});
