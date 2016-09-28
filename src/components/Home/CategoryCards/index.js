import React, { Component, PropTypes } from 'react';
import css from './styles.css';
import {
  ApplicationServicesIcon,
  DatabasesIcon,
  MessagingServicesIcon,
  OperatingSystemsIcon,
  ProgrammingLanguagesIcon,
  StorageIcon,
  AnalyticsIcon,
  // ApplicationFrameworkIcon,
  InfrastructureIcon,
  BaseImagesIcon,
  FeaturedImagesIcon,
  ToolsIcon,
} from 'common';
import {
  ANALYTICS_CATEGORY,
  APPLICATION_FRAMEWORK_CATEGORY,
  APPLICATION_INFRASTRUCTURE_CATEGORY,
  APPLICATION_SERVICES_CATEGORY,
  BASE_CATEGORY,
  DATABASE_CATEGORY,
  FEATURED_CATEGORY,
  LANGUAGES_CATEGORY,
  MESSAGING_CATEGORY,
  OS_CATEGORY,
  STORAGE_CATEGORY,
  TOOLS_CATEGORY,
} from 'lib/constants/landingPage';
import ChevronArrows from '../ChevronArrows';
import findKey from 'lodash/findKey';
const { func, object } = PropTypes;
const noOp = () => {};

export default class CategoryCards extends Component {
  static propTypes = {
    categories: object.isRequired,
    categoryDescriptions: object.isRequired,
    goToCategory: func.isRequired,
  }

  state = {
    firstCategoryIndex: 0,
  }

  setFirstCategoryIndex = (i) => () => {
    this.setState({
      firstCategoryIndex: i,
    });
  }

  getIcon = (category) => {
    switch (category.name) {
      case ANALYTICS_CATEGORY:
        return <AnalyticsIcon className={css.icon} />;
      // TODO Get real icon for framework
      case APPLICATION_FRAMEWORK_CATEGORY:
        return <ApplicationServicesIcon className={css.icon} />;
      case APPLICATION_INFRASTRUCTURE_CATEGORY:
        return <InfrastructureIcon className={css.icon} />;
      case APPLICATION_SERVICES_CATEGORY:
        return <ApplicationServicesIcon className={css.icon} />;
      case BASE_CATEGORY:
        return <BaseImagesIcon className={css.icon} />;
      case DATABASE_CATEGORY:
        return <DatabasesIcon className={css.icon} />;
      case FEATURED_CATEGORY:
        return <FeaturedImagesIcon className={css.icon} />;
      case LANGUAGES_CATEGORY:
        return <ProgrammingLanguagesIcon className={css.icon} />;
      case MESSAGING_CATEGORY:
        return <MessagingServicesIcon className={css.icon} />;
      case OS_CATEGORY:
        return <OperatingSystemsIcon className={css.icon} />;
      case STORAGE_CATEGORY:
        return <StorageIcon className={css.icon} />;
      case TOOLS_CATEGORY:
        return <ToolsIcon className={css.icon} />;
      default:
        return <OperatingSystemsIcon className={css.icon} />;
    }
  }

  renderCard = (category, index) => {
    const { categories, goToCategory } = this.props;
    const { name, description } = category;
    const categoryName = findKey(categories, (label) => label === name);
    const backgroundClass = `card${index}`;
    return (
      <div
        key={index}
        className={`${css.card} ${css[backgroundClass]}`}
        onClick={goToCategory(categoryName)}
      >
        <div className={css.iconWrapper}>
          {this.getIcon(category)}
        </div>
        <div className={css.cardContent}>
          <div className={css.cardTitle}>{name}</div>
          <div className={css.cardDescription}>{description}</div>
        </div>
      </div>
    );
  }

  render() {
    const {
      analytics,
      application_framework,
      application_services,
      application_infrastructure,
      base,
      database,
      featured,
      languages,
      messaging,
      os,
      storage,
      tools,
    } = this.props.categoryDescriptions;
    const categories = [
      featured,
      database,
      base,
      languages,
      application_services,
      messaging,
      os,
      analytics,
      application_framework,
      storage,
      tools,
      application_infrastructure,
    ];
    const numCategories = categories.length;
    const { firstCategoryIndex } = this.state;
    const numShowing = 6;
    const lastCategoryIndex = firstCategoryIndex + numShowing - 1;
    const displayCategories =
      categories.slice(firstCategoryIndex, lastCategoryIndex + 1);
    const isPreviousDisabled = firstCategoryIndex === 0;
    const isNextDisabled = lastCategoryIndex === numCategories - 1;
    const onClickNext = isNextDisabled ?
      noOp : this.setFirstCategoryIndex(lastCategoryIndex + 1);
    const onClickPrevious = isPreviousDisabled ?
      noOp : this.setFirstCategoryIndex(firstCategoryIndex - numShowing);
    return (
      <div>
        <div className={css.sectionTitleWrapper}>
          <div className={css.sectionTitle}>
            Everything you need to build applications for your business
          </div>
          <ChevronArrows
            isNextDisabled={isNextDisabled}
            isPreviousDisabled={isPreviousDisabled}
            onClickNext={onClickNext}
            onClickPrevious={onClickPrevious}
          />
        </div>
        <div className={css.categoryCards}>
          {displayCategories.map(this.renderCard)}
        </div>
      </div>
    );
  }
}
