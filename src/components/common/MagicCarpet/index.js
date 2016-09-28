import React, { Component, PropTypes } from 'react';
import css from './styles.css';
import ReactModal from 'react-modal';
import classnames from 'classnames';
import { merge } from 'lodash';
import { CloseIcon } from 'common';
import { LARGE } from 'lib/constants/sizes';
import { DULL } from 'lib/constants/variants';

const {
  any,
  bool,
  func,
  node,
  object,
  string,
} = PropTypes;

export default class MagicCarpet extends Component {
  static propTypes = {
    children: any,
    className: string,
    isOpen: bool,
    onRequestClose: func,
    style: object,
    title: string,
    titleIcon: node,
  }

  render() {
    const {
      children,
      className,
      isOpen,
      onRequestClose,
      style,
      title,
      titleIcon,
    } = this.props;
    const modalStyles = merge({
      overlay: {
        backgroundColor: 'rgba(0, 0, 0, .3)',
        // This z-index can't be higher than the z-index of the modal, so that
        // the modal will still work on top of the magic carpet
        zIndex: '800',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        // transition: 'opacity 10s ease',
      },
      content: {
        borderRadius: '2px',
        position: 'relative',
        bottom: 'null',
        top: 'null',
        width: '100%',
        height: '100vh',
        padding: 'null',
      },
    }, style);

    const modalClass = classnames({
      [className]: !!className,
      [css.magicCarpet]: true,
    });

    return (
      <ReactModal
        className={modalClass}
        style={modalStyles}
        isOpen={isOpen}
        onRequestClose={onRequestClose}
      >
        <div className={css.header}>
          <div></div>
          <div className={css.titleWrapper}>
            <span className={css.icon}>{titleIcon}</span>
            <div className={css.title}>{title}</div>
          </div>
          <div className={css.aboveLoading} onClick={onRequestClose}>
            <CloseIcon
              className={css.close}
              size={LARGE}
              variant={DULL}
            />
          </div>
        </div>
        <div className="wrapped">
          {children}
        </div>
      </ReactModal>
    );
  }
}
