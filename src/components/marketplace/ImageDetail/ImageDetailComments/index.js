import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { repositoryFetchComments } from 'actions/repository';
import BackButtonArea from 'common/BackButtonArea';
import CommentList from './CommentList';
import { DEFAULT_COMMENTS_PAGE_SIZE } from 'lib/constants/defaults';
import { getCurrentPage } from 'lib/utils/pagination';
import { Pagination } from 'common';
import ceil from 'lodash/ceil';
import get from 'lodash/get';
import makeRepositoryId from 'lib/utils/repo-image-name';
import routes from 'lib/constants/routes';
const { arrayOf, bool, func, number, shape, string } = PropTypes;

const mapState = ({ marketplace }, { location, params }) => {
  const { id, namespace, reponame } = params;
  const { community, certified } = marketplace.images;
  const page = getCurrentPage(location);
  const isCertified = !!id;
  let image;
  if (isCertified) {
    image = certified && certified[id];
    const { comments = {} } = image;
    const { isFetching, count } = comments;
    const results = get(comments, ['pages', page, 'results'], []);
    return {
      comments: results,
      count,
      name: image.name,
      isFetching,
      namespace: image.namespace,
      reponame: image.reponame,
    };
  }
  const repositoryId = makeRepositoryId({ namespace, reponame });
  // Community images don't have logos
  image = community && community[repositoryId];
  const { comments = {} } = image;
  const { isFetching, count } = comments;
  const results = get(comments, ['pages', page, 'results'], []);
  return {
    comments: results,
    count,
    isFetching,
    namespace: image.namespace,
    reponame: image.reponame,
  };
};

const dispatcher = {
  repositoryFetchComments,
};

@connect(mapState, dispatcher)
export default class ImageDetailComments extends Component {
  static propTypes = {
    comments: arrayOf(shape({
      id: number,
      user: string,
      comment: string,
      created_on: string,
      updated_on: string,
    })),
    count: number,
    name: string,
    isFetching: bool,
    namespace: string,
    reponame: string,
    repositoryFetchComments: func,
    params: shape({
      namespace: string,
      reponame: string,
    }).isRequired,
    location: shape({
      pathname: string,
    }).isRequired,
  }

  static contextTypes = {
    router: shape({
      push: func.isRequired,
    }).isRequired,
  }

  goToPage = (page) => {
    const { id } = this.props.params;
    const { namespace, reponame } = this.props;
    const { pathname, query, state } = this.props.location;
    const newQuery = { ...query, page };
    this.context.router.push({ pathname, query: newQuery, state });
    // Manually dispatch the action because a query change will not
    // trigger the `onEnter` event for the search route
    this.props.repositoryFetchComments({
      id,
      isCertified: !!id,
      namespace,
      reponame,
      ...newQuery,
    });
  };


  renderPagination() {
    const { query } = this.props.location;
    const currentPage = parseInt(query.page, 10) || 1;
    const { count, isFetching } = this.props;
    if (isFetching || !count) {
      return null;
    }
    const page_size =
      parseInt(query.page_size, 10) || DEFAULT_COMMENTS_PAGE_SIZE;
    const lastPage = ceil(count / page_size);
    return (
      <Pagination
        currentPage={currentPage}
        lastPage={lastPage}
        onChangePage={this.goToPage}
      />
    );
  }

  render() {
    const { id } = this.props.params;
    const { namespace, reponame } = this.props;
    let path;
    if (!!id) {
      path = routes.imageDetail({ id });
    } else {
      path = routes.communityImageDetail({ namespace, reponame });
    }
    return (
      <div>
        <BackButtonArea
          pathname={path}
          text="Back"
        />
        <CommentList {...this.props} />
        {this.renderPagination()}
      </div>
    );
  }
}
