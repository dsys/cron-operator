import React from 'react';
import { Checkbox, Tag, Intent } from '@blueprintjs/core';

export default class JobListFilters extends React.Component {
  handleChangeActive = () => {
    const { onChange, active } = this.props;
    onChange({ active: (active + 2) % 3 - 1 });
  };

  handleChangeSucceeded = () => {
    const { onChange, succeeded } = this.props;
    onChange({ succeeded: (succeeded + 2) % 3 - 1 });
  };

  handleChangeFailed = () => {
    const { onChange, failed } = this.props;
    onChange({ failed: (failed + 2) % 3 - 1 });
  };

  render() {
    const { active, succeeded, failed } = this.props;

    return (
      <div className="filters">
        <div className="checkbox">
          <Checkbox
            checked={active === 1}
            indeterminate={active === 0}
            onChange={this.handleChangeActive}
          >
            <Tag intent={Intent.PRIMARY}>active</Tag>
          </Checkbox>
        </div>
        <div className="checkbox">
          <Checkbox
            checked={succeeded === 1}
            indeterminate={succeeded === 0}
            onChange={this.handleChangeSucceeded}
          >
            <Tag intent={Intent.SUCCESS}>succeeded</Tag>
          </Checkbox>
        </div>
        <div className="checkbox">
          <Checkbox
            checked={failed === 1}
            indeterminate={failed === 0}
            onChange={this.handleChangeFailed}
          >
            <Tag intent={Intent.DANGER}>failed</Tag>
          </Checkbox>
        </div>
        <style jsx>{`
          .filters {
            display: flex;
            margin: 20px 0;
          }

          .checkbox {
            margin-right: 40px;
          }
        `}</style>
      </div>
    );
  }
}
