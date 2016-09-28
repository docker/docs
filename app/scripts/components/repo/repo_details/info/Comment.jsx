'use strict';

import styles from './Comment.css';
import React, { PropTypes, createClass } from 'react';
import { Link } from 'react-router';
import moment from 'moment';
import Card, { Block } from '@dux/element-card';
import Markdown from '@dux/element-markdown';
import FA from '../../../common/FontAwesome.jsx';
import { mkAvatarForNamespace } from 'utils/avatar';

var debug = require('debug')('RepositoryComment');

var CommentType = PropTypes.shape({
  updated_on: PropTypes.string.isRequired,
  created_on: PropTypes.string.isRequired,
  comment: PropTypes.string.isRequired,
  user: PropTypes.string.isRequired,
  id: PropTypes.number.isRequired
});

var UserType = PropTypes.shape({
  username: PropTypes.string.isRequired
});

var DeleteButton = createClass({
  propTypes: {
    deleteRepoComment: PropTypes.func.isRequired
  },
  render() {
    return (
      <div className={styles.commentDelete} onClick={this.props.deleteRepoComment}>
        <FA icon='fa-close'/>
      </div>
    );
  }
});

export default createClass({
  displayName: 'RepositoryComment',
  propTypes: {
    comment: CommentType,
    deleteRepoComment: PropTypes.func.isRequired,
    user: UserType
  },
  render() {
    var comment = this.props.comment;
    var maybeDeleteButton = null;

    if(this.props.user.username === this.props.comment.user) {
      maybeDeleteButton = <DeleteButton deleteRepoComment={this.props.deleteRepoComment}/>;
    }
    const avatar = mkAvatarForNamespace(this.props.comment.user);

    return (
      <Card>
        <Block>
          <div className={styles.commentHead}>
            <div>
              <img className={styles.avatar} src={avatar} />
            </div>
            <div className={styles.commentInfo}>
              <div className={styles.name}>{comment.user}</div>
              <div className={styles.time}>{moment(comment.created_on).fromNow()}</div>
            </div>
            {maybeDeleteButton}
          </div>
          <hr />
          <div className={styles.content}>
            <Markdown>{comment.comment}</Markdown>
          </div>
        </Block>
      </Card>
    );
  }
});
