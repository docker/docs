'use strict';

import React, { Component, PropTypes } from 'react';
import FontAwesome from 'components/common/fontAwesome';
import TagLabel from 'components/common/tagLabel';
import styles from './labelList.css';

/**
 * Shows a list of TagLabels truncated to max number
 */
export default class TagLabelList extends Component {
  static propTypes = {
    className: PropTypes.string,
    max: PropTypes.number.isRequired,
    labels: PropTypes.array.isRequired, // Strings or could be elements like <Link>blah</Link>
    onRemoveLabel: PropTypes.func,
    variant: PropTypes.oneOf([ // Note: remember to change TagLabel also
        'accessLevel',
        'members',
        'repoAccessLevel',
        'repository',
        'repositoryTag',
        'role',
        'selectedMembers'
    ])
  }

  static defaultProps = {
    max: 10
  }

  take(arr) {
    return arr.slice(0, this.props.max);
  }

  renderList() {
    const { labels, onRemoveLabel, variant } = this.props;
    let list;

    list = this.take(labels).map((label, i) => (
      <li key={ i }>
        <TagLabel variant={ variant }>
          { label }
          { onRemoveLabel &&
            <FontAwesome
              icon='fa-times-circle'
              className={ styles.remove }
              onClick={ (e) => {
                onRemoveLabel(label); // TODO handle non string labels
                e.preventDefault();
              } }
            />
          }
        </TagLabel>
      </li>
    ));

    if (labels.length > this.props.max) {
      list.push(<li key='more'><TagLabel variant={ variant }>+{ labels.length - this.props.max } more...</TagLabel></li>);
    }

    return list;
  }

  render() {
    return (
      <ul className={ styles.labelList + ' ' + this.props.className }>
        { this.renderList() }
      </ul>
    );
  }
}
