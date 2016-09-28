import React, { Component, PropTypes } from 'react';
import css from './styles.css';
import cn from 'classnames';
import { VelocityTransitionGroup } from 'velocity-react';

const { node, string } = PropTypes;

/**
 * Expand produces a small expanding box for content. It has a title and main
 * body copy; a user clicks on the title to expand the body copy via an
 * animation.
 *
 * This differs from ExapndingContentBox as EpxandingContentBox uses links at
 * the bottom of content to show/hide content.
 */
export default class Expand extends Component {

  static propTypes = {
    title: string,
    titleClass: string,
    contentClass: string,
    children: node,
  }

  state = {
    isExpanded: false,
  }

  toggle = () => {
    this.setState({ isExpanded: !this.state.isExpanded });
  }

  render() {
    const { title, titleClass, contentClass, children } = this.props;
    const { isExpanded } = this.state;

    const titleStyles = cn({
      [css.title]: true,
      [titleClass]: !!titleClass,
      [css.expanded]: isExpanded,
    });

    return (
      <div>
        <div className={titleStyles} onClick={this.toggle}>
          {title}
        </div>
        <div className={contentClass}>
          <VelocityTransitionGroup
            enter={{ animation: 'slideDown' }}
            leave={{ animation: 'slideUp' }}
          >
            {isExpanded ? children : undefined}
          </VelocityTransitionGroup>
        </div>
      </div>
    );
  }
}
