import React, { Component, PropTypes } from 'react';
import mergeClasses from 'classnames';
import LoadingIndicator from 'common/LoadingIndicator';
import css from './styles.css';


export default class FullscreenLoading extends Component {
  static propTypes = {
    className: PropTypes.string,
  }

  render() {
    const classNames = mergeClasses(this.props.className, css.main);

    return (
      <div className={classNames}>
        <div className={css.center} >
          <LoadingIndicator />
        </div>
      </div>
    );
  }
}
