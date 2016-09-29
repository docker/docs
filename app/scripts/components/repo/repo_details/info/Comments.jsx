'use strict';

import React, {
  PropTypes,
  createClass
} from 'react';
import { findDOMNode } from 'react-dom';
import _ from 'lodash';
import connectToStores from 'fluxible-addons-react/connectToStores';

import RepositoryCommentsStore from '../../../../stores/RepositoryCommentsStore';
import Comment from './Comment';
import classnames from 'classnames';
import addRepoComment from '../../../../actions/addRepoComment.js';
import getRepoComments from '../../../../actions/getRepoComments.js';
import deleteRepoComment from '../../../../actions/deleteRepoComment.js';
import Pagination from '../../../common/Pagination.jsx';
import { Button } from 'dux';

import styles from './Comments.css';

var debug = require('debug')('Comments');

let mkCommentElement = function(comment) {
  return (
    <Comment comment={comment}
             user={this.props.user}
             deleteRepoComment={this._deleteRepoComment(comment.id)}
             key={comment.id}/>
  );
};

var AddCommentButton = createClass({
  displayName: 'AddComment',
  contextTypes: {
    executeAction: PropTypes.func.isRequired
  },
  propTypes: {
    JWT: PropTypes.string.isRequired,
    commentAdded: PropTypes.func.isRequired
  },
  getInitialState: function() {
    return {
      commentBoxStatus: 'closed',
      commentBoxValue: ''
    };
  },
  _toggleComment(e){
    e.preventDefault();
    if (this.state.commentBoxStatus === 'open') {
      this.setState({
        commentBoxStatus: 'closed',
        commentBoxValue: ''
      });
    } else {
      this.setState({
        commentBoxStatus: 'open'
      }, () => {
        let node = findDOMNode(this.refs.commentTextArea);
        node.focus();
        node.scrollTop = node.scrollHeight;
      });
    }
  },
  _commentChange(e) {
    e.preventDefault();
    this.setState({
      commentBoxValue: e.target.value
    });
  },
  _leaveComment(e) {
    e.preventDefault();
    var comment = this.state.commentBoxValue;
    var repoShortName = this.props.namespace + '/' + this.props.name;
    this.setState({
      commentBoxValue: ''
    });
    this.context.executeAction(addRepoComment, {jwt: this.props.JWT, repoShortName: repoShortName, comment: comment});
    this.props.commentAdded();
  },
  render() {
    const addCommentClass = classnames({
      'row': true,
      [styles.addCommentBox]: true,
      [styles[this.state.commentBoxStatus]]: true
    });

    let commentButton = 'Add Comment';
    let buttonSize = 'small';
    let commentButtonIntent = 'primary';
    if (this.state.commentBoxStatus === 'open') {
      commentButton = 'Cancel';
      buttonSize = 'tiny';
      commentButtonIntent = 'alert';
    }
    return (
      <div>
        <div className="row">
          <div className={'columns large-12 ' + styles.button}>
            <Button intent={commentButtonIntent}
                    size={buttonSize}
                    onClick={this._toggleComment}>{commentButton}</Button>
          </div>
        </div>
        <div className={addCommentClass}>
          <div className="columns large-6">
            Add a Comment <span>(accepts markdown)</span>
            <form onSubmit={this._leaveComment}>
              <textarea value={this.state.commentBoxValue}
                        className={styles.addCommentBoxTextArea}
                        onChange={this._commentChange}
                        ref='commentTextArea'/>
              <Button type="submit" size='small'>Leave Comment</Button>
            </form>
          </div>
        </div>
      </div>
    );
  }
});


var RepoComments = createClass({
  displayName: 'RepoComments',
  contextTypes: {
    executeAction: PropTypes.func.isRequired
  },
  propTypes: {
    JWT: PropTypes.string,
    user: PropTypes.object,
    name: PropTypes.string.isRequired,
    namespace: PropTypes.string.isRequired,
    comments: PropTypes.object.isRequired
  },
  getInitialState() {
    return {
      currentPageNumber: 1
    };
  },
  _deleteRepoComment(commentid){
    return (e) => {
      e.preventDefault();
      var repoShortName = this.props.namespace + '/' + this.props.name;
      this.context.executeAction(deleteRepoComment, {
        jwt: this.props.JWT,
        repoShortName: repoShortName,
        commentid: commentid
      });
    };
  },
  _onChangePage(pageNumber) {
    this.setState({
      currentPageNumber: pageNumber
    }, function() {
      this.context.executeAction(getRepoComments, {JWT: this.props.JWT, namespace: this.props.namespace, repoName: this.props.name, pageNumber: pageNumber});
    });
  },
  _onCommentAdded() {
    this.setState({
      currentPageNumber: 1
    });
  },
  render() {

    // If logged in, display add comment form
    var maybeAddComment = null;
    if(this.props.JWT) {
      maybeAddComment = (
        <AddCommentButton JWT={this.props.JWT}
                          namespace={this.props.namespace}
                          name={this.props.name}
                          commentAdded={this._onCommentAdded}/>
      );
    }

    var maybePagination;
    if (this.props.comments.results && this.props.comments.results.length > 0) {
      maybePagination = (
        <div className='row'>
          <div className='large-12 columns'>
            <Pagination next={this.props.comments.next} prev={this.props.comments.prev}
                        onChangePage={this._onChangePage}
                        currentPage={this.state.currentPageNumber || 1}
                        pageSize={10} />
          </div>
        </div>
      );
    }

      return (
        <div>
          <div className="row">
            <div className="columns large-12">
              <hr />
            </div>
          </div>
          <div className="row">
            <div className="columns large-3">
              <div className={styles.header}>Comments ({this.props.comments.count})</div>
            </div>
          </div>
          {maybeAddComment}
          {maybePagination}
          <div className="row">
            <div className="columns large-12">
              {this.props.comments.results.map(mkCommentElement, this)}
            </div>
          </div>
          {maybePagination}
        </div>
    );
  }
});

export default connectToStores(RepoComments,
                               [
                                 RepositoryCommentsStore
                               ],
                               function({ getStore }, props) {
                                 return {comments: getStore(RepositoryCommentsStore).getState()};
                               });
