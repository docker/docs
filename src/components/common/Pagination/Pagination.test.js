import React from 'react';
import { mount, shallow } from 'enzyme';
import Pagination from './index.js';
import chai, { expect } from 'chai';
import chaiEnzyme from 'chai-enzyme';

chai.use(chaiEnzyme());
// The above imports and statement need to move somewhere global

describe('<Pagination />', () => {
  it('renders the pagination element', () => {
    const onChangePage = () => {};
    const wrapper = shallow(
      <Pagination
        currentPage={1}
        lastPage={10}
        maxVisible={10}
        onChangePage={onChangePage}
      />
    );
    expect(wrapper.find('.dpagination')).to.be.present();
    expect(wrapper.find('.dpagination')).to.have.exactly(1).descendants('ul');
    expect(wrapper.html()).to.contain('styles__paginationCentered');
  });

  it('hides the previous/last buttons when the current page is 1', () => {
    const onChangePage = () => {};
    const wrapper = mount(
      <Pagination
        currentPage={1}
        lastPage={10}
        maxVisible={10}
        onChangePage={onChangePage}
      />
    );
    expect(wrapper.find('.dpagination')).to.be.present();
    expect(wrapper.find('.dpagination')).to.have.exactly(2).descendants('svg');
    expect(wrapper.html()).to.contain('styles__nextPage');
    expect(wrapper.html()).to.not.contain('styles__previousPage');
    expect(wrapper.html()).to.not.contain('styles__firstPage');
  });

  it('hides the next/last buttons when current page = last page', () => {
    const onChangePage = () => {};
    const wrapper = mount(
      <Pagination
        currentPage={10}
        lastPage={10}
        maxVisible={10}
        onChangePage={onChangePage}
      />
    );
    expect(wrapper.find('.dpagination')).to.be.present();
    expect(wrapper.find('.dpagination')).to.have.exactly(2).descendants('svg');
    expect(wrapper.html()).to.not.contain('styles__nextPage');
    expect(wrapper.html()).to.contain('styles__previousPage');
    expect(wrapper.html()).to.contain('styles__firstPage');
  });

  it('shows all chevron buttons when applicable', () => {
    const onChangePage = () => {};
    const wrapper = mount(
      <Pagination
        currentPage={2}
        lastPage={4}
        maxVisible={10}
        onChangePage={onChangePage}
      />
    );
    expect(wrapper.find('.dpagination')).to.be.present();
    expect(wrapper.find('.dpagination')).to.have.exactly(4).descendants('svg');
    expect(wrapper.html()).to.contain('styles__nextPage');
    expect(wrapper.html()).to.contain('styles__previousPage');
    expect(wrapper.html()).to.contain('styles__firstPage');
  });

  it('only shows the maxVisible number of pages if lastPage is greater', () => {
    const onChangePage = () => {};
    const maxVisible = 4;
    const wrapper = mount(
      <Pagination
        currentPage={2}
        lastPage={10}
        maxVisible={maxVisible}
        onChangePage={onChangePage}
      />
    );
    expect(wrapper.find('.dpagination')).to.be.present();
    expect(wrapper.find('.dpagination'))
      .to.have.exactly(maxVisible).descendants('.dpage');
  });

  it('renders just one page if there is only one page', () => {
    const onChangePage = () => {};
    const wrapper = mount(
      <Pagination
        currentPage={1}
        lastPage={1}
        maxVisible={10}
        onChangePage={onChangePage}
      />
    );
    expect(wrapper.find('.dpagination')).to.be.present();
    expect(wrapper.find('.dpage').first()).to.have.text(1);
    expect(wrapper.find('.dpagination')).to.be.present();
    expect(wrapper.find('.dpagination'))
      .to.have.exactly(1).descendants('.dpage');
    expect(wrapper.html()).to.not.contain('styles__nextPage');
    expect(wrapper.html()).to.not.contain('styles__previousPage');
    expect(wrapper.html()).to.not.contain('styles__firstPage');
  });
});
