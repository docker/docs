import React from 'react';
import { shallow } from 'enzyme';
import Uptime from './index.js';
import chai, { expect } from 'chai';
import chaiEnzyme from 'chai-enzyme';

chai.use(chaiEnzyme());

describe('<Uptime />', () => {
  it('renders uptime component correctly', () => {
    const date = new Date();
    date.setMinutes(date.getMinutes() - 50);
    const wrapper = shallow(<Uptime since={+date} />);

    // time icon
    expect(wrapper.find('.duptime')).to.have.descendants('span');
    const html = wrapper.find('.duptime').html();

    expect(html).contain('an hour ago');
    expect(html).contain('styles__display');
    expect(html).contain('dicon');
    expect(html).contain('styles__tiny');
  });

  // TODO: Write a test for interval property & negative tests
});
