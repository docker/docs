import React, { Component } from 'react';
import css from './styles.css';
import { Card, AngledTitleBox } from 'common';
import { helpArticles } from 'lib/constants/landingPage';
const {
  // small,
  // large,
  blogPost,
} = helpArticles;

/* eslint-disable max-len */
export default class HelpArticlesCards extends Component {
  // Two card article layout - leaving for when we switch back
  // render() {
  //   return (
  //     <div className={css.articles}>
  //       <Card className={css.helpArticlesCards} shadow>
  //         <AngledTitleBox title="How To" className={css.titleBox} />
  //         <div className={css.helpArticlesHeadline}>{small.name}</div>
  //         <div className={css.helpArticlesDescription}>
  //           {small.description}
  //         </div>
  //       </Card>
  //       <Card className={css.helpArticlesCards} shadow>
  //         <AngledTitleBox title="Case Study" className={css.titleBox} />
  //         <div className={css.helpArticlesHeadline}>{large.name}</div>
  //         <div className={css.helpArticlesDescription}>
  //           {large.description}
  //         </div>
  //       </Card>
  //     </div>
  //   );
  // }

  render() {
    const blogPostLink = 'https://blog.docker.com/2016/06/docker-store/';
    return (
      <div className={css.articles}>
        <Card className={css.helpArticlesCards} shadow>
          <AngledTitleBox title="Blog Post" className={css.titleBox} />
          <div className={css.helpArticlesHeadline}>{blogPost.name}</div>
          <div className={css.helpArticlesDescription}>
            {blogPost.description}
          </div>
          <a href={blogPostLink} target="_blank" className={css.link}>
            Read More
          </a>
        </Card>
      </div>
    );
  }
}
/* eslint-enable */
