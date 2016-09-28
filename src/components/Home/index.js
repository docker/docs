import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import css from './styles.css';
import { categoryDescriptions } from 'lib/constants/landingPage';
import CategoryCards from './CategoryCards';
import DDCBanner from './DDCBanner';
import FeaturedContentRow from './FeaturedContentRow';
import HelpArticlesCards from './HelpArticlesCards';
import SearchWithAutocomplete from './SearchWithAutocomplete';
import routes from 'lib/constants/routes';
const { array, bool, func, object, shape } = PropTypes;

const mapStateToProps = ({ marketplace, root }) => {
  const { categories } = marketplace.filters;
  const { mostPopular, featured } = root.landingPage;
  return {
    categories,
    mostPopular,
    featured,
  };
};

@connect(mapStateToProps)
export default class Home extends Component {
  static propTypes = {
    autocomplete: shape({
      isFetching: bool,
      suggestions: array,
    }),
    categories: object.isRequired,
    location: object.isRequired,
    mostPopular: shape({
      isFetching: bool,
      images: array,
    }),
    featured: shape({
      isFetching: bool,
      images: array,
    }),
  }

  static contextTypes = {
    router: shape({
      push: func.isRequired,
    }).isRequired,
  }

  goToCategory = (category) => () => {
    const query = { category };
    const { state } = this.props.location;
    const pathname = routes.search();
    this.context.router.push({ pathname, query, state });
  }

  renderHero() {
    return (
      <div className={css.heroWrapper}>
        <div className="wrapped">
          <div className={css.heroContent}>
            <div className={css.title}>Search The Docker Store</div>
            <div className={css.subText}>
              Find Trusted and Enterprise Ready Containers
            </div>
            <SearchWithAutocomplete location={this.props.location} />
          </div>
        </div>
      </div>
    );
  }

  renderFeaturedRow1() {
    const { location, featured } = this.props;
    // TODO Kristie 7/18/16 Make this sound better :)
    const featuredHeadline = 'Featured Images';
    const featuredTitle = 'Image Spotlight';
    const featuredDescription =
      'Curated images, ready to use. Verified and secure.';
    return (
      <div className={`wrapped ${css.sectionWrapper}`}>
        <FeaturedContentRow
          description={featuredDescription}
          className={css.paddedRow}
          headline={featuredHeadline}
          images={featured.images}
          isFetching={featured.isFetching}
          location={location}
          title={featuredTitle}
        />
      </div>
    );
  }

  renderFeaturedRow2() {
    const { location, mostPopular } = this.props;
    const popularHeadline = 'Most Popular Containers on the Store';
    const popularTitle = 'Most Popular';
    const popularDescription =
      'Most popular content from our trusted partners on the Docker Store';
    return (
      <div className={`wrapped ${css.sectionWrapper}`}>
        <FeaturedContentRow
          description={popularDescription}
          className={css.paddedRow}
          headline={popularHeadline}
          images={mostPopular.images}
          isFetching={mostPopular.isFetching}
          location={location}
          title={popularTitle}
        />
      </div>
    );
  }


  renderCategoryCards() {
    const { categories } = this.props;
    return (
      <div className="wrapped">
        <CategoryCards
          categories={categories}
          categoryDescriptions={categoryDescriptions}
          goToCategory={this.goToCategory}
        />
      </div>
    );
  }

  render() {
    return (
      <div className={css.home}>
        {this.renderHero()}
        {this.renderFeaturedRow1()}
        <DDCBanner />
        {this.renderFeaturedRow2()}
        {this.renderCategoryCards()}
        <HelpArticlesCards />
      </div>
    );
  }
}
