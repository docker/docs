import React from 'react';
import { mount } from 'enzyme';
import TagList from './index.js';
import chai, { expect } from 'chai';
import chaiEnzyme from 'chai-enzyme';

chai.use(chaiEnzyme());

describe('<TagList />', () => {
  it('renders the TagList element', () => {
    const wrapper = mount(
      <TagList list={[
        'tag1',
        'tag2',
        'tag3',
        'tag4',
      ]} size={380}
      />
    );

    expect(wrapper.find('.dTagList')).to.be.present();
    expect(wrapper.props().list).to.deep.equal(
      ['tag1', 'tag2', 'tag3', 'tag4']
    );
    expect(wrapper.find('ul').children()).to.have.exactly(4).descendants('Tag');
  });
});
