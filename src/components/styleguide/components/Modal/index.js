import React, { Component } from 'react';
import { Modal, Button, Card } from 'common';
import asExample from '../../asExample';
import mdHeader from './header.md';
import mdApi from './api.md';

@asExample(mdHeader, mdApi)
export default class ModalDoc extends Component {
  state = { isOpen: false }

  openModal = () => {
    this.setState({ isOpen: true });
  }

  closeModal = () => {
    this.setState({ isOpen: false });
  }

  render() {
    return (
      <div>
        <Button onClick={this.openModal}>Open Modal</Button>

        <Modal
          isOpen={this.state.isOpen}
          onRequestClose={this.closeModal}
        >
          <Card
            title="Hello World"
            ghost
          >
            <Button
              outlined
              variant="panic"
              onClick={this.closeModal}
            >Cancel</Button>
          </Card>
        </Modal>

      </div>
    );
  }
}
