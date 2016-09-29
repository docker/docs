'use strict';
import React, { Component, PropTypes } from 'react';
import { Link } from 'react-router';
import _ from 'lodash';
import AutobuildBlankSlate from './AutobuildBlankSlate.jsx';
import AutobuildSourceRepositoriesStore from '../../stores/AutobuildSourceRepositoriesStore';
import selectSourceRepoForAutobuild from '../../actions/selectSourceRepoForAutobuild';
import ListSelector from '../common/ListSelector.jsx';
import FilterBar from '../filter/FilterBar.jsx';
import FA from 'common/FontAwesome';
import connectToStores from 'fluxible-addons-react/connectToStores';
import styles from './LinkedAccountSourcesForm.css';

const debug = require('debug')('COMPONENT:LinkedAccountSourcesForm');
const { array, func } = PropTypes;

class LinkedAccountSourcesForm extends Component{

  static propTypes = {
    repos: array
  }

  static contextTypes = {
    executeAction: func.isRequired
  }

  state = {
    selectedRepos: [],
    selectedUserOrg: {}
  }

  _handleUserOrgClick = (userOrOrg) => {
    //Clicked on user or org, show the repositories under that user or organization
    this.setState({
      selectedUserOrg: userOrOrg,
      selectedRepos: userOrOrg.repo_list
    });
  }

  _handleRepoClick = (item, currentType) => {
    //Replace all "/" with "-" if they exist, so route is always valid
    item.name = item.name.replace(/[/]/g, '-');
    this.props.history.pushState(null,
      `/add/automated-build/${currentType.toLowerCase()}/form/${this.state.selectedUserOrg.name}/${item.name}/`,
      {namespace: this.props.params.userNamespace}
    );
    //Set the source repository in the store, so the autobuild configuration form can get to it
    this.context.executeAction(selectSourceRepoForAutobuild, item);
  }

  _makeUserOrgList = (list) => {
    var userAndOrgsList = list.map(function(item, idx) {
      var imgAvatar = item.avatar_url;
      var selectedArrow;
      if (this.state.selectedUserOrg.name === item.name) {
        selectedArrow = <span className={'right ' + styles.arrowSelect}><FA icon='fa-chevron-right' /></span>;
      }
      return (
        <li key={idx} ref={item.name} tabstop='0' onClick={this._handleUserOrgClick.bind(null, item)}>
          <div className='list-item'>
            <img src={imgAvatar} width='23' height='23' />&nbsp;&nbsp;{item.name}
            {selectedArrow}
          </div>
        </li>);
    }, this);
    return userAndOrgsList;
  }

  _makeReposList = (list) => {
    var links = [];
    var currentType = this.props.type;
    var _this = this;
    const namespace = this.state.selectedUserOrg.name;
    if (currentType) {
        links = list.map(function(item, idx) {
          return (
            <li key={idx} onClick={_this._handleRepoClick.bind(null, item, currentType)}>
              <a>{item.name}</a>
            </li>);
        });
    }
    return links;
  }

  _filterRepos = (query) => {
    //e.preventDefault();
    if (query) {
      this.setState({
        selectedRepos: _.filter(this.state.selectedRepos, function (repo) {
          return repo.name.indexOf(query) !== -1;
        })
      });
    } else {
      this.setState({
        selectedRepos: this.state.selectedUserOrg.repo_list
      });
    }
  }

  componentDidMount = () => {
    //Set the first element as selected if there is none selected
    if (_.isEmpty(this.state.selectedUserOrg) && this.props.repos) {
      var selected = this.props.repos[0];
      this.setState({
        selectedUserOrg: selected,
        selectedRepos: selected.repo_list
      });
    }
  }

  componentWillReceiveProps = (nextProps) => {
    //Set the first element as selected if there is none selected
    if (nextProps.repos && nextProps.repos.length > 0) {
      var selected = nextProps.repos[0];
      this.setState({
        selectedUserOrg: selected,
        selectedRepos: selected.repo_list
      });
    }
  }

  render() {
    if (this.props.repos) {
      var currentUserAndOrgs = this._makeUserOrgList(this.props.repos);
      var currentRepos = [];
      if (currentUserAndOrgs && currentUserAndOrgs.length > 0) {
        currentRepos = this._makeReposList(this.state.selectedRepos);
      }
      var filterBar = (<FilterBar items={this.state.selectedRepos} onFilter={this._filterRepos}/>);
      return (
        <div>
          <br />
          <div className='row'>
            <div className='columns large-5'>
              <ListSelector header='Users/Organizations' items={currentUserAndOrgs} />
            </div>
            <div className='columns large-7'>
              <ListSelector header={filterBar} items={currentRepos} />
            </div>
          </div>
        </div>
      );
    } else {
      const slateItem = (
        <div>
          <h2>
            <i className="fa fa-exclamation-triangle"></i>&nbsp;An error occurred while trying to connect to {this.props.type}.
          </h2>
          <div>Link your account <Link to='/account/authorized-services/'>here</Link>.</div>
        </div>
      );
      return (<AutobuildBlankSlate slateItems={slateItem} />);
    }
  }
}

export default connectToStores(LinkedAccountSourcesForm,
  [
    AutobuildSourceRepositoriesStore
  ],
  function({ getStore }, props) {
    return getStore(AutobuildSourceRepositoriesStore).getState();
  });
