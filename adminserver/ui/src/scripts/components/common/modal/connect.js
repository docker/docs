'use strict';

import React, { Children, Component, PropTypes, cloneElement } from 'react';

export const connectModal = () => {
  return (DecoratedComponent) => class ModalDecorator extends Component {
    static displayName = 'DecoratedComponent'

    render() {
      const { props } = this;

      return (
        <ModalConnector>
          <DecoratedComponent { ...props } />
        </ModalConnector>
      );
    }
  };
};

class ModalConnector extends Component {
    static propTypes = {
      children: PropTypes.node.isRequired
    }


  static contextTypes = {
    showModal: PropTypes.func.isRequired,
    hideModal: PropTypes.func.isRequired
  }

  render() {
    const { children } = this.props;
    const props = {
      showModal: ::this.context.showModal,
      hideModal: ::this.context.hideModal
    };

    return (
      <div>
        { Children.map(children, (c) => cloneElement(c, props)) }
      </div>
    );
  }

}
