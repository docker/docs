import React, { PropTypes } from 'react';
import { Table,
  TableBody,
  TableHeader,
  TableHeaderColumn,
  TableRow,
  TableRowColumn,
} from 'material-ui/Table';
import {
  Link,
} from 'react-router';
import {
  Button,
  Modal,
  Card,
} from 'components/common';
import {
  marketplaceCreateRepository,
  marketplaceDeleteRepository,
} from 'actions/marketplace';
import { billingCreateProduct } from 'actions/billing';
import { connect } from 'react-redux';
import routes from 'lib/constants/routes';
import AdminHomeCreateRepositoryForm from './AdminHomeCreateRepositoryForm';
import get from 'lodash/get';
import css from './styles.css';

const { func, object, array } = PropTypes;

const mapStateToProps = ({ marketplace }, { params }) => {
  const summaries = get(marketplace, ['search', 'pages', 1, 'results'], []);
  return {
    params,
    summaries,
  };
};

const mapDispatch = {
  billingCreateProduct,
  marketplaceCreateRepository,
  marketplaceDeleteRepository,
};

@connect(mapStateToProps, mapDispatch)
export default class Admin extends React.Component {
  static propTypes = {
    params: object,
    summaries: array,
    billingCreateProduct: func,
    marketplaceCreateRepository: func,
    marketplaceDeleteRepository: func,
  }

  static contextTypes = {
    router: object.isRequired,
  }

  state = { isModalOpen: false }

  handleSubmit = (values) => {
    const defaultValues = {
      categories: [],
      platforms: [],
      default_version: {},
    };
    this.props.marketplaceCreateRepository({
      ...defaultValues,
      ...values,
      source: 'official',
    }).then(res => {
      const id = res.value.product_id;
      return this.props.billingCreateProduct({ id, body: {
        name: `${values.namespace}-${values.reponame}`,
        label: values.display_name,
        publisher_id: values.publisher.id,
        rate_plans: [{
          name: 'free',
          label: 'Free',
          duration: 1,
          duration_period: 'year',
          trial: 0,
          trial_period: 'none',
          currency: 'USD',
        }],
      } }).then(() => {
        this.context.router.push(routes.adminRepository({ id }));
      });
    }).catch(err => {
      // TODO (jmorgan): show error text in form if creating a product failed
      console.log(err);
    });
  }

  openModal = () => {
    this.setState({ isModalOpen: true });
  }

  closeModal = () => {
    this.setState({ isModalOpen: false });
  }

  deleteRepository = (id) => () => {
    const result = confirm('Are you sure you want to delete this repository?');
    if (result) {
      this.props.marketplaceDeleteRepository({ id });
    }
  };

  render() {
    return (
      <div>
        <div className={css.header}>
          <h2>{this.props.summaries.length} Repositories</h2>
          <div className={css.new}>
            <div className={css.buttonwrapper}>
              <Button onClick={this.openModal} className={css.nomargin}>
                Add Repository
              </Button>
            </div>
          </div>
        </div>
        <Table selectable={false}>
          <TableHeader displaySelectAll={false} adjustForCheckbox={false}>
            <TableRow>
              <TableHeaderColumn>Name</TableHeaderColumn>
              <TableHeaderColumn>Publisher</TableHeaderColumn>
              <TableHeaderColumn>Repository</TableHeaderColumn>
              <TableHeaderColumn className={css.smallColumn}>
                Status
              </TableHeaderColumn>
              <TableHeaderColumn className={css.smallColumn} />
              <TableHeaderColumn className={css.smallColumn} />
            </TableRow>
          </TableHeader>
          <TableBody displayRowCheckbox={false}>
            {this.props.summaries.map(s => {
              return (
                <TableRow key={s.id}>
                  <TableRowColumn>{s.display_name}</TableRowColumn>
                  <TableRowColumn>{s.publisher.name}</TableRowColumn>
                  <TableRowColumn className={css.monospace}>
                    store/{s.namespace}/{s.reponame}
                  </TableRowColumn>
                  <TableRowColumn className={css.smallColumn}>
                    LIVE
                  </TableRowColumn>
                  <TableRowColumn className={css.smallColumn}>
                    <Link to={routes.adminRepository({ id: s.id })}>
                      Edit →
                    </Link>
                  </TableRowColumn>
                  <TableRowColumn className={css.smallColumn}>
                    <a href="#" onClick={this.deleteRepository(s.id)}>
                      Delete ×
                    </a>
                  </TableRowColumn>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
        <Modal
          isOpen={this.state.isModalOpen}
          onRequestClose={this.closeModal}
          className={css.modal}
        >
          <Card title="Add Repository to the Store">
            <AdminHomeCreateRepositoryForm
              onSubmit={this.handleSubmit}
            />
          </Card>
        </Modal>
      </div>
    );
  }
}
