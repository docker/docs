import React, { Component, PropTypes } from 'react';
import styles from './styles.css';
import ReactModal from 'react-modal';
import classnames from 'classnames';
import { merge } from 'lodash';

const {
  any,
  string,
  bool,
  func,
  object,
} = PropTypes;

export default class Modal extends Component {
  static propTypes = {
    className: string,
    children: any,
    isOpen: bool,
    onRequestClose: func,
    onAfterOpen: func,
    style: object,
  }

  render() {
    const {
      className,
      children,
      isOpen,
      onRequestClose,
      onAfterOpen,
      style,
    } = this.props;
    const modalStyles = merge({
      overlay: {
        backgroundColor: 'rgba(0, 0, 0, .3)',
        zIndex: '1000',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
      },
      content: {
        borderRadius: '2px',
        position: 'relative',
        bottom: 'null',
        top: 'null',
        width: '740px',
        maxHeight: '80%',
        padding: 'null',
      },
    }, style);

    const modalClass = classnames({
      [className]: !!className,
      [styles.modalClass]: true,
    });

    return (
      <ReactModal
        className={modalClass}
        style={ modalStyles }
        isOpen={isOpen}
        onRequestClose={onRequestClose}
        onAfterOpen={onAfterOpen}
      >
        { children }
      </ReactModal>
    );
  }
}
