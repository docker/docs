import React from 'react';
import { mount } from 'enzyme';
import Slider from './index.js';
import chai, { expect } from 'chai';
import chaiEnzyme from 'chai-enzyme';

chai.use(chaiEnzyme());

describe('<Slider />', () => {
  it('render Slider with Handle, Track, Mark, Step', () => {
    const wrapper = mount(<Slider min={32} max={1024} step={16} />);
    const html = wrapper.html();
    const slider = wrapper.find('.rc-slider');
    expect(slider).to.have.exactly(5).descendants('div');
    expect(html).to.contain('rc-slider');
    expect(html).to.contain('rc-slider-handle');
    expect(html).to.contain('rc-slider-track');
    expect(html).to.contain('rc-slider-mark');
    expect(html).to.contain('rc-slider-step');
  });
});
