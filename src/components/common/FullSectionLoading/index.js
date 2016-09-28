import React, { Component, PropTypes } from 'react';
import mergeClasses from 'classnames';
import LoadingIndicator from 'common/LoadingIndicator';
import css from './styles.css';


export default class FullSectionLoading extends Component {
  static propTypes = {
    className: PropTypes.string,
    title: PropTypes.string,
  }

  render() {
    const classNames = mergeClasses(this.props.className, css.main);
    let title;
    if (this.props.title) {
      title = (
        <div className={css.title}>{this.props.title}</div>
      );
    }
    return (
      <div className={classNames}>
        <div className={css.center} >
          {title}
          <LoadingIndicator />
        </div>
      </div>
    );
  }
}
