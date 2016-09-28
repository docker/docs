/* eslint-disable react/sort-comp */
/* eslint-disable no-underscore-dangle */
import React, { Component, PropTypes } from 'react';
import { SMALL } from 'lib/constants/sizes';
import { SearchIcon } from '../Icon';
const scrollIntoView = require('dom-scroll-into-view');
import classnames from 'classnames';
const { any, array, bool, func, node, object, shape, string } = PropTypes;
import css from './styles.css';

const _shouldItemRender = () => { return true; };

/* This is a heavily modified version of react-autocomplete
 * to enable the the search use-case and our specific needs.
 * https://github.com/reactjs/react-autocomplete is the original
 *
 * TODO Kristie 5/11/16 Make this from a fork of the original
 * repo and NPM publish
 */
export default class AutocompleteSearchBar extends Component {
  static propTypes = {
    autoSelectMatch: bool,
    classNames: shape({
      icon: string,
      input: string,
      menu: string,
      menuTitle: string,
      wrapper: string,
    }),
    getItemValue: func.isRequired,
    id: string.isRequired,
    inputProps: object,
    isLoading: bool,
    items: array,
    labelText: string,
    menuTitle: node,
    onChange: func.isRequired,
    onSelect: func.isRequired,
    onSubmit: func,
    placeholder: string,
    renderItem: func.isRequired,
    renderMenu: func,
    shouldItemRender: func,
    sortItems: func,
    value: any,
    wrapperProps: object,
  }

  static defaultProps = {
    // automatically select a match with autocomplete
    autoSelectMatch: false,
    classNames: {},
    inputProps: {},
    isLoading: false,
    items: [],
    labelText: '',
    shouldItemRender: _shouldItemRender,
    value: '',
    wrapperProps: {},
  }

  state = {
    isOpen: false,
    highlightedIndex: null,
  }

  componentWillMount() {
    this.id = `autocomplete-${this.props.id}`;
    this._ignoreBlur = false;
    this._performAutoCompleteOnUpdate = false;
    this._performAutoCompleteOnKeyUp = false;
  }

  componentWillReceiveProps() {
    this._performAutoCompleteOnUpdate = true;
  }

  componentDidUpdate(prevProps, prevState) {
    if (this.state.isOpen && !prevState.isOpen) {
      this.setMenuPositions();
    }
    if (this.state.isOpen && this._performAutoCompleteOnUpdate) {
      this._performAutoCompleteOnUpdate = false;
      this.maybeAutoCompleteText();
    }

    this.maybeScrollItemIntoView();
  }

  // Define Key Down Handlers

  _onEnter = () => {
    if (this.state.highlightedIndex === null) {
      // input has focus but no menu item is selected + enter is hit ->
      // close the menu, highlight whatever's in input
      this.setState({
        isOpen: false,
      }, () => {
        this.refs.input.select();
        // call an OnSubmit function to submit the search
        if (this.props.onSubmit) {
          this.props.onSubmit(this.refs.input.value);
        }
      });
    } else {
      // text entered + menu item has been highlighted + enter is hit ->
      // update value to that of selected menu item, close the menu
      const item = this.getFilteredItems()[this.state.highlightedIndex];
      const value = this.props.getItemValue(item);
      this.setState({
        isOpen: false,
        highlightedIndex: null,
      }, () => {
        // this.refs.input.focus() // TODO: file issue
        this.refs.input.setSelectionRange(
          value.length,
          value.length
        );
        this.props.onSelect(value, item);
      });
    }
  }

  _onArrowDown = (event) => {
    event.preventDefault();
    const { highlightedIndex } = this.state;
    const index = (
      highlightedIndex === null ||
      highlightedIndex === this.getFilteredItems().length - 1
    ) ? 0 : highlightedIndex + 1;
    this._performAutoCompleteOnKeyUp = true;
    this.setState({
      highlightedIndex: index,
      isOpen: true,
    });
  }

  _onArrowUp = (event) => {
    event.preventDefault();
    const { highlightedIndex } = this.state;
    const index = (
      highlightedIndex === 0 ||
      highlightedIndex === null
    ) ? this.getFilteredItems().length - 1 : highlightedIndex - 1;
    this._performAutoCompleteOnKeyUp = true;
    this.setState({
      highlightedIndex: index,
      isOpen: true,
    });
  }

  _onEscape = () => {
    this.setState({
      highlightedIndex: null,
      isOpen: false,
    });
  }

  getKeyDownHandlers = () => {
    return {
      ArrowDown: this._onArrowDown,
      ArrowUp: this._onArrowUp,
      Enter: this._onEnter,
      Escape: this._onEscape,
    };
  }

  getFilteredItems = () => {
    let items = this.props.items;

    if (this.props.shouldItemRender) {
      items = items.filter((item) => (
        this.props.shouldItemRender(item, this.props.value)
      ));
    }

    if (this.props.sortItems) {
      items.sort((a, b) => (
        this.props.sortItems(a, b, this.props.value)
      ));
    }

    return items;
  }

  maybeAutoCompleteText = () => {
    const { autoSelectMatch } = this.props;
    const items = this.getFilteredItems();
    if (!autoSelectMatch || !this.props.value || items.length === 0) {
      return;
    }
    const { highlightedIndex } = this.state;
    const matchedItem = highlightedIndex !== null ?
      items[highlightedIndex] : items[0];
    const itemValue = this.props.getItemValue(matchedItem);
    const itemValueDoesMatch = (itemValue.toLowerCase().indexOf(
      this.props.value.toLowerCase()
    ) === 0);
    if (itemValueDoesMatch) {
      const inputNode = this.refs.input;
      const setSelection = () => {
        inputNode.value = itemValue;
        inputNode.setSelectionRange(this.props.value.length, itemValue.length);
      };
      if (highlightedIndex === null) {
        this.setState({ highlightedIndex: 0 }, setSelection);
      } else {
        setSelection();
      }
    }
  }

  setMenuPositions = () => {
    const inputNode = this.refs.input;
    const rect = inputNode.getBoundingClientRect();
    const computedStyle = global.window.getComputedStyle(inputNode);
    const marginBottom = parseInt(computedStyle.marginBottom, 10) || 0;
    const marginLeft = parseInt(computedStyle.marginLeft, 10) || 0;
    const marginRight = parseInt(computedStyle.marginRight, 10) || 0;
    this.setState({
      menuTop: rect.bottom + marginBottom,
      menuLeft: rect.left + marginLeft,
      menuWidth: rect.width + marginLeft + marginRight,
    });
  }

  highlightItemFromMouse = (index) => {
    this.setState({ highlightedIndex: index });
  }

  selectItemFromMouse = (item) => {
    const value = this.props.getItemValue(item);
    this.setState({
      isOpen: false,
      highlightedIndex: null,
    }, () => {
      this.props.onSelect(value, item);
      this.refs.input.focus();
      this.setIgnoreBlur(false);
    });
  }

  setIgnoreBlur = (ignore) => {
    this._ignoreBlur = ignore;
  }

  _renderMenu = (items, value, style) => {
    const { classNames, menuTitle, renderMenu } = this.props;
    if (renderMenu) {
      return renderMenu(items, value, style);
    }
    let title;
    if (items && items.length && menuTitle) {
      const menuTitleClasses = classnames('autocompleteMenuTitle', {
        [css.menuTitle]: true,
        [classNames.menuTitle]: !!classNames.menuTitle,
      });
      title = <div className={menuTitleClasses}>{menuTitle}</div>;
    }
    const menuClasses = classnames('autocompleteMenu', {
      [css.menu]: true,
      [classNames.menu]: !!classNames.menu,
    });
    // TODO Kristie 5/11/16 implement loading state
    return (
      <div className={menuClasses}>
        {title}
        {items}
      </div>
    );
  }

  maybeScrollItemIntoView = () => {
    if (this.state.isOpen && this.state.highlightedIndex !== null) {
      const itemNode = this.refs[`item-${this.state.highlightedIndex}`];
      const menuNode = this.refs.menu;
      scrollIntoView(itemNode, menuNode, { onlyScrollIfNeeded: true });
    }
  }

  handleKeyDown = (event) => {
    const keyHandlers = this.getKeyDownHandlers();
    if (keyHandlers[event.key]) {
      keyHandlers[event.key].call(this, event);
    } else {
      const { selectionStart, value } = event.target;
      if (value === this.state.value) {
        // Nothing changed, no need to do anything. This also prevents
        // our workaround below from nuking user-made selections
        return;
      }
      this.setState({
        highlightedIndex: null,
        isOpen: true,
      }, () => {
        // Restore caret position before autocompletion process
        // to work around a setSelectionRange bug in IE (#80)
        this.refs.input.selectionStart = selectionStart;
      });
    }
  }

  handleChange = (event) => {
    this._performAutoCompleteOnKeyUp = true;
    this.props.onChange(event, event.target.value);
  }

  handleKeyUp = () => {
    if (this._performAutoCompleteOnKeyUp) {
      this._performAutoCompleteOnKeyUp = false;
      this.maybeAutoCompleteText();
    }
  }

  renderMenu = () => {
    const items = this.getFilteredItems().map((item, index) => {
      const element = this.props.renderItem(
        item,
        this.state.highlightedIndex === index,
        { cursor: 'default' },
      );
      return React.cloneElement(element, {
        onMouseDown: () => this.setIgnoreBlur(true),
        onMouseEnter: () => this.highlightItemFromMouse(index),
        onClick: () => this.selectItemFromMouse(item),
        ref: `item-${index}`,
      });
    });
    const style = {
      left: this.state.menuLeft,
      top: this.state.menuTop,
      minWidth: this.state.menuWidth,
    };
    const menu = this._renderMenu(items, this.props.value, style);
    return React.cloneElement(menu, { ref: 'menu' });
  }

  handleInputBlur = () => {
    if (this._ignoreBlur) {
      return;
    }
    this.setState({
      isOpen: false,
      highlightedIndex: null,
    });
  }

  handleInputFocus = () => {
    if (this._ignoreBlur) {
      return;
    }
    this.setState({ isOpen: true });
  }

  handleInputClick = () => {
    if (!this.state.isOpen) {
      this.setState({ isOpen: true });
    }
  }

  render() {
    const {
      classNames,
      inputProps,
      labelText,
      placeholder,
      value,
      wrapperProps,
    } = this.props;
    const wrapperClasses = classnames('autocompleteSearchBar', {
      [css.wrapper]: true,
      [classNames.wrapper]: !!classNames.wrapper,
    });
    const inputClasses = classnames('autocompleteInput', {
      [css.input]: true,
      [classNames.input]: !!classNames.input,
    });
    const iconClasses = classnames('autocompleteIcon', {
      [css.icon]: true,
      [classNames.icon]: !!classNames.icon,
    });
    return (
      <div className={wrapperClasses} {...wrapperProps} >
        <label htmlFor={this.id} ref="label">
          {labelText}
        </label>
        <SearchIcon className={iconClasses} size={SMALL} />
        <input
          {...inputProps}
          aria-autocomplete="both"
          className={inputClasses}
          id={this.id}
          onBlur={this.handleInputBlur}
          onChange={(event) => this.handleChange(event)}
          onClick={this.handleInputClick}
          onFocus={this.handleInputFocus}
          onKeyDown={(event) => this.handleKeyDown(event)}
          onKeyUp={(event) => this.handleKeyUp(event)}
          placeholder={placeholder}
          ref="input"
          role="combobox"
          value={value}
        />
        {this.state.isOpen && this.renderMenu()}
      </div>
    );
  }
}
/* eslint-enable react/sort-comp */
