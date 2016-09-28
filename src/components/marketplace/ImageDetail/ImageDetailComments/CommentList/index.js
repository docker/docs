import React, { PropTypes, Component } from 'react';
import { List } from 'common';
import { Link } from 'react-router';
import Comment from '../Comment';
import css from './styles.css';
import moment from 'moment';
const { arrayOf, bool, number, shape, string } = PropTypes;

export default class CommentList extends Component {
  static propTypes = {
    comments: arrayOf(shape({
      id: number,
      user: string,
      comment: string,
      created_on: string,
      updated_on: string,
    })),
    count: number,
    isPreview: bool,
    isFetching: bool,
    linkTo: string,
    page: number,
  }

  static defaultProps = {
    isPreview: false,
  }

  renderIsFetching() {
    return (
      <div>Fetching...</div>
    );
  }

  renderHeader() {
    const { isPreview, linkTo, comments, count } = this.props;
    const countText = typeof count === 'undefined' ? '' : `(${count})`;
    if (isPreview && !count) {
      return <div></div>;
    }
    if (isPreview) {
      return (
        <div className={css.previewHeader}>
          <div className={css.title}>Most Recent Comments</div>
          <Link to={linkTo} className={css.link}>
            {`All Comments ${countText}`}
          </Link>
        </div>
      );
    }
    let lastSubmitted;
    if (count) {
      // the first is the last updated
      const { updated_on } = comments[0];
      lastSubmitted = (
        <div className={css.lastSubmitted}>
          {`Last submitted ${moment(updated_on).fromNow()}`}
        </div>
      );
    }

    return (
      <div>
        <div className={css.allComments}>{`All Comments ${countText}`}</div>
        { lastSubmitted }
      </div>
    );
  }

  render() {
    // TODO Kristie 3/29/16 Handle Errors
    const { comments, isFetching } = this.props;
    if (isFetching) {
      return this.renderIsFetching();
    }
    let pageContent;
    if (comments && comments.length) {
      pageContent = comments.map((comment) =>
        <Comment key={comment.id} {...comment} />
      );
    } else {
      pageContent = <div>There are no comments yet.</div>;
    }
    return (
      <div>
        { this.renderHeader() }
        <List>
          { pageContent }
        </List>
      </div>

    );
  }
}
