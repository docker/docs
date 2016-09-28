import React from 'react';
import { mount } from 'enzyme';
import AutocompleteSearchBar from './index.js';
import chai, { expect } from 'chai';
import chaiEnzyme from 'chai-enzyme';
import sinon from 'sinon';

chai.use(chaiEnzyme());

const item1 = { id: 1, name: 'item1' };
const item2 = { id: 2, name: 'item2' };
const items = [item1, item2];
const onChange = () => {};
const onSubmit = () => {};
const onSelect = () => {};
const getItemValue = (item) => item.name;
const renderItem = (item, isHighlighted) => {
  const classname = isHighlighted ? 'highlightedResult' : 'result';
  return <div key={item.id} className={classname}>{item.name}</div>;
};
const id = 'autocomplete';


describe('<AutocompleteSearchBar />', () => {
  it('renders the search bar element and input', () => {
    const wrapper = mount(
      <AutocompleteSearchBar
        getItemValue={getItemValue}
        id={id}
        items={items}
        onChange={onChange}
        onSelect={onSelect}
        onSubmit={onSubmit}
        renderItem={renderItem}
        value="test"
      />
    );
    expect(wrapper.find('input')).to.have.length(1);
    expect(wrapper.html()).to.contain('autocompleteSearchBar');
    expect(wrapper.html()).to.contain('autocompleteInput');
  });

  it('renders the search icon inside the search bar', () => {
    const wrapper = mount(
      <AutocompleteSearchBar
        getItemValue={getItemValue}
        id={id}
        items={items}
        onChange={onChange}
        onSelect={onSelect}
        onSubmit={onSubmit}
        renderItem={renderItem}
        value="test"
      />
    );
    expect(wrapper.find('svg')).to.have.length(1);
    expect(wrapper.find('svg')).to.include.className('dicon');
    expect(wrapper.html()).to.contain('styles__icon');
    expect(wrapper.html()).to.contain('styles__small');
  });

  it('receives the value prop', () => {
    const wrapper = mount(
      <AutocompleteSearchBar
        getItemValue={getItemValue}
        id={id}
        items={items}
        onChange={onChange}
        onSelect={onSelect}
        onSubmit={onSubmit}
        renderItem={renderItem}
        value="test"
      />
    );
    // it is a controlled input so the value is not visible in the HTML
    expect(wrapper.prop('value')).to.equal('test');
  });

  it('calls onChange and renders the menu when the input is changed', () => {
    const onChangeSpy = sinon.spy();
    const onSelectSpy = sinon.spy();
    const wrapper = mount(
      <AutocompleteSearchBar
        getItemValue={getItemValue}
        id={id}
        items={items}
        onChange={onChangeSpy}
        onSelect={onSelectSpy}
        onSubmit={onSubmit}
        renderItem={renderItem}
        value="test"
      />
    );
    wrapper.find('input').simulate('focus');
    // this simulates typing in "test" to the input
    wrapper.find('input').simulate('change', { target: { value: 'item1' } });
    // Items should be showing
    expect(onChangeSpy.calledOnce).to.equal(true);
    expect(wrapper.html()).to.contain('styles__menu');
    expect(wrapper.html()).to.contain('result');
  });

  it('calls the onSubmit method', () => {
    const onSubmitSpy = sinon.spy();
    const wrapper = mount(
      <AutocompleteSearchBar
        getItemValue={getItemValue}
        id={id}
        items={items}
        onChange={onChange}
        onSelect={onSelect}
        onSubmit={onSubmitSpy}
        renderItem={renderItem}
        value="test"
      />
    );
    wrapper.find('input').simulate('focus');
    wrapper.find('input').simulate('keyDown', {
      key: 'Enter',
      keyCode: 13,
      type: 'keydown',
      which: 13,
    });
    expect(onSubmitSpy.calledOnce).to.equal(true);
  });

  it('calls the onSelect method', () => {
    const onSelectSpy = sinon.spy();
    const wrapper = mount(
      <AutocompleteSearchBar
        getItemValue={getItemValue}
        id={id}
        items={items}
        onChange={onChange}
        onSelect={onSelectSpy}
        onSubmit={onSubmit}
        renderItem={renderItem}
        value="test"
      />
    );
    wrapper.find('input').simulate('focus');
    // Simulate the down arrow keypress
    wrapper.find('input').simulate('keyDown', {
      key: 'ArrowDown',
      keyCode: 40,
      type: 'keydown',
      which: 40,
    });
    // should be a highlighted element
    expect(wrapper.html()).to.contain('highlightedResult');
    // simulate hitting enter
    wrapper.find('input').simulate('keyDown', {
      key: 'Enter',
      keyCode: 13,
      type: 'keydown',
      which: 13,
    });
    expect(onSelectSpy.calledOnce).to.equal(true);
  });
});
