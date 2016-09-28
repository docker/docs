import React from 'react';
import { mount, shallow } from 'enzyme';
import Button from './index.js';
import { DockerFlatIcon } from '../Icon';
import chai, { expect } from 'chai';
import chaiEnzyme from 'chai-enzyme';

chai.use(chaiEnzyme());

describe('<Button />', () => {
  it('renders the button element', () => {
    const wrapper = shallow(<Button>Test</Button>);
    expect(wrapper.find('button')).to.have.length(1);
  });

  it('renders secondary outlined button', () => {
    const wrapper = shallow(<Button variant="secondary" outlined>Test</Button>);
    expect(wrapper.find('button')).to.have.text('Test');
    expect(wrapper.find('button')).to.include.className('dbutton');
    expect(wrapper.html()).to.contain('styles__button');
    expect(wrapper.html()).to.contain('styles__secondary');
    expect(wrapper.html()).to.contain('styles__outlined');
  });

  it('renders dull link button', () => {
    const wrapper = shallow(<Button variant="dull" text>Test</Button>);
    expect(wrapper.html()).to.contain('styles__button');
    expect(wrapper.html()).to.contain('styles__dull');
    expect(wrapper.html()).to.contain('styles__text');
  });

  it('renders panic disabled button', () => {
    const wrapper = mount(<Button variant="panic" disabled>Aaaaaaah</Button>);
    expect(wrapper.html()).to.contain('styles__button');
    expect(wrapper.html()).to.contain('styles__panic');
    expect(wrapper.find('button')).to.be.disabled();
  });

  it('renders a warning icon button', () => {
    const wrapper = mount(
      <Button
        variant="warn"
        icon
      >
        <DockerFlatIcon />
      </Button>
    );
    expect(wrapper.html()).to.contain('styles__button');
    expect(wrapper.html()).to.contain('styles__warn');
    expect(wrapper.find('button')).to.have.exactly(1).descendants('svg');
    expect(wrapper.find('button')).to.have.prop('icon', true);
    expect(wrapper.find('button')).to.not.be.disabled();
  });

  it('renders a primary button with an icon on the left', () => {
    const wrapper = mount(
      <Button
        icon="left"
      >
        <DockerFlatIcon />
        Dockerhub
      </Button>
    );
    expect(wrapper.html()).to.contain('styles__button');
    expect(wrapper.html()).to.contain('styles__primary');
    expect(wrapper.find('button')).to.have.exactly(1).descendants('svg');
    expect(wrapper.find('button')).to.have.prop('icon', 'left');
    expect(wrapper.find('button')).to.have.text('Dockerhub');
    expect(wrapper.find('button')).to.not.be.disabled();
  });
});
