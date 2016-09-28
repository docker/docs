import React, { Component, PropTypes } from 'react';
import css from './styles.css';
import cn from 'classnames';
import { findDOMNode } from 'react-dom';
import { CopyIcon } from 'common';
import makeRepositoryId from 'lib/utils/repo-image-name';

const { bool, string } = PropTypes;

/**
 * CopyPullCommand is a pull command component with a copy button
 */
export default class CopyPullCommand extends Component {

  static propTypes = {
    hasInstruction: bool,
    codeClassName: string,
    // Either provide namespace and reponame
    namespace: string,
    reponame: string,
    tag: string,
    // or provide full command
    fullCommand: string,
    wrapperClassName: string,
  }

  static defaultProps = {
    hasInstruction: true,
  }

  copyContentsToClipboard = () => {
    const input = findDOMNode(this.refs.content);
    input.select();
    try {
      // Returns a boolean describing if the copy worked or not
      document.execCommand('copy');
    } catch (err) {
      // console.log('Error trying to copy.');
    }
  }

  render() {
    const {
      codeClassName,
      fullCommand,
      hasInstruction,
      namespace,
      reponame,
      tag,
      wrapperClassName,
    } = this.props;

    const codeStyles = cn({
      [css.pullCode]: true,
      [codeClassName]: !!codeClassName,
    });

    const instructionsStyles = cn({
      [css.pullCommand]: true,
      [wrapperClassName]: !!wrapperClassName,
    });

    const name = makeRepositoryId({ namespace, reponame });
    const displayTag = (!tag || tag === 'latest') ? '' : `:${tag}`;
    const copyButton = (
      <div className={css.copyButton} onClick={this.copyContentsToClipboard}>
        <CopyIcon />
      </div>
    );
    const pullCommand = fullCommand || `docker pull ${name}${displayTag}`;
    let instructions;
    if (hasInstruction) {
      instructions = (
        <div className={instructionsStyles}>
          Copy and paste to pull this image:
        </div>
      );
    }
    return (
      <div>
        {instructions}
        <div className={codeStyles}>
          <input
            className={css.input}
            readOnly
            ref="content"
            value={pullCommand}
          />
          {copyButton}
        </div>
      </div>
    );
  }
}
