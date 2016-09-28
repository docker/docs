'use strict';

import React, { PropTypes } from 'react';
import classnames from 'classnames';

/*
 <div className='content active row' style={{paddingTop: 0}}>
 <div className='large-9 columns'>
 <div className='panel'>
 NO CONTENT HERE
 </div>
 </div>
 </div>
 */

export default React.createClass({
  displayName: 'DashboardTabs',
  propTypes: {
    myrepos: React.PropTypes.element.isRequired,
    contribs: React.PropTypes.element.isRequired,
    starred: React.PropTypes.element.isRequired,
    activity: React.PropTypes.element.isRequired
  },
  getInitialState() {
    return {
      tabType: 'myRepos'
    };
  },
  _getActiveTab() {
    switch(this.state.tabType) {
      case 'myRepos':
        //Show grid of my repositories
        return this.props.myrepos;
      case 'contrib':
        //Show grid of contributed repositories
        return this.props.contribs;
      case 'starred':
        //Show grid of my starred repositories
        return this.props.starred;
      case 'activity':
        //Show my activity feed
        return this.props.activity;
    }
  },
  _getTitleClassNames(type) {
    return classnames({
      'active': this.state.tabType === type,
      'tab-title': true
    });
  },
  _handleTabClick(type, evt) {
    evt.preventDefault();
    this.setState({
      tabType: type
    });
  },
  render() {
    var activeTab = this._getActiveTab();

    return (
      <div>
        <div className="row">
          <ul className="tabs large-12 columns">
            <li className={this._getTitleClassNames('myRepos')} onClick={this._handleTabClick.bind(null, 'myRepos')}>My Repositories</li>
            <li className={this._getTitleClassNames('contrib')} onClick={this._handleTabClick.bind(null, 'contrib')}>Contributed</li>
            <li className={this._getTitleClassNames('starred')} onClick={this._handleTabClick.bind(null, 'starred')}>Starred</li>
            <li className={this._getTitleClassNames('activity')} onClick={this._handleTabClick.bind(null, 'activity')}>Recent Activity</li>
          </ul>
        </div>
        <div>
          <div className="row">
            <div className="large-12 columns">
              <div className="tabs-content">
                  {activeTab}
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
});
