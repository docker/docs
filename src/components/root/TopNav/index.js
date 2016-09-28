import React, { Component, PropTypes } from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';
import omit from 'lodash/omit';
import get from 'lodash/get';
import isEmpty from 'lodash/isEmpty';
import debounce from 'lodash/debounce';
import {
  AutocompleteSearchBar,
  Avatar,
  DockerFlatIcon,
  DockerStoreBetaIcon,
  DropdownIcon,
  MagicCarpet,
  Menu,
} from 'common';
import {
  marketplaceFetchAutocompleteSuggestions,
  marketplaceSearch,
} from 'actions/marketplace';
import { accountLogout, accountToggleMagicCarpet } from 'actions/account';
import { rootChangeGlobalSearchValue } from 'actions/root';
import Subscriptions from 'account/Subscriptions';
import BillingProfile from 'account/BillingProfile';
import { LARGE } from 'lib/constants/sizes';
import { PRIMARY } from 'lib/constants/variants';
import { SUBSCRIPTIONS, BILLING } from 'lib/constants/overlays';
import formatCategories from 'lib/utils/format-categories';
import { getLoginRedirectURL } from 'lib/utils/url-utils';
import routes from 'lib/constants/routes';
import { DEFAULT_DEBOUNCE_TIME } from 'lib/constants/defaults';
import css from './styles.css';
const { array, bool, func, shape, string, object } = PropTypes;

const ADMIN = 'admin';

const mapStateToProps = ({ account, marketplace, root, publish }) => {
  const { filters } = marketplace;
  const { currentUser, magicCarpet, isCurrentUserWhitelisted } = account;
  const publisherSignup = get(publish, ['signup', 'results']);
  const publisherVendor = get(publish, ['publishers', 'results']);
  const approvedPublisher =
    (
      publisherSignup.status === 'reviewed' ||
      !isEmpty(publisherVendor)
    );
  return {
    approvedPublisher,
    autocomplete: root && root.autocomplete || {},
    categories: filters && filters.categories || {},
    currentUser: currentUser || {},
    isCurrentUserWhitelisted,
    globalSearch: root && root.globalSearch || '',
    magicCarpet,
  };
};

const dispatcher = {
  logout: accountLogout,
  accountToggleMagicCarpet,
  marketplaceFetchAutocompleteSuggestions,
  marketplaceSearch,
  rootChangeGlobalSearchValue,
};

@connect(mapStateToProps, dispatcher)
export default class TopNav extends Component {
  static propTypes = {
    approvedPublisher: bool,
    autocomplete: shape({
      isFetching: bool,
      suggestions: array,
    }),
    categories: object,
    currentUser: object,
    isCurrentUserWhitelisted: bool,
    globalSearch: string,
    location: shape({
      query: object,
      pathname: string,
    }).isRequired,
    logout: func.isRequired,
    magicCarpet: string,
    accountToggleMagicCarpet: func.isRequired,
    marketplaceFetchAutocompleteSuggestions: func.isRequired,
    marketplaceSearch: func.isRequired,
    rootChangeGlobalSearchValue: func.isRequired,
    showTransparentNavBar: bool,
  }

  static contextTypes = {
    router: shape({
      push: func.isRequired,
    }).isRequired,
  }

  constructor(props) {
    super(props);
    // Make sure there is one debounce function per component instance
    // http://stackoverflow.com/a/28046731/5965502
    this.debouncedFetchAutocompleteSuggestions = debounce(
      this.debouncedFetchAutocompleteSuggestions,
      DEFAULT_DEBOUNCE_TIME
    );
  }

  componentWillUnmount() {
    this.debouncedFetchAutocompleteSuggestions.cancel();
  }

  onSelect = (value) => {
    switch (value) {
      case 'logout':
        this.logout();
        break;
      case 'publish':
        this.context.router.push({ pathname: routes.publisher() });
        break;
      case SUBSCRIPTIONS:
        this.openMagicCarpet(SUBSCRIPTIONS);
        break;
      case BILLING:
        this.openMagicCarpet(BILLING);
        break;
      case ADMIN:
        this.context.router.push({ pathname: routes.admin() });
        break;
      default:
        return;
    }
  }

  onSearchQueryChange = (e, value) => {
    // Remove this synthetic event from the pool so that we can still access the
    // event asyncronously (for debouncing)
    // https://facebook.github.io/react/docs/events.html#event-pooling
    e.persist();
    // change the value in the search bar
    this.props.rootChangeGlobalSearchValue({ value });
    // fetch new suggestions (debounced)
    this.debouncedFetchAutocompleteSuggestions(e, value);
  }

  onSelectAutosuggestItem = (value, item) => {
    const { id } = item;
    // Jump to the product detail page for this result
    const detail = routes.imageDetail({ id });
    this.context.router.push(detail);
  }

  debouncedFetchAutocompleteSuggestions = (e, value) => {
    this.props.marketplaceFetchAutocompleteSuggestions({ q: value });
  }

  // Global Search Bar form has been submitted
  search = (q) => {
    // fire search action and transition to search results page
    this.props.marketplaceSearch({ q });
    const pathname = routes.search();
    const { state } = this.props.location;
    // search from global search bar will have a query (q) (no page num or size)
    const query = { q };
    this.context.router.push({ pathname, query, state });
  }

  goToLogIn = () => {
    // After they login, redirect them back to whatever page they are on
    this.context.router.push(getLoginRedirectURL());
  }

  openMagicCarpet(carpet) {
    this.props.accountToggleMagicCarpet({ magicCarpet: carpet });
    const { query, pathname } = this.props.location;
    this.context.router.push({
      pathname,
      query: { ...query, overlay: carpet },
    });
  }

  closeMagicCarpet = () => {
    this.props.accountToggleMagicCarpet({ magicCarpet: '' });
    const { query, pathname } = this.props.location;
    this.context.router.push({
      pathname,
      query: omit(query, 'overlay'),
    });
  }

  logout = () => {
    // TODO Arunan 03/22/2016 Figure out if we need to clear the entire state
    //  Francisco:
    //   re above: We are not persisting any state across page refreshes
    //   so it doesn't matter really. We can handle this in the (newly created)
    //   LOGOUT_{REQ|ACK|ERR} action reducers when necessary
    this.props.logout().then(() => {
      // redirect home
      analytics.reset();
      window.location = routes.beta();
    });
  }

  renderLoggedInMenu() {
    const { gravatar_url, username } = this.props.currentUser;
    const avatarTrigger = (
      <div className={css.loggedIn}>
        <Avatar
          className={css.avatar}
          src={gravatar_url}
          style={{ width: '40px' }}
        />
        <strong className={css.username}>{username}</strong>
        <DropdownIcon />
      </div>
    );

    const menuOptions = [
      { label: 'Subscriptions', value: SUBSCRIPTIONS },
      { label: 'Billing', value: BILLING },
      { label: (<div className={css.logout}>Log out</div>), value: 'logout' },
    ];

    if (this.props.approvedPublisher) {
      menuOptions.unshift({
        label: 'Publisher Center',
        value: 'publish',
      });
    }

    return (
      <Menu
        className={css.accountDropdown}
        trigger={avatarTrigger}
        onSelect={this.onSelect}
        items={menuOptions}
      />
    );
  }

  renderMagicCarpet = () => {
    let content;
    let title;
    const { magicCarpet } = this.props;
    if (magicCarpet === SUBSCRIPTIONS) {
      title = 'Subscriptions';
      content = (<Subscriptions />);
    } else if (magicCarpet === BILLING) {
      title = 'Billing';
      content = (<BillingProfile />);
    }
    return (
      <MagicCarpet
        isOpen={!!magicCarpet}
        onRequestClose={this.closeMagicCarpet}
        title={title}
        titleIcon={<DockerFlatIcon size={LARGE} variant={PRIMARY} />}
      >
        {content}
      </MagicCarpet>
    );
  }

  renderAutocompleteItem = (item, isHighlighted) => {
    const { id, name, categories } = item;
    const catNames = formatCategories(categories);
    const catText = catNames ? ` in ${catNames}` : '';
    const itemClass = isHighlighted ? css.highlightedResult : css.result;
    return (
      <div className={itemClass} key={id} id={id}>
        <span className={css.resultName}>{name}</span>
        <span className={css.resultCategories}>{catText}</span>
      </div>
    );
  };

  render() {
    const { currentUser } = this.props;
    const isLoggedIn = currentUser && currentUser.id || false;
    let menu = (
      <div className={css.buttonLink} onClick={this.goToLogIn}>
        Log In
      </div>
    );
    if (isLoggedIn) {
      menu = this.renderLoggedInMenu();
    }
    const {
      approvedPublisher,
      autocomplete,
      globalSearch,
      isCurrentUserWhitelisted,
      showTransparentNavBar,
    } = this.props;
    let searchBar;
    if (!showTransparentNavBar && isCurrentUserWhitelisted) {
      const { suggestions = [] } = autocomplete;
      const menuTitle = (
        <div className={css.menuTitle}>
          Suggested Results
        </div>
      );
      const classNames = { input: css.input, menu: css.menu };
      const getItemValue = (item) => item.id;
      searchBar = (
        <AutocompleteSearchBar
          classNames={classNames}
          getItemValue={getItemValue}
          id="global-search"
          items={suggestions}
          menuTitle={menuTitle}
          onChange={this.onSearchQueryChange}
          onSelect={this.onSelectAutosuggestItem}
          onSubmit={this.search}
          ref="autocomplete"
          renderItem={this.renderAutocompleteItem}
          value={globalSearch}
        />
      );
    }
    const classes = showTransparentNavBar ?
      `${css.topnav} ${css.transparent}` : css.topnav;
    // TODO Kristie 6/14/16 Switch beta to home when beta is over
    let homeLink = routes.beta();
    let browse;
    if (isLoggedIn && isCurrentUserWhitelisted) {
      homeLink = routes.home();
      browse = (
        <Link className={css.buttonLink} to={routes.browse()}>
          Browse
        </Link>
      );
    }
    let publisher;
    if (approvedPublisher) {
      publisher = (
        <Link className={css.buttonLink} to={routes.publisherAddProduct()}>
          Publish a Product
        </Link>
      );
    } else {
      publisher = (
        <Link className={css.buttonLink} to={routes.publisherSignup()}>
          Become a Publisher
        </Link>
      );
    }

    const feedback = (
      <a
        className={css.buttonLink}
        href="mailto:store-feedback@docker.com?Subject=Store%20Beta%20Feedback"
      >Feedback?</a>
    );
    return (
      <header className={classes}>
        <div className={css.wrappedColumns}>
          <Link to={homeLink} className={css.buttonLink}>
            <DockerStoreBetaIcon />
          </Link>
          <span className={css.search}>{searchBar}</span>
          <div className={css.accountLinks}>
            {browse}
            {publisher}
            {feedback}
            {menu}
          </div>
        </div>
        {this.renderMagicCarpet()}
      </header>
    );
  }
}
