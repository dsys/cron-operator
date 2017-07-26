import React from 'react';
import Time from './Time';
import { Intent, Tag, Button } from '@blueprintjs/core';

export default function CronJob({ data }) {
  const schedule = data.spec.schedule;
  const name = data.metadata.name;

  return (
    <div className="cron-job">
      <div className="tags">
        <div className="tag">
          <Tag intent={Intent.WARNING}>
            {schedule}
          </Tag>
        </div>
      </div>
      <div className="info">
        <div className="name">
          {name}
        </div>
        <div className="times pt-text-muted">
          <div>
            last scheduled <Time time={new Date()} />
          </div>
        </div>
      </div>
      <div className="actions">
        <Button iconName="redo">Run Now</Button>
      </div>
      <style jsx>{`
        .cron-job {
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

        .cron-job:hover .actions {
          display: block;
        }
      `}</style>
    </div>
  );
}
