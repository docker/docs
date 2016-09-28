import React, { PropTypes, Component } from 'react';
import { connect } from 'react-redux';
import {
  Card,
  BackButtonArea,
  ToggleSwitch,
} from 'components/common';
import { marketplaceEditRepository } from 'actions/marketplace';
import AdminRepositoryEditMetadataForm from './AdminRepositoryEditMetadataForm';
import css from './styles.css';

const mapStateToProps = ({ marketplace }, { params }) => {
  const details = marketplace.images.certified[params.id] || {};
  details.logo_url = details.logo_url || {};
  return {
    params,
    categories: marketplace.filters.categories,
    platforms: marketplace.filters.platforms,
    details,
  };
};

const mapDispatch = {
  marketplaceEditRepository,
};

@connect(mapStateToProps, mapDispatch)
export default class AdminRepository extends Component {
  static propTypes = {
    params: PropTypes.object.isRequired,
    categories: PropTypes.object.isRequired,
    platforms: PropTypes.object.isRequired,
    details: PropTypes.object.isRequired,
  }

  onSubmit = (values) => {
    marketplaceEditRepository({
      id: this.props.details.id,
      ...values,
    });
    // TODO (jmorgan): update publisher_id and publisher_name in billing
    // once the billing API supports PATCHing products
  };

  render() {
    return (
      <div>
          <BackButtonArea pathname="/admin" text="Back to Repositories" />
          <div className={css.header}>
            <h1>{this.props.details.display_name}</h1>
            <div className={css.switchWrapper}>
              <div>
                <ToggleSwitch
                  toggled
                  label="LIVE"
                  labelPosition="left"
                  disabled
                />
              </div>
            </div>
          </div>
          <Card title="Repository Metadata">
            <AdminRepositoryEditMetadataForm
              onSubmit={this.onSubmit}
              categories={this.props.categories}
              platforms={this.props.platforms}
              initialValues={this.props.details}
            />
          </Card>
      </div>
    );
  }
}
