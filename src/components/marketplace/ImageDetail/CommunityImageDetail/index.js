import React, { Component, PropTypes } from 'react';
import css from './styles.css';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import {
  BackButtonArea,
  Card,
  CodeBlock,
  ExpandingContentBox,
  FetchingError,
  Markdown,
} from 'common';
import CommentList from
  'components/marketplace/ImageDetail/ImageDetailComments/CommentList';
import makeRepositoryId from 'lib/utils/repo-image-name';
import routes from 'lib/constants/routes';
import { formatBucketedNumber, bytesToSize } from 'lib/utils/formatNumbers';
import get from 'lodash/get';
const { bool, number, object, shape, string } = PropTypes;
const backButtonProps = {
  pathname: routes.home(),
  text: 'Home',
};

const mapStateToProps = ({ marketplace }, { params }) => {
  const { community } = marketplace && marketplace.images;
  const repositoryId = makeRepositoryId(params);
  return {
    image: community && community[repositoryId] || {},
  };
};

@connect(mapStateToProps)
export default class CommunityImageDetail extends Component {
  static propTypes = {
    image: shape({
      can_edit: bool,
      comments: object,
      description: string,
      error: string,
      full_description: string,
      isFetching: bool,
      last_updated: string,
      namespace: string,
      pull_count: number,
      reponame: string,
      tags: object,
    }),
    location: shape({
      pathname: string,
    }).isRequired,
  }

  renderHeader() {
    const {
      description,
      namespace,
      pull_count,
      reponame,
    } = this.props.image;
    let popularityDisplay;
    if (pull_count) {
      const numPulls = formatBucketedNumber(pull_count);
      const pullOrPulls = numPulls === '1' ? ' Pull' : ' Pulls';
      popularityDisplay = (
        <span><b>{numPulls}</b> {pullOrPulls}</span>
      );
    }
    return (
      <div>
        <div className={css.name}>
          {makeRepositoryId({ namespace, reponame })}
        </div>
        <div className={css.shortDescription}>{description}</div>
        <div className={css.categoryAndPopularity}>
          {popularityDisplay}
        </div>
      </div>
    );
  }

  renderDescription() {
    const { full_description } = this.props.image;
    if (!full_description) {
      return (
        <div>There is no available description for this repository.</div>
      );
    }
    return (
      <ExpandingContentBox>
        <Markdown className={css.scaleImages} rawMarkdown={full_description} />;
      </ExpandingContentBox>
    );
  }

  renderDetailsCard() {
    return (
      <div className={css.detailsCard}>
        <div>
          {this.renderHeader()}
          <hr />
          {this.renderDescription()}
        </div>
      </div>
    );
  }

  renderComments() {
    const { comments, namespace, reponame } = this.props.image;
    if (!comments) {
      return null;
    }
    const { isFetching, count, pages } = comments;
    // This is only page 1 of comments since it is a preview
    const results = isFetching ? [] : pages && pages[1].results;
    const linkTo = routes.communityImageDetailComments({
      namespace,
      reponame,
    });
    return (
      <div className={css.commentsBox}>
        <CommentList
          comments={results}
          count={count}
          isFetching={isFetching}
          isPreview
          linkTo={linkTo}
        />
      </div>
    );
  }

  // Link to view all tags
  renderViewAllVersions() {
    const { namespace, reponame, tags } = this.props.image;
    if (!tags || !tags.count) {
      // Do not display anything if there are no tags
      return null;
    }
    const { count } = tags;
    let numTags;
    if (count) {
      numTags = ` (${count})`;
    }
    const tagsLink = routes.communityImageDetailTags({ namespace, reponame });
    return (
      <div className={css.viewAllVersions}>
        <Link to={tagsLink}>View All Versions {numTags}</Link>
      </div>
    );
  }

  renderPullCommand() {
    const { reponame, namespace } = this.props.image;
    const name = makeRepositoryId({ namespace, reponame });
    // TODO Kristie 3/28/16 Make pull command copyable
    return (
      <div>
        <div className={css.pullCommand}>To pull this image use command</div>
        <CodeBlock className={css.pullCode}>{`docker pull ${name}`}</CodeBlock>
      </div>
    );
  }

  renderTagInfo() {
    const { tags } = this.props.image;
    if (!tags) {
      return null;
    }
    // Find the most recent tag, and display information
    const tag = get(tags, ['pages', 1, 'results', 0], {});
    if (!tag || !tag.name) {
      return null;
    }
    const { full_size, name } = tag;
    let size;
    if (full_size) {
      size = <span> ({bytesToSize(full_size)})</span>;
    }
    return (
      <div className={css.tagsCard}>
        <Card shadow>
          <div className={css.flexRow}>
            <div className={css.tagDetails}>{name}{size}</div>
            {this.renderViewAllVersions()}
          </div>
          <div className={css.notScanned}>Not Docker Verified or Scanned</div>
          <hr />
          {this.renderPullCommand()}
        </Card>
      </div>
    );
  }

  renderIsFetching() {
    // TODO Kristie 5/5/16 Get designs for actual fetching
    return (
      <div>
        <BackButtonArea {...backButtonProps} />
        <div>Fetching...</div>
      </div>
    );
  }

  render() {
    const { error, isFetching } = this.props.image;
    if (isFetching) {
      return this.renderIsFetching();
    }
    // TODO Kristie 3/28/16 Proper 404 Page handling
    if (error) {
      return (
        <div>
          <BackButtonArea {...backButtonProps} />
          <FetchingError resource="this repository" />
        </div>
      );
    }
    return (
      <div>
        <BackButtonArea {...backButtonProps} />
        <div className={css.detailsAndTagsCards}>
          <div>
            {this.renderDetailsCard()}
            {this.renderComments()}
          </div>
          <div>
            {this.renderTagInfo()}
          </div>
        </div>
      </div>
    );
  }
}
