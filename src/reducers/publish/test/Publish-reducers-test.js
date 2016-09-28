import publishReducer, { DEFAULT_STATE } from 'reducers/publish';
import {
  REPOSITORY_FETCH_IMAGE_TAGS,
  REPOSITORY_FETCH_REPOSITORIES_FOR_NAMESPACE,
} from 'actions/repository';
import {
  PUBLISH_FETCH_PRODUCT_DETAILS,
  PUBLISH_FETCH_PRODUCT_LIST,
  PUBLISH_FETCH_PRODUCT_TIERS,
  PUBLISH_GET_PUBLISHERS,
  PUBLISH_GET_SIGNUP,
  PUBLISH_GET_VENDOR_AGREEMENT,
  PUBLISH_SUBSCRIBE,
} from 'actions/publish';
import { expect } from 'chai';

describe('publish reducer', () => {
  it('should return the initial state', () => {
    expect(publishReducer(undefined, {})).to.deep.equal(DEFAULT_STATE);
  });

  //----------------------------------------------------------------------------
  // PUBLISH_SUBSCRIBE
  //----------------------------------------------------------------------------
  it('should handle PUBLISH_SUBSCRIBE_ACK', () => {
    const payload = {
      key: 'value',
    };
    const action = {
      type: `${PUBLISH_SUBSCRIBE}_ACK`,
      payload,
    };

    const reducer = publishReducer(undefined, action);
    expect(reducer.signup.results).to.deep.equal(payload);
  });

  it('should handle PUBLISH_SUBSCRIBE_ERR', () => {
    const payload = {};
    const action = {
      type: `${PUBLISH_SUBSCRIBE}_ACK`,
      payload,
    };

    const reducer = publishReducer(undefined, action);
    expect(reducer.signup.error).to.exist;
  });

  //----------------------------------------------------------------------------
  // PUBLISH_GET_SIGNUP
  //----------------------------------------------------------------------------
  it('should handle PUBLISH_GET_SIGNUP_ACK', () => {
    const payload = {
      key: 'value',
    };
    const action = {
      type: `${PUBLISH_GET_SIGNUP}_ACK`,
      payload,
    };

    const reducer = publishReducer(undefined, action);
    expect(reducer.signup.error).to.be.empty;
    expect(reducer.signup.results).to.deep.equal(payload);
  });

  it('should handle PUBLISH_GET_SIGNUP_ERR', () => {
    const payload = {
      response: {
        error: {
          message: 'Failed to fetch signup info',
        },
      },
    };
    const action = {
      type: `${PUBLISH_GET_SIGNUP}_ERR`,
      payload,
    };

    const reducer = publishReducer(undefined, action);
    expect(reducer.signup.error).to.exist;
    expect(reducer.signup.results).to.be.empty;
  });

  //----------------------------------------------------------------------------
  // PUBLISH_FETCH_PRODUCT_LIST
  //----------------------------------------------------------------------------
  it('should handle PUBLISH_FETCH_PRODUCT_LIST_ACK', () => {
    const payload = {
      key: 'value',
    };
    const action = {
      type: `${PUBLISH_FETCH_PRODUCT_LIST}_ACK`,
      payload,
    };

    const reducer = publishReducer(undefined, action);
    expect(reducer.productList.error).to.be.empty;
    expect(reducer.productList.results).to.deep.equal(payload);
  });

  it('should handle PUBLISH_FETCH_PRODUCT_LIST_ERR', () => {
    const payload = {
      response: {
        error: {
          message: 'Failed to fetch signup info',
        },
      },
    };
    const action = {
      type: `${PUBLISH_FETCH_PRODUCT_LIST}_ERR`,
      payload,
    };

    const reducer = publishReducer(undefined, action);
    expect(reducer.productList.error).to.exist;
    expect(reducer.productList.results).to.be.empty;
  });

  //----------------------------------------------------------------------------
  // PUBLISH_GET_PUBLISHERS
  //----------------------------------------------------------------------------
  it('should handle PUBLISH_GET_PUBLISHERS_ACK', () => {
    const payload = {
      key: 'value',
    };
    const action = {
      type: `${PUBLISH_GET_PUBLISHERS}_ACK`,
      payload,
    };

    const reducer = publishReducer(undefined, action);
    expect(reducer.publishers.error).to.be.empty;
    expect(reducer.publishers.results).to.deep.equal(payload);
  });

  it('should handle PUBLISH_GET_PUBLISHERS_ERR', () => {
    const payload = {
      response: {
        error: {
          message: 'Failed to fetch signup info',
        },
      },
    };
    const action = {
      type: `${PUBLISH_GET_PUBLISHERS}_ERR`,
      payload,
    };

    const reducer = publishReducer(undefined, action);
    expect(reducer.publishers.error).to.exist;
    expect(reducer.publishers.results).to.be.empty;
  });

  //----------------------------------------------------------------------------
  // PUBLISH_GET_VENDOR_AGREEMENT
  //----------------------------------------------------------------------------
  it('should handle PUBLISH_GET_VENDOR_AGREEMENT_ACK', () => {
    const payload = {
      html: '<html></html>',
    };
    const action = {
      type: `${PUBLISH_GET_VENDOR_AGREEMENT}_ACK`,
      payload,
    };

    const reducer = publishReducer(undefined, action);
    expect(reducer.submit.agreement.results).to.deep.equal(payload.html);
  });

  it('should handle PUBLISH_GET_VENDOR_AGREEMENT_ERR', () => {
    const payload = {
      response: {
        error: {
          message: 'Failed to fetch vendor agreement',
        },
      },
    };
    const action = {
      type: `${PUBLISH_GET_VENDOR_AGREEMENT}_ERR`,
      payload,
    };

    const reducer = publishReducer(undefined, action);
    expect(reducer.submit.agreement.error).to.deep.equal(
      payload.response.error.message
    );
  });

  //--------------------------------------------------------------------------
  // PUBLISH_FETCH_PRODUCT_DETAILS
  //--------------------------------------------------------------------------
  it('should handle PUBLISH_FETCH_PRODUCT_DETAILS_ACK', () => {
    const payload = {
      key: 'value',
    };
    const action = {
      type: `${PUBLISH_FETCH_PRODUCT_DETAILS}_ACK`,
      payload,
    };

    const reducer = publishReducer(undefined, action);
    expect(reducer.currentProductDetails.results).to.deep.equal(payload);
  });

  it('should handle PUBLISH_FETCH_PRODUCT_DETAILS_ERR', () => {
    const payload = {
      response: {
        error: {
          message: 'Failed to fetch submitted repositories',
        },
      },
    };
    const action = {
      type: `${PUBLISH_FETCH_PRODUCT_DETAILS}_ERR`,
      payload,
    };

    const reducer = publishReducer(undefined, action);
    expect(reducer.currentProductDetails.error).to.deep.equal(
      payload.response.error.message
    );
  });

  //--------------------------------------------------------------------------
  // REPOSITORY_FETCH_IMAGE_TAGS
  //--------------------------------------------------------------------------
  it('should handle REPOSITORY_FETCH_IMAGE_TAGS_REQ', () => {
    const testnamespace = 'testnamespace';
    const testreponame = 'testreponame';
    const action = {
      type: `${REPOSITORY_FETCH_IMAGE_TAGS}_REQ`,
      meta: {
        namespace: testnamespace,
        reponame: testreponame,
      },
    };

    const reducer = publishReducer(undefined, action);
    expect(reducer.submit.tags).to.have.keys(testnamespace);
    expect(reducer.submit.tags[testnamespace]).to.have.keys(testreponame);
    expect(reducer.submit.tags[testnamespace][testreponame].isFetching)
      .to.be.true;
  });

  it('should handle REPOSITORY_FETCH_IMAGE_TAGS_ACK', () => {
    const testnamespace = 'testnamespace';
    const testreponame = 'testreponame';
    const payload = {
      results: [{
        name: '1.0',
      }, {
        name: '2.0',
      }],
    };
    const action = {
      type: `${REPOSITORY_FETCH_IMAGE_TAGS}_ACK`,
      meta: {
        namespace: testnamespace,
        reponame: testreponame,
      },
      payload,
    };

    const reducer = publishReducer(undefined, action);
    expect(reducer.submit.tags[testnamespace][testreponame].results)
      .to.deep.equal(['1.0', '2.0']);
    expect(reducer.submit.tags[testnamespace][testreponame].isFetching)
      .to.be.false;
  });

  it('should handle REPOSITORY_FETCH_IMAGE_TAGS_ERR', () => {
    const testnamespace = 'testnamespace';
    const testreponame = 'testreponame';
    const payload = {
      response: {
        error: {
          message: 'Failed to fetch repository tags',
        },
      },
    };
    const action = {
      type: `${REPOSITORY_FETCH_IMAGE_TAGS}_ERR`,
      meta: {
        namespace: testnamespace,
        reponame: testreponame,
      },
      payload,
    };

    const reducer = publishReducer(undefined, action);
    expect(reducer.submit.tags[testnamespace][testreponame].error)
    .to.deep.equal(payload.response.error.message);
    expect(reducer.submit.tags[testnamespace][testreponame].isFetching)
      .to.be.false;
  });

  //--------------------------------------------------------------------------
  // REPOSITORY_FETCH_REPOSITORIES_FOR_NAMESPACE
  //--------------------------------------------------------------------------
  it('should handle REPOSITORY_FETCH_REPOSITORIES_FOR_NAMESPACE_REQ', () => {
    const testnamespace = 'testnamespace';
    const action = {
      type: `${REPOSITORY_FETCH_REPOSITORIES_FOR_NAMESPACE}_REQ`,
      meta: {
        namespace: testnamespace,
      },
    };

    const reducer = publishReducer(undefined, action);
    expect(reducer.submit.repositories).to.have.keys(testnamespace);
    expect(reducer.submit.repositories[testnamespace]).be.an('object');
    expect(reducer.submit.repositories[testnamespace].isFetching).to.be.true;
  });

  it('should handle REPOSITORY_FETCH_REPOSITORIES_FOR_NAMESPACE_ACK', () => {
    const testnamespace = 'testnamespace';
    const payload = {
      results: [{
        name: 'community',
      }, {
        name: 'enterprise',
      }],
    };
    const action = {
      type: `${REPOSITORY_FETCH_REPOSITORIES_FOR_NAMESPACE}_ACK`,
      meta: {
        namespace: testnamespace,
      },
      payload,
    };

    const reducer = publishReducer(undefined, action);
    expect(reducer.submit.repositories[testnamespace].results).to.deep.equal([
      'community',
      'enterprise',
    ]);
    expect(reducer.submit.repositories[testnamespace].isFetching).to.false;
  });

  it('should handle REPOSITORY_FETCH_REPOSITORIES_FOR_NAMESPACE_ERR', () => {
    const testnamespace = 'testnamespace';
    const payload = {
      response: {
        error: {
          message: 'Failed to fetch repositories',
        },
      },
    };
    const action = {
      type: `${REPOSITORY_FETCH_REPOSITORIES_FOR_NAMESPACE}_ERR`,
      meta: {
        namespace: testnamespace,
      },
      payload,
    };

    const reducer = publishReducer(undefined, action);
    expect(reducer.submit.repositories[testnamespace].error).to.deep.equal(
      payload.response.error.message
    );
    expect(reducer.submit.repositories[testnamespace].isFetching).to.false;
  });

  //--------------------------------------------------------------------------
  // PUBLISH_FETCH_PRODUCT_TIERS
  //--------------------------------------------------------------------------
  it('should handle PUBLISH_FETCH_PRODUCT_TIERS_ACK', () => {
    const payload = {
      key: 'value',
    };
    const action = {
      type: `${PUBLISH_FETCH_PRODUCT_TIERS}_ACK`,
      payload,
    };

    const reducer = publishReducer(undefined, action);
    expect(reducer.currentProductTiers.results).to.deep.equal(payload);
  });

  it('should handle PUBLISH_FETCH_PRODUCT_TIERS_ERR', () => {
    const payload = {
      response: {
        error: {
          message: 'Failed to fetch submitted tiers',
        },
      },
    };
    const action = {
      type: `${PUBLISH_FETCH_PRODUCT_TIERS}_ERR`,
      payload,
    };

    const reducer = publishReducer(undefined, action);
    expect(reducer.currentProductTiers.error).to.deep.equal(
      payload.response.error.message
    );
  });
});
