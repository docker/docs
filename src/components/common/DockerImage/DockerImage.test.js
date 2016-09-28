import React from 'react';
import { shallow } from 'enzyme';
import DockerImage from './index.js';
import chai, { expect } from 'chai';
import chaiEnzyme from 'chai-enzyme';

chai.use(chaiEnzyme());

describe('<DockerImage />', () => {
  it('renders a Docker image component (<icon> <repository name>)', () => {
    const wrapper = shallow(<DockerImage imageName="debian" />);

    expect(wrapper.find('.ddockerImage').html()).to.contain('styles__regular');
    expect(wrapper.find('.ddockerImage'))
      .to.have.exactly(2).descendants('span');
    // TODO Arunan 03/14/2016 Figure out a less messier way to do this
    expect(wrapper.find('span').nodes[0].props.className)
      .to.contain('styles__icon');
    expect(wrapper.find('span').nodes[1].props.className)
      .to.contain('styles__name');
    expect(wrapper.html()).to.contain('debian');
  });

  it('renders a standalone Docker image component', () => {
    const wrapper = shallow(
      <DockerImage standalone size="xlarge" imageName="debian" />
    );

    expect(wrapper.find('.ddockerImage').html())
      .to.contain('styles__standalone');
    expect(wrapper.find('.ddockerImage').html())
      .to.contain('styles__xlarge');
    expect(wrapper.find('.ddockerImage'))
      .to.have.exactly(1).descendants('span');
    expect(wrapper.find('span').nodes[0].props.className)
      .to.contain('styles__icon');
    expect(wrapper.html()).to.not.contain('debian');
  });

  it('renders a Docker Image with text', () => {
    const wrapper = shallow(
      <DockerImage size="large" imageName="library/ubuntu:trusty" />
    );

    expect(wrapper.find('.ddockerImage').html())
      .to.contain('styles__large');
    expect(wrapper.find('.ddockerImage'))
      .to.have.exactly(2).descendants('span');
    expect(wrapper.find('span').nodes[0].props.className)
      .to.contain('styles__icon');
    expect(wrapper.find('span').nodes[1].props.className)
      .to.contain('styles__name');
    expect(wrapper.html()).not.to.contain('library');
    expect(wrapper.html()).to.contain('ubuntu:trusty');
  });
});
