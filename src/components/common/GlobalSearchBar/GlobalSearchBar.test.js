import React from 'react';
import { mount, shallow } from 'enzyme';
import GlobalSearchBar from './index.js';
import chai, { expect } from 'chai';
import chaiEnzyme from 'chai-enzyme';
import sinon from 'sinon';

chai.use(chaiEnzyme());

const onChange = () => {};
const onSubmit = () => {};

describe('<GlobalSearchBar />', () => {
  it('renders the search bar element', () => {
    const wrapper = mount(
      <GlobalSearchBar onChange={onChange} onSubmit={onSubmit} />
    );
    expect(wrapper.find('form')).to.have.length(1);
    expect(wrapper.find('input')).to.have.length(1);
    expect(wrapper.html()).to.contain('dGlobalSearchBar');
  });

  it('renders the search icon inside the search bar', () => {
    const wrapper = mount(
      <GlobalSearchBar onChange={onChange} onSubmit={onSubmit} />
    );
    expect(wrapper.find('svg')).to.have.length(1);
    expect(wrapper.find('svg')).to.include.className('dicon');
    expect(wrapper.html()).to.contain('styles__icon');
    expect(wrapper.html()).to.contain('styles__small');
  });

  it('receives the value prop', () => {
    const wrapper = mount(
      <GlobalSearchBar onChange={onChange} onSubmit={onSubmit} value="test" />
    );
    // it is a controlled input so the value is not visible in the HTML
    expect(wrapper.prop('value')).to.equal('test');
  });

  it('calls the onChange method', () => {
    const onChangeSpy = sinon.spy();
    const wrapper = mount(
      <GlobalSearchBar
        onChange={onChangeSpy}
        onSubmit={onSubmit}
      />
    );
    wrapper.find('input').simulate('change', { target: { value: 'test' } });
    expect(onChangeSpy.calledOnce).to.equal(true);
    // There is a parent component that must control the state that we can't
    // simulate at this level so we can't check the value
  });

  it('calls the onSubmit method', () => {
    const onSubmitSpy = sinon.spy();
    const wrapper = mount(
      <GlobalSearchBar
        onChange={onChange}
        onSubmit={onSubmitSpy}
      />
    );
    wrapper.find('form').simulate('submit');
    expect(onSubmitSpy.calledOnce).to.equal(true);
  });

  it('renders a placeholder', () => {
    const wrapper = shallow(
      <GlobalSearchBar
        onChange={onChange}
        onSubmit={onSubmit}
        placeholder="test"
      />
    );
    expect(wrapper.find('input')).to.have.length(1);
    expect(wrapper.html()).to.contain('placeholder');
  });
});
