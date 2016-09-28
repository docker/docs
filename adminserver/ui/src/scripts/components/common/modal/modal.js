'use strict';

import React, { Component, PropTypes } from 'react';
const { func, node } = PropTypes;
import cn from 'classnames';
import styles from './modal.css';

export default class Modal extends Component {

  static propTypes = {
    children: node.isRequired,
    component: node
  }

  static childContextTypes = {
    showModal: func.isRequired,
    hideModal: func.isRequired
  }

  state = {
    visible: false,
    component: null
  }

  getChildContext() {
    return {
      showModal: ::this.showModal,
      hideModal: ::this.hideModal
    };
  }

  componentDidMount() {
    document.addEventListener('keydown', ::this.handleKeypress, true);
  }

  handleKeypress(evt) {
    if (this.state.visible && evt.which === 27) {
      this.hideModal();
    }
  }

  handleBackdropClick(evt) {
    const { target } = evt;
    // Depending on the browser the target may be the container or backdrop
    if (target === this.refs.container || target === this.refs.backdrop) {
      this.hideModal();
    }
  }

  showModal(component) {
    let key = component.key;
    if (!key) {
      // If there is no key on the component, we assume we can set a random unique key for it
      // The purpose of this key is to force it to NOT reuse the old component if one exists
      key = Math.random();
      component = React.cloneElement(component, { key });
    }
    this.setState({ visible: true, component });
  }

  hideModal() {
    this.setState({ visible: false });
  }

  render() {
    const classes = cn({
      [styles.visible]: this.state.visible,
      [styles.modalBackdrop]: true
    });

    return (
      <div className={ styles.contentWrapper }>
        { this.props.children }
        <div className={ classes } onClick={ ::this.handleBackdropClick } ref='backdrop'>
          <div className={ styles.modalContainer } ref='container'>
            { this.state.component }
          </div>
        </div>
      </div>
    );
  }
}
