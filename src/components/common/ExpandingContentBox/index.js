import React, { Component, PropTypes } from 'react';
import css from './styles.css';
import classnames from 'classnames';
// import { VelocityComponent } from 'velocity-react';

const { node, number, string } = PropTypes;

export default class ExpandingContentBox extends Component {
  static propTypes = {
    children: node,
    height: number,
    // these classes will be applied to the parent container
    showClass: string,
    hideClass: string,
    // the text to use for show/hide links
    showText: string,
    hideText: string,
  }

  static defaultProps = {
    height: 200,
    showText: 'Show more',
    hideText: 'Show less',
  }

  state = {
    isExpanded: false,
  }

  toggle = () => {
    this.setState({ isExpanded: !this.state.isExpanded });
  }

  render() {
    const {
      children,
      height,
      showClass,
      hideClass,
      showText,
      hideText,
    } = this.props;
    const { isExpanded } = this.state;

    let styles = {};
    if (!isExpanded) {
      styles.maxHeight = height;
    }

    const className = classnames({
      [css.content]: true,
      [showClass]: isExpanded && !!showClass,
      [hideClass]: !isExpanded && !!hideClass,
    });

    return (
      <div>
        <div style={styles} className={className}>
          {children}
        </div>
        <div className={css.toggleText} onClick={this.toggle}>
          {isExpanded ? hideText : showText}
        </div>
      </div>
    );
  }
}
