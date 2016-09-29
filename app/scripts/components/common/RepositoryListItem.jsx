'use strict';

import styles from './RepositoryListItem.css';
import React, { PropTypes, createClass, Component } from 'react';
let { string, number } = PropTypes;
import { Link } from 'react-router';
import classnames from 'classnames';
import FA from 'common/FontAwesome';
import numeral from 'numeral';
import { mkAvatarForNamespace, isOfficialAvatarURL } from 'utils/avatar';
import isFinite from 'lodash/lang/isFinite';

var debug = require('debug')('RepositoryListItem');

/* TODO: Should dedupe these render methods */

/**
 * RepositoryListItem has to handle two versions of a Repository data
 * structure. One comes from ElasticSearch, and one comes from Postgres
 *
 * ElasticSearch:
 *
 * -- Official
 * {
 *  "repo_name": "ubuntu",
 *  "is_offical": true,
 *  "is_automated": false,
 *  "short_description": "",
 *  "repo_owner": null,
 *  "pull_count": 9,
 *  "star_count": 2
 * }
 *
 * -- Normal
 * {
 *   "repo_name": "cpuguy83/ubuntu",
 *   "is_offical": false,
 *   "is_automated": false,
 *   "short_description": "ubuntu but more awesome",
 *   "repo_owner": null,
 *   "pull_count": 2,
 *   "star_count": 0
 * }
 *
 * Postgres:
 *
 * --Official
 * --Normal
 * {
 *   "last_updated": null,
 *   "pull_count": 9,
 *   "star_count": 2,
 *   "can_edit": true,
 *   "is_automated": false,
 *   "is_private": false,
 *   "full_description": null,
 *   "description": "",
 *   "status": 1,
 *   "namespace": "library",
 *   "name": "ubuntu",
 *   "user": "jlhawn"
 * }
 */

class Stats extends Component {
  static propTypes: {
    title: string,
    value: number,
    inBuckets: bool
  }

  formatNumber(n) {
    if (n < 1000) {
      return numeral(n).format('0a');
    }
    return numeral(n).format('0.0a').toUpperCase();
  }

  /**
   * Buckets to format the number:
   * 10M +
   * 5M +
   * 1M +
   * 500k +
   * 100k +
   * 50k +
   * 10k +
   */
  formatBucketedNumber(n) {
    if (n < 10000) {
      return this.formatNumber(n);
    } else if(n < 50000) {
      return '10K+';
    } else if(n < 100000) {
      return '50K+';
    } else if(n < 500000) {
      return '100K+';
    } else if(n < 1000000) {
      return '500K+';
    } else if(n < 5000000) {
      return '1M+';
    } else if(n < 10000000) {
      return '5M+';
    } else {
      return '10M+';
    }
  }

  render() {
    const { value, title, inBuckets } = this.props;

    let statValue;

    if (isFinite(value)) {
      if (inBuckets) {
        statValue = this.formatBucketedNumber(value);
      } else {
        statValue = this.formatNumber(value);
      }
    } else {
      return null;
    }

    return (
      <div className={styles.stats}>
        <div className={styles.labelValue}>
          <div className={styles.value}>{statValue}</div>
          <div className={styles.subLabel}>{title.toUpperCase()}</div>
        </div>
      </div>
    );
  }
}

var PostgresRepo = createClass({
  displayName: 'PostgresRepo',
  propTypes: {
    namespace: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired,
    status: PropTypes.number,
    description: PropTypes.string,
    fullDescription: PropTypes.string,
    isPrivate: PropTypes.bool,
    isAutomated: PropTypes.bool,
    isOfficial: PropTypes.bool,
    starCount: PropTypes.number,
    pullCount: PropTypes.number
  },
  render() {
    /**
     * Since Official Repos don't show their namespace (library),
     * we use this to chop that part off if we need to by redefining
     * this variable.
     */
    var repoDisplayName = this.props.namespace + '/' + this.props.name;
    var linkTo = `/r/${repoDisplayName}/`;

    /**
     * visibility can be 'official', 'public' or 'private', with an
     * appropriate class to match.
     */
    var visibility = null;
    var visibilityClasses = {};
    /**
     * autobuild is set if this repository is an automated build, otherwise
     * we rely on `null` to render nothing.
     */
    var autobuild = null;

    if(this.props.isOfficial || this.props.namespace === 'library') {

      visibility = 'official';
      visibilityClasses = classnames({
        [styles.official]: true
      });
      repoDisplayName = this.props.name || this.props.repo_name;
      linkTo = `/_/${this.props.name}/`;

    } else {

      visibility = 'private';
      if(!this.props.isPrivate) {
        visibility = 'public';
      }
      visibilityClasses = classnames({
        [styles.public]: !this.props.isPrivate,
        [styles.private]: this.props.isPrivate
      });

      if(this.props.isAutomated) {
        autobuild = <span className={styles.automated}> | automated build</span>;
      }
    }

    const avatar = mkAvatarForNamespace(this.props.namespace, this.props.name);
    const avatarClass = classnames({
      [styles.avatar]: true,
      [styles.officialAvatar]: isOfficialAvatarURL(avatar)
    });

    return (
      <li key={repoDisplayName}
          className={styles.repositoryListItem}>
        <Link to={linkTo} className={styles.flexible}>
        <div className={styles.head}>
          <div className={avatarClass}><img className={styles.img} src={avatar} /></div>
          <div className={styles.title}>
            <div className={styles.labels}>
            <div className={styles.repoName}>{repoDisplayName}</div>
              <span className={visibilityClasses}>{visibility}</span>
              {autobuild}</div>
          </div>
        </div>
        <Stats title='stars' value={this.props.starCount} inBuckets={false} />
        <Stats title='pulls' value={this.props.pullCount} inBuckets={!this.props.isPrivate} />
        <div className={styles.action}>
          <FA icon='fa-chevron-right' size='lg'/>
          <div className={styles.text}>DETAILS</div>
        </div>
        </Link>
      </li>
    );
  }
});


var ElasticRepo = createClass({
  displayName: 'ElasticRepo',
  propTypes: {
    repoName: PropTypes.string.isRequired,
    isOfficial: PropTypes.bool.isRequired,
    isAutomated: PropTypes.bool.isRequired,
    pullCount: PropTypes.number.isRequired,
    starCount: PropTypes.number.isRequired
  },
  render() {
    var repoDisplayName = this.props.repoName;
    var linkTo;

    /**
     * visibility can be 'official', 'public' or 'private', with an
     * appropriate class to match.
     */
    var visibility = null;
    var visibilityClasses = {};
    /**
     * autobuild is set if this repository is an automated build, otherwise
     * we rely on `null` to render nothing.
     */
    var autobuild = null;

    if(this.props.isOfficial) {

      visibility = 'official';
      visibilityClasses = classnames({
        [styles.official]: true
      });
      linkTo = `/_/${this.props.repoName}/`;

    } else {

      let [namespace, splat] = this.props.repoName.split('/');
      linkTo = `/r/${namespace}/${splat}/`;

      if(this.props.isPrivate) {
        visibility = 'private';
      } else {
        visibility = 'public';
      }
      visibilityClasses = classnames({
        [styles.public]: !this.props.isPrivate,
        [styles.private]: this.props.isPrivate
      });

      if(this.props.isAutomated) {
        autobuild = <span className={styles.automated}> | automated build</span>;
      }
    }

    const avatar = mkAvatarForNamespace(this.props.namespace, this.props.repoName);
    const avatarClass = classnames({
      [styles.avatar]: true,
      [styles.officialAvatar]: isOfficialAvatarURL(avatar)
    });

    return (
      <li key={repoDisplayName}
          className={styles.repositoryListItem}>
        <Link to={linkTo} className={styles.flexible}>
        <div className={styles.head}>
          <div className={avatarClass}><img src={avatar} /></div>
          <div className={styles.title}>
            <div className={styles.repoName}>{repoDisplayName}</div>
            <div className={styles.labels}>
              <span className={visibilityClasses}>{visibility}</span>
              {autobuild}</div>
          </div>
        </div>
        <Stats title='stars' value={this.props.starCount} inBuckets={false} />
        <Stats title='pulls' value={this.props.pullCount} inBuckets={!this.props.isPrivate} />
        <div className={styles.action}>
          <FA icon='fa-chevron-right' size='lg' />
          <div className={styles.text}>DETAILS</div>
        </div>
        </Link>
      </li>
    );
  }
});

export default class RepositoryListItem extends Component {
  render() {
    debug(this.props);
    if(this.props.repoName) {
      // Render a repo from ElasticSearch
      return (
        <ElasticRepo {...this.props} />
      );
    } else {
      // Render a repo from Postgres
      return (
        <PostgresRepo {...this.props} />
      );
    }
  }
}
