import React, { PropTypes, Component } from 'react';
import { ListItem, Avatar, Uptime, Markdown } from 'common';
import css from './styles.css';
const { number, string } = PropTypes;

export default class Comment extends Component {
  static propTypes = {
    id: number,
    user: string,
    comment: string,
    created_on: string,
    updated_on: string,
  }

  render() {
    const { comment, created_on, user } = this.props;
    const avatarUrl = [
      process.env.DOCKERCLOUD_HOST,
      '/v2/users/',
      user,
      '/avatar/',
    ].join('');

    return (
      <ListItem>
        <div className={css.header}>
          <Avatar size={40} src={avatarUrl} className={css.avatar} />
          <div className={css.user}>{user}</div>
          <Uptime since={created_on} />
        </div>
        <Markdown rawMarkdown={comment} />
      </ListItem>
    );
  }
}
