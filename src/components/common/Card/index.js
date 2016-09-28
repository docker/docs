import React, { PropTypes } from 'react';
import css from './styles.css';
import mergeClasses from 'classnames';

const {
  string,
  node,
  bool,
} = PropTypes;

const Card = (props) => {
  /* eslint-disable no-use-before-define */
  const {
    children,
    className,
    ghost,
    hover,
    shadow,
    shy,
    title = '',
    ...restProps,
  } = props;
  /* eslint-enable no-use-before-define */

  const classNames = mergeClasses(
    'dcard',
    className,
    {
      [css.shy]: shy && !ghost,
      [css.card]: !ghost && !shadow,
      [css.shadowCard]: shadow,
      [css.shadowCardHover]: shadow && hover,
    }
  );

  return (
    <div className={classNames} {...restProps} >
      {Boolean(title.length) && <h3 className={css.title}>{title}</h3>}
      <div className={css.content}>{children}</div>
    </div>
  );
};

Card.propTypes = {
  children: node.isRequired,
  className: string,
  ghost: bool,
  hover: bool,
  shadow: bool,
  shy: bool,
  title: string,
};

export default Card;
