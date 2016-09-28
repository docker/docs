import React, { Component, PropTypes } from 'react';
import css from './styles.css';
import { sizes } from 'lib/constants';

const { TINY, SMALL, REGULAR, LARGE, XLARGE } = sizes;

import getIcon from './getIcon';

export default class NodeProvider extends Component {
  static propTypes = {
    size: PropTypes.oneOf([TINY, SMALL, REGULAR, LARGE, XLARGE]),
    standalone: PropTypes.bool,
    gray: PropTypes.bool,
    providerLabel: PropTypes.string,
    providerName: PropTypes.string.isRequired,
    className: PropTypes.string,
  }

  render() {
    const {
      className = '',
      size = REGULAR,
      gray,
      providerName,
      providerLabel,
      standalone,
    } = this.props;
    const Icon = getIcon(providerName);

    const styles = [
      'dnodeProvider',
      `${className}`,
      `${css[size]}`,
      `${gray ? css.gray : ''}`,
      `${standalone ? css.standalone : ''}`,
    ].join(' ');

    const displayText = standalone ? null :
      <span className={css.name}>{ providerLabel || providerName }</span>;

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
