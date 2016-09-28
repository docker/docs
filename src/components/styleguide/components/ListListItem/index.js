import React, { Component } from 'react';
import { findWhere } from 'lodash';
import asExample from '../../asExample';
import mdHeader from './header.md';
import mdApi from './api.md';
import { List, ListItem } from 'common';

@asExample(mdHeader, mdApi)
export default class ListListItemDoc extends Component {
  constructor(props) {
    super(props);
    this.state = {
      listWithSelectionItems: [{
        uuid: 'uuid-1',
        selected: false,
        content: `
          I got my head checked
          By a jumbo jet
          It wasn't easy
          But nothing i-is
          No`,
      }, {
        uuid: 'uuid-2',
        selected: true,
        content: `
          I got my head down
          When I was young
          It's not my problem
          It's not my problem
        `,
      }, {
        uuid: 'uuid-3',
        disabled: true,
        content: `
          Woooo hooo!
          Woooo hooo!
          Woooo hooo!
        `,
      }],
    };
  }

  resolveItemSelection = (ev, selected, uuid) => {
    const el = findWhere(this.state.listWithSelectionItems, { uuid });
    el.selected = selected;
    this.setState(el);
  }

  render() {
    const resolve = this.resolveItemSelection;
    const selectionList = this.state.listWithSelectionItems;
    /* eslint-disable max-len */
    return (
      <div>
        <h3>Simple List</h3>
        <List>
          <ListItem>Provided by <a href="http://www.google.com/?q=i%27m+feeling+curious" target="_blank">www.google.com/?q=i%27m+feeling+curious</a></ListItem>
          <ListItem>The name “hamburger” actually came from Hamburg, the second largest city in Germany. In the late 1700s, sailors who traveled between Hamburg and New York City often ate hard slabs of salted minced beef, which they called “Hamburg steak.”</ListItem>
          <ListItem>An average professional football game lasts 3 hours and 12 minutes, but if you tally up the time when the ball is actually in play, the action amounts to a mere 11 minutes. Part of the discrepancy has to do with the basic rules of American football.</ListItem>
          <ListItem>The first woman to win a Nobel Prize was Marie Curie, who won the Nobel Prize in Physics in 1903 with her husband, Pierre Curie, and Henri Becquerel. Curie is also the only woman to have won multiple Nobel Prizes; in 1911, she won the Nobel Prize in Chemistry.</ListItem>
          <ListItem>The rule of thumb used by most antique dealers is that anything about 100 years or older is an antique. Items that are old, but not quite that old, are called vintage.</ListItem>
        </List>
        <br />
        <h3>Selectable List</h3>
        <List selectable hover>
          <div>Selected: {selectionList.filter(({ selected }) => selected).map(({ uuid }) => uuid).join(' ')}</div>
          {selectionList.map(({ uuid, content, selected, disabled }) => {
            return (
              <ListItem
                id={uuid}
                key={uuid}
                selected={!!selected}
                disabled={!!disabled}
                onSelect={resolve}
              >
                {content}
              </ListItem>
            );
          })}
        </List>
      </div>
    );
    /* eslint-enable max-len */
  }
}
