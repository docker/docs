import React, { Component, PropTypes } from 'react';
import css from './styles.css';
import { sizes } from 'lib/constants';

const { TINY, SMALL, REGULAR, LARGE, XLARGE } = sizes;

import getIcon from './getIcon';

function normalizeName(name = '') {
  return name.replace('library/', '');
}

export default class DockerImage extends Component {
  static propTypes = {
    size: PropTypes.oneOf([TINY, SMALL, REGULAR, LARGE, XLARGE]),
    standalone: PropTypes.bool,
    imageName: PropTypes.string.isRequired,
    className: PropTypes.string,
  }

  render() {
    const {
      className = '',
      size = REGULAR,
      imageName,
      standalone,
    } = this.props;
    const Icon = getIcon(imageName);

    const styles = 'ddockerImage ' +
      ` ${className}` +
      ` ${css[size]}` +
      ` ${standalone ? css.standalone : ''}`;

    const displayText = standalone ? null :
      <span className={css.name}>{ normalizeName(imageName) }</span>;

    return (
      <div className={styles}>
        <span className={css.icon}>
          <Icon size={size} />
        </span>
        {displayText}
      </div>
    );
  }
}
