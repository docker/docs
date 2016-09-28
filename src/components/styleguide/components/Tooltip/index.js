import React, { Component } from 'react';
import { Tooltip, DockerCloudIcon } from 'common';
import asExample from '../../asExample';
import mdHeader from './header.md';
import mdApi from './api.md';
import css from './styles.css';


@asExample(mdHeader, mdApi)
export default class TooltipDoc extends Component {
  render() {
    const basic = <div>Hello World</div>;
    const fancy = (
      <DockerCloudIcon size="large" variant="secondary" />
    );
    return (
      <div className="clearfix">
        <h3>Tooltip</h3>
        <Tooltip
          content={basic}
          placement="top"
        >
          <div className={css.trigger}>
            Tooltip on top
          </div>
        </Tooltip>
        <Tooltip
          content={basic}
          placement="bottom"
        >
          <div className={css.trigger}>
            Tooltip on bottom
          </div>
        </Tooltip>
        <Tooltip
          content={basic}
          placement="right"
        >
          <div className={css.trigger}>
            Tooltip on right
          </div>
        </Tooltip>
        <Tooltip
          content={basic}
          placement="left"
        >
          <div className={css.trigger}>
            Tooltip on left
          </div>
        </Tooltip>
        <Tooltip
          content={basic}
          placement="topLeft"
        >
          <div className={css.trigger}>
            Tooltip on top left
          </div>
        </Tooltip>
        <Tooltip
          content={basic}
          placement="bottomLeft"
        >
          <div className={css.trigger}>
            Tooltip on bottom left
          </div>
        </Tooltip>
        <Tooltip
          content={basic}
          placement="topRight"
        >
          <div className={css.trigger}>
            Tooltip on top right
          </div>
        </Tooltip>
        <Tooltip
          content={basic}
          placement="topLeft"
        >
          <div className={css.trigger}>
            Tooltip on top left
          </div>
        </Tooltip>
        <Tooltip
          content={fancy}
          placement="top"
        >
          <div className={css.trigger}>
            Any content
          </div>
        </Tooltip>
        <Tooltip
          content={fancy}
          placement="bottom"
          trigger={['click']}
        >
          <div className={css.trigger}>
            Click me!
          </div>
        </Tooltip>
        <Tooltip
          content={basic}
          placement="top"
          theme="dark"
        >
          <div className={css.trigger}>
            Dark Theme tooltip!
          </div>
        </Tooltip>
        <Tooltip
          content={fancy}
          placement="top"
          theme="dark"
        >
          <div className={css.trigger}>
            Dark Theme tooltip with SVG!
          </div>
        </Tooltip>
        <Tooltip
          content={fancy}
          placement="bottom"
          trigger={['focus']}
        >
          <input
            id="tooltip-input"
            placeholder="Focus on me!"
            className={css.input}
          />
        </Tooltip>
      </div>
    );
  }
}
