import React, { Component, PropTypes } from 'react';
import { ListItem } from 'common';
import css from './styles.css';
import classnames from 'classnames';
import TagInfoArea from '../TagInfoArea';
const { bool, number, shape, string } = PropTypes;

export default class UnscannedTagRow extends Component {
  static propTypes = {
    canEdit: bool,
    tag: shape({
      name: string.isRequired,
      full_size: number.isRequired,
      id: number,
      last_updated: string,
    }),
  }

  render() {
    const { tag } = this.props;
    const { name } = tag;
    const classes = classnames({
      [css.listItem]: true,
    });
    return (
      <ListItem key={name}>
        <div className={classes}>
          <TagInfoArea tag={tag} />
        </div>
      </ListItem>
    );
  }
}
