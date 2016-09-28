import React from 'react';
import { mount } from 'enzyme';
import Select from './index.js';
import chai, { expect } from 'chai';
import chaiEnzyme from 'chai-enzyme';

chai.use(chaiEnzyme());

describe('<Select />', () => {
  it('renders a Select element', () => {
    const wrapper = mount(<Select />);

    expect(wrapper.find('.Select')).to.have.className('is-searchable');
    expect(wrapper.find('.Select-input'))
      .to.have.exactly(1).descendants('input');
  });
});
