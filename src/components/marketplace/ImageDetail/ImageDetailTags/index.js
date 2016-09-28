import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import BackButtonArea from 'common/BackButtonArea';
import { Card, Pagination } from 'common';
import RepositoryTagsList from './RepositoryTagsList';
import ImageHeaderWithLogo from '../ImageHeaderWithLogo';
import ceil from 'lodash/ceil';
import get from 'lodash/get';
import capitalize from 'lodash/capitalize';
import { getCurrentPage } from 'lib/utils/pagination';
import getLogo from 'lib/utils/get-largest-logo';
import { DEFAULT_TAGS_PAGE_SIZE } from 'lib/constants/defaults';
import { repositoryFetchImageTags } from 'actions/repository';
import makeRepositoryId from 'lib/utils/repo-image-name';
import routes from 'lib/constants/routes';
const { array, bool, func, number, object, shape, string } = PropTypes;

const mapStateToProps = ({ marketplace }, { location, params }) => {
  const { id, namespace, reponame } = params;
  const page = getCurrentPage(location);
  const { certified, community } = marketplace && marketplace.images;
  const isCertified = !!id;
  let image;
  if (isCertified) {
    image = certified && certified[id];
    const { tags = {}, logo_url } = image;
    const { isFetching, count } = tags;
    const results = get(tags, ['pages', page, 'results'], []);
    return {
      name: image.name,
      namespace: image.namespace,
      reponame: image.reponame,
      tags: {
        isFetching,
        count,
        results,
      },
      logo: getLogo(logo_url),
    };
  }
  const repositoryId = makeRepositoryId({ namespace, reponame });
  // Community images don't have logos
  image = community && community[repositoryId];
  const { tags = {} } = image;
  const { isFetching, count } = tags;
  const results = get(tags, ['pages', page, 'results'], []);
  return {
    namespace,
    reponame,
    tags: {
      isFetching,
      count,
      results,
    },
  };
};

const dispatcher = {
  repositoryFetchImageTags,
};

@connect(mapStateToProps, dispatcher)
export default class ImageDetailTags extends Component {
  static propTypes = {
    repositoryFetchImageTags: func.isRequired,
    tags: shape({
      isFetching: bool,
      count: number,
      results: array,
    }),
    location: shape({
      pathname: string,
    }).isRequired,
    logo: string,
    namespace: string,
    reponame: string,
    name: string,
    params: object.isRequired,
  }

  static contextTypes = {
    router: shape({
      push: func.isRequired,
    }).isRequired,
  }

  goToPage = (page) => {
    const { id } = this.props.params;
    const { reponame, namespace } = this.props;
    const { pathname, query, state } = this.props.location;
    const newQuery = { ...query, page };
    this.context.router.push({ pathname, query: newQuery, state });
    // Manually dispatch the action because a query change will not
    // trigger the `onEnter` event for the search route
    this.props.repositoryFetchImageTags({
      id,
      isCertified: !!id,
      namespace,
      reponame,
      ...newQuery,
    });
  };

  mkTagTable() {
    const { logo, reponame, name, tags } = this.props;
    const { count, results } = tags;
    if (!count) {
      return null;
    }
    const displayName = name || reponame;
    const tagsOrVersions = name ? 'Versions' : 'Tags';
    const title = `${capitalize(displayName)} ${tagsOrVersions}`;
    const helpText = `${count} ${tagsOrVersions} available`;
    return (
      <div>
        <ImageHeaderWithLogo
          title={title}
          logo={logo}
          helpText={helpText}
        />
        <RepositoryTagsList tags={results} />
      </div>
    );
  }

  renderPagination() {
    const { query } = this.props.location;
    const currentPage = parseInt(query.page, 10) || 1;
    const { count, isFetching } = this.props.tags;
    if (isFetching || !count) {
      return null;
    }
    const page_size =
      parseInt(query.page_size, 10) || DEFAULT_TAGS_PAGE_SIZE;
    const lastPage = ceil(count / page_size);
    return (
      <Pagination
        currentPage={currentPage}
        lastPage={lastPage}
        onChangePage={this.goToPage}
      />
    );
  }

  renderBlankSlate() {
    // TODO Kristie 3/25/16 Get designs for actual blank slate
    return (
      <Card>
        <div>Empty repository.</div>
      </Card>
    );
  }

  renderIsFetching() {
    // TODO Kristie 3/25/16 Get designs for actual fetching
    return (
      <div>Fetching...</div>
    );
  }

  render() {
    const { id } = this.props.params;
    const { namespace, reponame } = this.props;
    const { isFetching, count } = this.props.tags;
    let path;
    if (!!id) {
      path = routes.imageDetail({ id });
    } else {
      path = routes.communityImageDetail({ namespace, reponame });
    }
    let pagination;
    let pageContent;
    if (isFetching) {
      pageContent = this.renderIsFetching();
    } else if (!count) {
      pageContent = this.renderBlankSlate();
    } else {
      pagination = this.renderPagination();
      pageContent = this.mkTagTable();
    }
    return (
      <div>
        <BackButtonArea pathname={path} text="Back" />
        {pageContent}
        {pagination}
      </div>
    );
  }
}
