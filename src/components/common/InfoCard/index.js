import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import { Card, InfoIcon } from 'common';
import styles from './styles.css';
import { variants, sizes } from 'lib/constants';
import { isEmpty } from 'lodash';

const { SMALL } = sizes;
const { DULL } = variants;
const { string } = PropTypes;

const InfoCard = ({ text, linkText, url, path, target, ...props }) => {
  const hasLink = !isEmpty(linkText) && (!isEmpty(url) || !isEmpty(path));
  const isExternal = !isEmpty(url);

  let link;
  if (hasLink && isExternal) {
    link = (
      <a
        href={url}
        className={styles.cardLink}
        target={target}
      >
        {linkText}
      </a>
    );
  } else if (hasLink) {
    link = <Link to={path} className={styles.cardLink}>{linkText}</Link>;
  }

  return (
    <Card shy className={styles.wrapper} {...props}>
      <div className={styles.cardContainer}>
        <InfoIcon size={SMALL} variant={DULL} className={styles.icon} />
        <div className={styles.cardBody}>
          <div className={styles.cardText}>{text}</div>
          {link}
        </div>
      </div>
    </Card>
  );
};

InfoCard.propTypes = {
  text: string.isRequired,
  linkText: string,
  url: string,
  path: string,
  target: string,
};

InfoCard.defaultProps = {
  target: '_blank',
};

export default InfoCard;
