import React, { Component, PropTypes } from 'react';
import ReactDOM from 'react-dom';
import { uniqueId } from 'lodash';
import Tag from './Tag';
import './styles.css';

export default class TagList extends Component {
  static propTypes = {
    list: PropTypes.array.isRequired,
    size: PropTypes.number.isRequired,
  }

  state = { lastVisibleitem: 100, remaining: 0 };


  componentDidMount() {
    const listState = this.hideOverflowingTags();
    // MaximeHeckel 05-04-2016 we need the component mounted
    // to get the size of each tag

    // eslint-disable-next-line react/no-did-mount-set-state
    this.setState({
      lastVisibleitem: listState.lastVisibleitem,
      remaining: listState.remaining,
    });
  }

  calculateTotalSize() {
    const list = this.props.list;
    let tagListSize = 0;
    list.forEach((tag, index) => {
      tagListSize += ReactDOM.findDOMNode(
        this.refs[`tag_element_${index}`]
      ).offsetWidth;
    });
    return tagListSize;
  }

  hideOverflowingTags() {
    const totalSize = this.calculateTotalSize();
    const maxSize = this.props.size;
    const list = this.props.list;
    const remainingTextSize = 70;
    const padding = 4;
    const listState = {};
    let lastVisibleitem = this.props.list.length;
    let currentSize = 0;
    let remainings = 0;

    if (totalSize > maxSize) {
      list.forEach((tag, index) => {
        currentSize += ReactDOM.findDOMNode(
          this.refs[`tag_element_${index}`]
        ).offsetWidth + padding;
        if (currentSize < maxSize - remainingTextSize) {
          lastVisibleitem = index + 1;
        } else {
          remainings ++;
        }
      });
    }
    listState.remaining = remainings;
    listState.lastVisibleitem = lastVisibleitem;
    return listState;
  }

  render() {
    const {
      list,
      size,
    } = this.props;

    const styles = 'dTagList';
    const ulStyle = {
      width: size,
    };

    let Tags;

    return (
      <div className={styles}>
        <ul style={ulStyle}>
          {list.map((tag, index) => {
            if (index + 1 === this.props.list.length
                && this.state.remaining > 0) {
              Tags = (
                <Tag
                  key={uniqueId()}
                  name={ `... + ${this.state.remaining}  more` }
                  last={index + 1 === this.props.list.length}
                />
              );
            } else {
              Tags = (
                <Tag
                  key={uniqueId()}
                  ref={ `tag_element_${index}` }
                  name={tag}
                  last={index + 1 === this.props.list.length}
                />
              );
            }
            return Tags;
          })}
        </ul>
      </div>
    );
  }
}
