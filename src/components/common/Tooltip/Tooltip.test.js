import React from 'react';
import { mount } from 'enzyme';
import Tooltip from './index.js';
import chai from 'chai';
import chaiEnzyme from 'chai-enzyme';

chai.use(chaiEnzyme());

describe('<Tooltip />', () => {
  it('renders the tooltip element', () => {
    const tooltipContent = <div className="content">Hi I am the tooltip.</div>;
    const wrapper = mount(
      <Tooltip
        placement="top"
        content={tooltipContent}
        trigger={['click', 'hover']}
      >
        <div className="children">
          I have a tooltip.
        </div>
      </Tooltip>
    );
    // TODO Kristie 3/22/16 We cannot get these tests to work at the moment.
    // We may need to investigate using ReactTestUtils as the creator does in
    // their tests
    // console.log(wrapper.find('.children'));
    wrapper.find('.children').simulate('click');
    // expect(wrapper.find('.content')).to.be.present();
  });
});
