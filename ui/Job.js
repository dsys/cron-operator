import React from 'react';
import Time from './Time';
import { Intent, Tag, Button } from '@blueprintjs/core';

export default function Job({ data }) {
  const { metadata, status } = data;
  const name = metadata.name;
  const active = status.active || 0;
  const succeeded = status.succeeded || 0;
  const failed = status.failed || 0;
  const startTime = status.startTime ? new Date(status.startTime) : null;
  const completionTime = status.completionTime
    ? new Date(status.completionTime)
    : null;

  return (
    <div className="job">
      <div className="tags">
        {!active &&
          !succeeded &&
          !failed &&
          <div className="tag">
            <Tag>pending</Tag>
          </div>}
        {active
          ? <div className="tag">
              <Tag intent={Intent.PRIMARY}>
                {active} active
              </Tag>
            </div>
          : null}
        {succeeded
          ? <div className="tag">
              <Tag intent={Intent.SUCCESS}>
                {succeeded} succeeded
              </Tag>
            </div>
          : null}
        {failed
          ? <div className="tag">
              <Tag intent={Intent.DANGER}>
                {failed} failed
              </Tag>
            </div>
          : null}
      </div>
      <div className="info">
        <div className="name">
          {name}
        </div>

        <div className="times pt-text-muted">
          {startTime
            ? <div>
                started <Time time={startTime} />
              </div>
            : <div>not started</div>}
          {completionTime
            ? <div>
                completed <Time time={completionTime} />
              </div>
            : <div>not completed</div>}
        </div>
      </div>
      <div className="actions">
        <Button iconName="redo">Re-run</Button>
      </div>
      <style jsx>{`
        .job {
          display: flex;
          border-top: 1px solid #eee;
          padding: 20px;
          align-items: center;
        }

        .tags {
          display: flex;
          flex-direction: column;
          width: 100px;
        }

        .tag {
          margin: 2px 0;
        }

        .info {
          margin: 0 20px;
          flex: 1 0 0;
        }

        .name {
          font-weight: bold;
        }

        .actions {
          display: none;
        }

        .job:hover .actions {
          display: block;
        }
      `}</style>
    </div>
  );
}
