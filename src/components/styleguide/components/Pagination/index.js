import React, { Component } from 'react';
import { Input, Pagination } from 'common';
import asExample from '../../asExample';
import mdHeader from './header.md';
import mdApi from './api.md';

@asExample(mdHeader, mdApi)
export default class PaginationDoc extends Component {
  state = {
    currentPage: 1,
    lastPage: 10,
    maxVisible: 5,
  }

  onChangePage = (page) => {
    this.setState({
      currentPage: page,
    });
  }

  updateLastPage = (event) => {
    this.setState({
      lastPage: event.target.value,
    });
  }

  updateMaxVisible = (event) => {
    this.setState({
      maxVisible: event.target.value,
    });
  }

  render() {
    const { currentPage, lastPage, maxVisible } = this.state;
    const maxVisibleInput = (
      <div>
        <div>Set Max Visible</div>
        <Input
          id="max-visible-input"
          hintText="Max Visible"
          value={maxVisible}
          onChange={this.updateMaxVisible}
        />
      </div>
    );
    const lastPageInput = (
      <div>
        <div>Set Last Page</div>
        <Input
          id="lastpage-input"
          label="Last Page"
          value={lastPage}
          onChange={this.updateLastPage}
        />
      </div>

    );
    return (
      <div>
        {maxVisibleInput}
        {lastPageInput}
        <Pagination
          currentPage={currentPage}
          lastPage={lastPage}
          maxVisible={maxVisible}
          onChangePage={this.onChangePage}
        />
      </div>
    );
  }
}
