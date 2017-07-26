import React from 'react';
import moment from 'moment';
import { Tooltip, Position } from '@blueprintjs/core';

export default function Time({ time }) {
  const abs = time.toISOString();
  const rel = moment(time).fromNow();

  return (
    <div>
      <Tooltip content={abs} position={Position.RIGHT} inline={true}>
        <span>
          {rel}
        </span>
      </Tooltip>
      <style jsx>{`
        div {
          display: inline-block;
        }

        span {
          border-bottom: 1px dotted #eee;
        }
      `}</style>
    </div>
  );
}
