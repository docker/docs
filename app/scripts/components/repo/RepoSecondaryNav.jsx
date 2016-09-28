'use strict';
/**
 * TODO: Official Repo's route is linking to regular Tags route (/r/library/)
 *  Should create another route path /_/:reponame/tags to keep consistent.
 */
import styles from './RepoSecondaryNav.css';
import React, { Component, PropTypes } from 'react';
import { Link } from 'react-router';
const { string, bool } = PropTypes;
import { SecondaryNav } from 'dux';
import FA from '../common/FontAwesome';
import LiLink from '../common/LiLink';
var debug = require('debug')('RepoSecondaryNav::');

export default class RepoSecondaryNav extends Component {
  render() {
    const {user, splat, isOfficialRoute} = this.props;
    if (isOfficialRoute) {
      return (
        <SecondaryNav>
          <ul className='left'>
            <LiLink to={`/_/${splat}/`} key='repoDetailsInfo' onlyActiveOnIndex>Repo Info</LiLink>
            <LiLink to={`/r/${user}/${splat}/tags/`} key='repoDetailsTags'>Tags</LiLink>
          </ul>
        </SecondaryNav>
      );
    } else {
      let settingsLinks = null;
      let automated = null;
      let buildSettings = null;
      let repoSettings = null;

      if(this.props.isAutomated) {
        buildSettings = <LiLink to={`/r/${user}/${splat}/~/settings/automated-builds/`}
                                key='autobuildSettings'>Build Settings</LiLink>;
        automated = [(<LiLink to={`/r/${user}/${splat}/~/dockerfile/`} key='dockerfile'>Dockerfile</LiLink>),
          (<LiLink to={`/r/${user}/${splat}/builds/`} key='buildsMain'>Build Details</LiLink>)
        ];
      }

      if(this.props.canEdit) {
        repoSettings = (<LiLink to={`/r/${user}/${splat}/~/settings/`} key='repoSettingsMain' onlyActiveOnIndex>Settings</LiLink>);
        settingsLinks = [(buildSettings),
          (<LiLink to={`/r/${user}/${splat}/~/settings/collaborators/`} key='collaborators'>Collaborators</LiLink>),
          (<LiLink to={`/r/${user}/${splat}/~/settings/webhooks/`} key='webhooks'>Webhooks</LiLink>)
          ];
      }

      return (
        <SecondaryNav>
          <ul className='left'>
            <LiLink to={`/r/${user}/${splat}/`} key='repoDetailsInfo' onlyActiveOnIndex>Repo Info</LiLink>
            <LiLink to={`/r/${user}/${splat}/tags/`} key='repoDetailsTags'>Tags</LiLink>
              {automated}
              {settingsLinks}
              {repoSettings}
          </ul>
        </SecondaryNav>
      );
    }
  }
}

RepoSecondaryNav.propTypes = {
  JWT: string,
  user: string.isRequired,
  splat: string.isRequired,
  canEdit: bool,
  isAutomated: bool.isRequired
};
