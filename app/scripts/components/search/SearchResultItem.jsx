'use strict';

import React from 'react';
import Badge from '../Badge.jsx';
import StatsComponent from '../StatsComponent.jsx';
var debug = require('debug')('COMPONENT:SearchResultItem');

//TODO: <FieldValueItems/> will go under the ul in item info, will be a bunch of key value pairs reused across
//TODO: Logged out views will have the `owner/reponame` (think about this)
//TODO: Star icon should be passed to badge as d-`iconname` where `d-` is for the docker font icons

var SearchResultItem = React.createClass({
  render: function() {
    var resultItem = this.props.resultItem;

    //Push badges based on result item
    var badges = [];
    var officialBadge = <li><Badge name="Official" category="dux-badge images" /></li>;
    var autobuildBadge = <li><Badge name="Autobuild" category="dux-badge builds" /></li>;

    // jscs:disable requireCamelCaseOrUpperCaseIdentifiers
    if (resultItem.is_official) {
      badges.push(officialBadge);
    } else if (resultItem.is_automated) {
      badges.push(autobuildBadge);
    }

    //TODO: repo_owner is null atm, since API performance degrades if we try to get it
    //<li><StatsComponent statsItemName="OWNER" value={resultItem.repo_owner}/></li>

    return (
      <li className="search-list-item row" onClick={this.props.onClick}>
        <div className="large-12 columns">
          <div className="row">
            <div className="logo large-2 columns">
                <i className="d-logo"></i>
            </div>
            <div className="search-item-basic-info large-2 columns">
              <div className="row">
                <div className="search-item-name large-12 columns">{resultItem.repo_name}</div>
              </div>
              <div className="row">
                <ul className="search-item-badges large-12 columns">
                  {badges}
                </ul>
              </div>
            </div>
            <div className="search-item-badges-stats large-4 columns">
              <ul className="small-block-grid-2">
                <li><StatsComponent statsItemName="DOWNLOADS" value={resultItem.pull_count}/></li>
                <li><StatsComponent statsItemName="STARS" value={resultItem.star_count}/></li>
              </ul>
            </div>
          </div>
          <div className="row search-bar-details">
            <div className="large-12 columns">
              <p>{resultItem.short_description}</p>
            </div>
          </div>
        </div>
      </li>
    );
    // jscs:enable
  }
});

module.exports = SearchResultItem;
