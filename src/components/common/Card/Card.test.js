import React from 'react';
import { shallow } from 'enzyme';
import Card from './index.js';
import chai, { expect } from 'chai';
import chaiEnzyme from 'chai-enzyme';

chai.use(chaiEnzyme());

describe('<Card />', () => {
  it('renders a default card', () => {
    const wrapper = shallow(<Card children={[]} />);
    expect(wrapper.find('.dcard')).to.be.present();
  });

  it('renders a card with title', () => {
    const wrapper = shallow(
      <Card children={[]} title="Test">
        <p>test step 1</p>
        <p>test step 2</p>
        <p>test step 3</p>
      </Card>
    );
    expect(wrapper.find('.dcard')).to.have.exactly(1).descendants('h3');
    expect(wrapper.find('h3')).to.have.text('Test');
    expect(wrapper.find('.dcard')).to.have.exactly(3).descendants('p');
  });

  it('renders a card with a `shy` title', () => {
    const wrapper = shallow(
      <Card children={[]} shy title="Shy Title Test">
        <p>test step 1</p>
        <p>test step 2</p>
        <p>test step 3</p>
      </Card>
    );
    expect(wrapper.html()).to.contain('styles__shy');
    expect(wrapper.find('.dcard')).to.have.exactly(1).descendants('h3');
    expect(wrapper.find('h3')).to.have.text('Shy Title Test');
    expect(wrapper.find('.dcard')).to.have.exactly(3).descendants('p');
  });

  it('renders a card with no title', () => {
    const wrapper = shallow(
      <Card children={[]}>
        <p>test step 1</p>
        <p>test step 2</p>
        <p>test step 3</p>
      </Card>
    );
    expect(wrapper.find('.dcard')).to.not.contain('h3');
    expect(wrapper.find('.dcard')).to.have.exactly(3).descendants('p');
  });
});
